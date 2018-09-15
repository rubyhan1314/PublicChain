package BLC

//getblocks 意为 “给我看一下你有什么区块”（在比特币中，这会更加复杂）
type GetBlocks struct {
	AddrFrom string
}