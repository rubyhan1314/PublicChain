package main

import (
	"./BLC"
	"fmt"
)

func main() {

	//7.测试创世区块存入数据库
	blockchain:=BLC.CreateBlockChainWithGenesisBlock("Genesis Block..")
	fmt.Println(blockchain)
	defer blockchain.DB.Close()
	//8.测试新添加的区块
	blockchain.AddBlockToBlockChain("Send 100RMB to wangergou")
	blockchain.AddBlockToBlockChain("Send 100RMB to lixiaohua")
	blockchain.AddBlockToBlockChain("Send 100RMB to rose")
	fmt.Println(blockchain)
	blockchain.PrintChains()

}