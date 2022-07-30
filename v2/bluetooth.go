package v2

import (
	"errors"
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
)

func FindDevice(adapter *bluetooth.Adapter, name string) (bluetooth.Addresser, error) {
	var result bluetooth.Addresser

	count := 0

	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if device.LocalName() == "GB03" {
			fmt.Printf("GB03 found\n")

			if err := adapter.StopScan(); err != nil {
				fmt.Printf("error stopping scan: %s\n", err)
				return
			}
			result = device.Address
			return
		}

		count += 1
		if count > 100 {
			fmt.Printf("device not found\n")
			if err := adapter.StopScan(); err != nil {
				fmt.Printf("error stopping scan: %s\n", err)
			}
			return
		}
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func SendCommands(queue [][]byte) error {
	var adapter = bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return err
	}

	uuid1, err := bluetooth.ParseUUID("0000ae30-0000-1000-8000-00805f9b34fb")
	if err != nil {
		return err
	}

	/*
		uuid2, err := bluetooth.ParseUUID("0000af30-0000-1000-8000-00805f9b34fb")
		if err != nil {
			return err
		}
	*/

	// TX_CHARACTERISTIC_UUID = '0000ae01-0000-1000-8000-00805f9b34fb'
	// gb01print.py uses:        0000AE01-0000-1000-8000-00805F9B34FB
	uuid3, err := bluetooth.ParseUUID("0000ae01-0000-1000-8000-00805f9b34fb")
	if err != nil {
		return err
	}

	deviceAddr, err := FindDevice(adapter, "GB03")
	if err != nil {
		return err
	}
	if deviceAddr == nil {
		return errors.New("device not found")
	}

	fmt.Printf("connecting to device\n")
	device, err := adapter.Connect(deviceAddr, bluetooth.ConnectionParams{})
	if err != nil {
		return err
	}
	defer func() {
		if err := device.Disconnect(); err != nil {
			fmt.Printf("error disconnecting: %s\n", err)
		}
	}()

	// services I'm finding:
	// ae300000-0000-0000-0000-000000000000
	// ae3a0000-0000-0000-0000-000000000000

	fmt.Printf("discovering services\n")
	services, err := device.DiscoverServices([]bluetooth.UUID{uuid1})
	// services, err := device.DiscoverServices([]bluetooth.UUID{})
	if err != nil {
		return err
	}

	if len(services) == 0 {
		return errors.New("no services")
	}

	for _, service := range services {
		fmt.Printf("discovering characteristics on: %s\n", service)
		chs, err := service.DiscoverCharacteristics([]bluetooth.UUID{})
		if err != nil {
			return err
		}

		for _, ch := range chs {
			fmt.Printf("found characteristic: %s\n", ch)
		}
	}

	/*
		discovering characteristics on: ae300000-0000-0000-0000-000000000000
		found characteristic: ae010000-0000-0000-0000-000000000000
		found characteristic: ae020000-0000-0000-0000-000000000000
		found characteristic: ae030000-0000-0000-0000-000000000000
		found characteristic: ae040000-0000-0000-0000-000000000000
		found characteristic: ae050000-0000-0000-0000-000000000000
		found characteristic: ae100000-0000-0000-0000-000000000000
		discovering characteristics on: ae3a0000-0000-0000-0000-000000000000
		found characteristic: ae3b0000-0000-0000-0000-000000000000
		found characteristic: ae3c0000-0000-0000-0000-000000000000
	*/

	// SEMBRA SIA ROTTO ON MACOS, PORCODIO: https://github.com/tinygo-org/bluetooth/issues/68

	fmt.Printf("trying to send commands\n")
	for _, service := range services {
		fmt.Printf("service: %v\n", service)

		chs, err := service.DiscoverCharacteristics([]bluetooth.UUID{uuid3})
		if err != nil {
			return err
		}

		if len(chs) < 1 {
			fmt.Printf("no characteristics found\n")
			continue
		}

		tx := chs[0]
		for _, cmd := range queue {
			sendbuf := cmd

			for len(sendbuf) != 0 {
				partlen := 20
				if len(sendbuf) < 20 {
					partlen = len(sendbuf)
				}

				part := sendbuf[:partlen]
				sendbuf = sendbuf[partlen:]
				fmt.Printf("sending chunk...\n")
				if _, err := tx.WriteWithoutResponse(part); err != nil {
					return err
				}
				time.Sleep(time.Millisecond * 10)
			}
		}
	}

	return nil
}
