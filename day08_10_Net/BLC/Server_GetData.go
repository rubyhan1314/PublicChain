package BLC
//用于某个块或交易的请求，它可以仅包含一个块或交易的 ID。
type GetData struct {
	AddrFrom string
	Type     string
	Hash       []byte
}
