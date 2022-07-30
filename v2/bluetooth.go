package v2

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

func findDevice(adapter *bluetooth.Adapter, name string) (bluetooth.Addresser, error) {
	var result bluetooth.Addresser

	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if device.LocalName() == "GB03" {
			if err := adapter.StopScan(); err != nil {
				fmt.Printf("error stopping scan: %s\n", err)
				return
			}
			result = device.Address
		}
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func sendCommands(queue [][]byte) error {
	var adapter = bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return err
	}

	uuid1, err := bluetooth.ParseUUID("0000ae30-0000-1000-8000-00805f9b34fb")
	if err != nil {
		return err
	}

	uuid2, err := bluetooth.ParseUUID("0000af30-0000-1000-8000-00805f9b34fb")
	if err != nil {
		return err
	}

	// TX_CHARACTERISTIC_UUID = '0000ae01-0000-1000-8000-00805f9b34fb'
	uuid3, err := bluetooth.ParseUUID("0000ae01-0000-1000-8000-00805f9b34fb")
	if err != nil {
		return err
	}

	deviceAddr, err := findDevice(adapter, "GB03")
	if err != nil {
		return err
	}

	device, err := adapter.Connect(deviceAddr, bluetooth.ConnectionParams{})
	if err != nil {
		return err
	}
	defer func() {
		if err := device.Disconnect(); err != nil {
			fmt.Printf("error disconnecting: %s\n", err)
		}
	}()

	services, err := device.DiscoverServices([]bluetooth.UUID{uuid1, uuid2})
	if err != nil {
		return err
	}

	for _, service := range services {
		fmt.Printf("service: %v\n", service)
		chs, err := service.DiscoverCharacteristics([]bluetooth.UUID{uuid3})
		if err != nil {
			return err
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
				if _, err := tx.WriteWithoutResponse(part); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func sendPayload(payload []byte) error {
	var adapter = bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return err
	}

	/*
		POSSIBLE_SERVICE_UUIDS = [
		'0000ae30-0000-1000-8000-00805f9b34fb',
		'0000af30-0000-1000-8000-00805f9b34fb',
		]
	*/

	uuid1, err := bluetooth.ParseUUID("0000ae30-0000-1000-8000-00805f9b34fb")
	if err != nil {
		return err
	}

	uuid2, err := bluetooth.ParseUUID("0000af30-0000-1000-8000-00805f9b34fb")
	if err != nil {
		return err
	}

	// TX_CHARACTERISTIC_UUID = '0000ae01-0000-1000-8000-00805f9b34fb'
	uuid3, err := bluetooth.ParseUUID("0000ae01-0000-1000-8000-00805f9b34fb")
	if err != nil {
		return err
	}

	err = adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		fmt.Printf("found device: %s %d %s\n", device.Address.String(), device.RSSI, device.LocalName())

		if device.LocalName() == "GB03" {
			// do stuff
			params := bluetooth.ConnectionParams{}
			d, err := adapter.Connect(device.Address, params)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}

			services, err := d.DiscoverServices([]bluetooth.UUID{uuid1, uuid2})
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}

			for _, service := range services {
				chs, err := service.DiscoverCharacteristics([]bluetooth.UUID{uuid3})
				if err != nil {
					fmt.Printf("error: %s\n", err)
					continue
				}

				tx := chs[0]
				sendbuf := payload

				for len(sendbuf) != 0 {
					partlen := 20
					if len(sendbuf) < 20 {
						partlen = len(sendbuf)
					}
					part := sendbuf[:partlen]
					sendbuf = sendbuf[partlen:]
					if _, err := tx.WriteWithoutResponse(part); err != nil {
						fmt.Printf("error writing: %s\n", err)
					}
				}
			}

		}
	})
	if err != nil {
		return err
	}

	return nil
}
