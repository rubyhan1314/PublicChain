package main

import (
	"mypublicchain/day01_01_Base_Prototype/BLC"
	"fmt"
	"crypto/sha256"
)

func main() {
	//1.测试Block
	//block:=BLC.NewBlock("I am a block",make([]byte,32,32),1)
	//fmt.Printf("Heigth:%x\n",block.Height)
	//fmt.Printf("Data:%s\n",block.Data)

	//2.测试创世区块
	//genesisBlock :=BLC.CreateGenesisBlock("Genesis Block..")
	//fmt.Printf("Heigth:%x\n",genesisBlock.Height)
	//fmt.Printf("PrevBlockHash:%x\n",genesisBlock.PrevBlockHash)
	//fmt.Printf("Data:%s\n",genesisBlock.Data)

	//3.测试区块链
	//genesisBlockChain := BLC.CreateBlockChainWithGenesisBlock("Genesis Block..")
	//fmt.Println(genesisBlockChain)
	//fmt.Println(genesisBlockChain.Blocks)
	//fmt.Printf("Heigth:%x\n",genesisBlockChain.Blocks[0].Height)
	//fmt.Printf("PrevBlockHash:%x\n",genesisBlockChain.Blocks[0].PrevBlockHash)
	//fmt.Printf("Data:%s\n",genesisBlockChain.Blocks[0].Data)

	//4.测试添加新区块
	blockChain:=BLC.CreateBlockChainWithGenesisBlock("Genesis Block..")
	blockChain.AddBlockToBlockChain("Send 1BTC To Wangergou",blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,blockChain.Blocks[len(blockChain.Blocks)-1].Hash)
	blockChain.AddBlockToBlockChain("Send 3BTC To lixiaohua",blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,blockChain.Blocks[len(blockChain.Blocks)-1].Hash)
	blockChain.AddBlockToBlockChain("Send 5BTC To rose",blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,blockChain.Blocks[len(blockChain.Blocks)-1].Hash)

	for _, block := range blockChain.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}

	hash:=sha256.Sum256([]byte("HelloWorld"))
	fmt.Printf("%x\n",hash)
}
