package bt

import (
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
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

	characteristics, err := printService.DiscoverCharacteristics([]bluetooth.UUID{*WriteUUID})
	if err != nil {
		return fmt.Errorf("failed to get write characteristic: %w", err)
	}

	tx := characteristics[0]

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
	fmt.Println("Waiting 30 seconds for the printer to finish")
	time.Sleep(30 * time.Second)

	return nil
}
