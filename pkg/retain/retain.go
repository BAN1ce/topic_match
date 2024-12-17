package retain

const (
	HandingSendRetain = byte(iota)
	HandingSendRetainOnlyNewConnection
	HandingNoSendRetain
)
