package bt

import "tinygo.org/x/bluetooth"

var (
	PrintServiceUUID *bluetooth.UUID
	WriteUUID        *bluetooth.UUID
	NotificationUUID *bluetooth.UUID
)

func mustParseUUID(s string) *bluetooth.UUID {
	uuid, err := bluetooth.ParseUUID(s)
	if err != nil {
		panic(err)
	}
	return &uuid
}
