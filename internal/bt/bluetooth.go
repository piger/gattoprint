package bt

import (
	"fmt"
	"log"
	"time"

	"tinygo.org/x/bluetooth"
)

const (
	printWaitTime = 30 * time.Second      // time to wait for the printer to finish printing.
	sendDelay     = 10 * time.Millisecond // time to wait between each write command sent to the printer.
)

// chunks split the slice `s` in chunks of the given size.
func chunks(s []byte, size int) [][]byte {
	var result [][]byte
	l := len(s)

	for i := 0; i < l; i += size {
		end := i + size
		if end > l {
			end = l
		}
		result = append(result, s[i:end])
	}

	return result
}

// SendCommands reads from the `commands` channel and send commands to the printer,
// ensuring a certain chunk size (20 bytes) and a small delay between each write.
func SendCommands(adapter *bluetooth.Adapter, address bluetooth.Addresser, commands chan []byte) error {
	device, err := adapter.Connect(address, bluetooth.ConnectionParams{})
	if err != nil {
		return fmt.Errorf("failed to connect to printer: %w", err)
	}
	defer func() {
		if err := device.Disconnect(); err != nil {
			fmt.Printf("warning: error disconnecting from printer: %s", err)
		}
	}()

	services, err := device.DiscoverServices([]bluetooth.UUID{*PrintServiceUUID})
	if err != nil {
		return fmt.Errorf("failed to get print service: %w", err)
	}

	printService := services[0]

	chs, err := printService.DiscoverCharacteristics([]bluetooth.UUID{*WriteUUID, *NotificationUUID})
	if err != nil {
		return fmt.Errorf("failed to get write characteristic: %w", err)
	}

	tx, notif := chs[0], chs[1]
	notifChan := make(chan struct{}, 1)

	if err := notif.EnableNotifications(func(buf []byte) {
		notifChan <- struct{}{}
	}); err != nil {
		return fmt.Errorf("error enabling notifications: %w", err)
	}

	log.Println("sending commands to printer")

	for cmd := range commands {
		fmt.Print(".")

		for _, chunk := range chunks(cmd, 20) {
			fmt.Print("+")
			if _, err := tx.WriteWithoutResponse(chunk); err != nil {
				return err
			}
			time.Sleep(sendDelay)
		}
	}

	fmt.Println()
	log.Println("waiting for the printer to finish printing")

	t := time.NewTimer(printWaitTime)
	defer t.Stop()

	select {
	case <-t.C:
		log.Printf("%s passed but the printer didn't signal that finished printing", printWaitTime)
	case <-notifChan:
	}

	return nil
}
