package BLC

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)
//step2:修改Block的交易类型
type Block struct {
	//字段：
	//高度Height：其实就是区块的编号，第一个区块叫创世区块，高度为0
	Height int64
	//上一个区块的哈希值ProvHash：
	PrevBlockHash []byte
	//交易数据Data：目前先设计为[]byte,后期是Transaction
	//Data [] byte
	Txs [] *Transaction
	//时间戳TimeStamp：
	TimeStamp int64
	//哈希值Hash：32个的字节，64个16进制数
	Hash []byte

	Nonce int64
}

func NewBlock(txs []*Transaction,provBlockHash []byte,height int64) *Block{
	//创建区块
	block:=&Block{height,provBlockHash,txs,time.Now().Unix(),nil,0}
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

func CreateGenesisBlock(txs []*Transaction) *Block{
	return NewBlock(txs,make([] byte,32,32),0)
}

//将区块序列化，得到一个字节数组---区块的行为，设计为方法
func (block *Block) Serilalize() []byte {
	//1.创建一个buffer
	var result bytes.Buffer
	//2.创建一个编码器
	encoder := gob.NewEncoder(&result)
	//3.编码--->打包
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

//反序列化，得到一个区块---设计为函数
func DeserializeBlock(blockBytes [] byte) *Block {
	var block Block
	var reader = bytes.NewReader(blockBytes)
	//1.创建一个解码器
	decoder := gob.NewDecoder(reader)
	//解包
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

//step4：新增方法
//将Txs转为[]byte
func (block *Block) HashTransactions()[]byte{
	var txHashes [][] byte
	var txHash [32]byte
	for _,tx :=range block.Txs{
		txHashes = append(txHashes,tx.TxID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes,[]byte{}))
	return txHash[:]
}










