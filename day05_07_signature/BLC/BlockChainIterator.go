package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

//1.新增一个结构体
type BlockChainIterator struct {
	CurrentHash [] byte //当前区块的hash
	DB          *bolt.DB //数据库
}




//3.获取区块
func (bcIterator *BlockChainIterator) Next() *Block {
	block:=new(Block)
	//1.打开数据库并读取
	err :=bcIterator.DB.View(func(tx *bolt.Tx) error {
		//2.打开数据表
		b:=tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil{
			//3.根据当前hash获取数据并反序列化
			blockBytes:=b.Get(bcIterator.CurrentHash)
			block = DeserializeBlock(blockBytes)
			//4.更新当前的hash
			bcIterator.CurrentHash = block.PrevBlockHash
		}

		return nil
	})
	if err != nil{
		log.Panic(err)
	}
	return block
}


