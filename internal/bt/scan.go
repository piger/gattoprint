package bt

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
	"tinygo.org/x/bluetooth"
)

func FindDevice(name string, adapter *bluetooth.Adapter) (bluetooth.Addresser, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("error configuring terminal: %w", err)
	}
	defer func() {
		if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil {
			fmt.Printf("error restoring terminal state: %s", err)
		}
	}()

	keyPress := make(chan struct{}, 1)
	go func() {
		b := make([]byte, 1)
		if _, err := os.Stdin.Read(b); err != nil {
			// let's continue instead of returning here, so that
			// the channel is not left hanging forever.
			fmt.Printf("error reading stdin: %s", err)
		}
		keyPress <- struct{}{}
	}()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	errChan := make(chan error, 1)
	jobChan := make(chan bluetooth.Addresser, 1)
	go func() {
		err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
			if result.LocalName() == name {
				if err := adapter.StopScan(); err != nil {
					fmt.Printf("error stopping scan: %s", err)
				}
				jobChan <- result.Address
			}
		})
		if err != nil {
			errChan <- fmt.Errorf("error scanning for BLE: %w", err)
		}
	}()

	// progress bar
	chars := []string{"|", "/", "-", "\\"}
	l := len(chars)
	idx := 0

	// clear progress bar
	defer func() {
		fmt.Printf("\r%s\r", strings.Repeat(" ", 50))
	}()

	for {
		select {
		case <-keyPress:
			return nil, errors.New("scan cancelled")

		case result := <-jobChan:
			return result, nil

		case <-ticker.C:
			fmt.Printf("\r%s Scanning...", chars[idx])
			idx++
			if idx >= l {
				idx = 0
			}

		case err := <-errChan:
			return nil, err
		}
	}
}
