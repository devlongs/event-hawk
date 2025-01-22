package notifications

type Event struct {
	BlockNumber uint64
	TxHash      string
	Data        map[string]interface{}
}
