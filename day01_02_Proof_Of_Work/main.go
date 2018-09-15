package main

import (
	"./BLC"
	"fmt"
	"strconv"
)

func main() {
	//1.测试Block
	//block:=BLC.NewBlock("I am a block",make([]byte,32,32),1)
	//fmt.Println(block)
	//2.测试创世区块
	//genesisBlock :=BLC.CreateGenesisBlock("Genesis Block..")
	//fmt.Println(genesisBlock)

	//3.测试区块链
	//genesisBlockChain := BLC.CreateBlockChainWithGenesisBlock()
	//fmt.Println(genesisBlockChain)
	//fmt.Println(genesisBlockChain.Blocks)
	//fmt.Println(genesisBlockChain.Blocks[0])

	//4.测试添加新区块
	blockChain:=BLC.CreateBlockChainWithGenesisBlock("Genesis Block..")
	blockChain.AddBlockToBlockChain("Send 100RMB To Wangergou",blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,blockChain.Blocks[len(blockChain.Blocks)-1].Hash)
	blockChain.AddBlockToBlockChain("Send 300RMB To lixiaohua",blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,blockChain.Blocks[len(blockChain.Blocks)-1].Hash)
	blockChain.AddBlockToBlockChain("Send 500RMB To rose",blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,blockChain.Blocks[len(blockChain.Blocks)-1].Hash)

	fmt.Println(blockChain)
	for _, block := range blockChain.Blocks {
		pow := BLC.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.IsValid()))
	}

	/*
	// 5.检测pow
	//1.创建一个big对象 0000000.....00001
	target := big.NewInt(1)
	fmt.Printf("0x%x\n",target) //0x1

	//2.左移256-bits位
	target = target.Lsh(target, 256-BLC.TargetBit)

	fmt.Printf("0x%x\n",target) //61
	//61位：0x1000000000000000000000000000000000000000000000000000000000000
	//64位：0x0001000000000000000000000000000000000000000000000000000000000000

	s1:="HelloWorld"
	hash:=sha256.Sum256([]byte(s1))
	fmt.Printf("0x%x\n",hash)
	*/
}
