package bt

func init() {
	PrintServiceUUID = mustParseUUID("ae300000-0000-0000-0000-000000000000")
	WriteUUID = mustParseUUID("ae010000-0000-0000-0000-000000000000")
	NotificationUUID = mustParseUUID("ae020000-0000-0000-0000-000000000000")
}
