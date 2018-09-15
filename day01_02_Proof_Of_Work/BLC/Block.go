package BLC

import (
	"time"

)

//step1:创建Block结构体
type Block struct {
	//字段：
	//高度Height：其实就是区块的编号，第一个区块叫创世区块，高度为0
	Height int64
	//上一个区块的哈希值ProvHash：
	PrevBlockHash []byte
	//交易数据Data：目前先设计为[]byte,后期是Transaction
	Data [] byte
	//时间戳TimeStamp：
	TimeStamp int64
	//哈希值Hash：32个的字节，64个16进制数
	Hash []byte

	Nonce int64
}

//step2：创建新的区块
func NewBlock(data string,provBlockHash []byte,height int64) *Block{
	//创建区块
	block:=&Block{height,provBlockHash,[]byte(data),time.Now().Unix(),nil,0}
	//step5：设置block的hash和nonce
	//设置哈希
	//block.SetHash()
	//调用工作量证明的方法，并且返回有效的Hash和Nonce
	pow:=NewProofOfWork(block)
	hash,nonce:=pow.Run()
	block.Hash = hash
	block.Nonce = nonce


	return block
}


//step4:创建创世区块：
func CreateGenesisBlock(data string) *Block{
	return NewBlock(data,make([] byte,32,32),0)
}










