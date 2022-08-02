package bt

// It looks like the `cbgo` library doesn't support 16 bits UUID, leading
// to macOS reading different UUIDs from the printer than what is read on
// Linux; see also: https://github.com/tinygo-org/bluetooth/issues/68

func init() {
	PrintServiceUUID = mustParseUUID("ae300000-0000-0000-0000-000000000000")
	WriteUUID = mustParseUUID("ae010000-0000-0000-0000-000000000000")
	NotificationUUID = mustParseUUID("ae020000-0000-0000-0000-000000000000")
}
