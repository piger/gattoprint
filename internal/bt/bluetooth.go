package bt

import (
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
)

const (
	printWaitTime = 30 * time.Second
)

type PrinterCommand []byte

func SendCommands(adapter *bluetooth.Adapter, address bluetooth.Addresser, commands [][]byte) error {
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

	characteristics, err := printService.DiscoverCharacteristics([]bluetooth.UUID{*WriteUUID, *NotificationUUID})
	if err != nil {
		return fmt.Errorf("failed to get write characteristic: %w", err)
	}

	tx, notif := characteristics[0], characteristics[1]
	notifChan := make(chan struct{}, 1)

	if err := notif.EnableNotifications(func(buf []byte) {
		notifChan <- struct{}{}
	}); err != nil {
		return fmt.Errorf("error enabling notifications: %w", err)
	}

	fmt.Println("Sending commands:")
	for _, cmd := range commands {
		sendbuf := cmd

		fmt.Print(".")

		for len(sendbuf) != 0 {
			partlen := 20
			if len(sendbuf) < 20 {
				partlen = len(sendbuf)
			}

			part := sendbuf[:partlen]
			sendbuf = sendbuf[partlen:]
			fmt.Print("+")
			if _, err := tx.WriteWithoutResponse(part); err != nil {
				return err
			}
			time.Sleep(time.Millisecond * 10)
		}
	}
	fmt.Println()
	fmt.Println("Waiting for the printer to finish printing...")

	t := time.NewTimer(printWaitTime)
	defer t.Stop()

	select {
	case <-t.C:
		fmt.Printf("%s passed but the printer didn't signal that finished printing", printWaitTime)
	case <-notifChan:
	}

	return nil
}
