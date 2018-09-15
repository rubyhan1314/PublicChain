package main

import (
	"./BLC"
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
	//blockChain:=BLC.CreateBlockChainWithGenesisBlock()
	//blockChain.AddBlockToBlockChain("Send 100RMB To Wangergou",blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,blockChain.Blocks[len(blockChain.Blocks)-1].Hash)
	//blockChain.AddBlockToBlockChain("Send 300RMB To lixiaohua",blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,blockChain.Blocks[len(blockChain.Blocks)-1].Hash)
	//blockChain.AddBlockToBlockChain("Send 500RMB To rose",blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,blockChain.Blocks[len(blockChain.Blocks)-1].Hash)
	//
	//fmt.Println(blockChain)

	//5.测试序列化和反序列化
	//block:=BLC.NewBlock("helloworld",make([]byte,32,32),0)
	//data:=block.Serilalize()
	//fmt.Println(block)
	//fmt.Println(data)
	//block2:=BLC.DeserializeBlock(data)
	//fmt.Println(block2)

	//6.创建区块，存入数据库
	//打开数据库
	//block:=BLC.NewBlock("helloworld",make([]byte,32,32),0)
	//db,err := bolt.Open("my.db",0600,nil)
	//if err != nil{
	//	log.Fatal(err)
	//}
	//
	//defer db.Close()
	//
	//err = db.Update(func(tx *bolt.Tx) error {
	//	//获取bucket，没有就创建新表
	//	b := tx.Bucket([]byte("blocks"))
	//	if b == nil{
	//		b,err = tx.CreateBucket([] byte("blocks"))
	//		if err !=nil{
	//			log.Panic("创建表失败")
	//		}
	//	}
	//	//添加数据
	//	err  = b.Put([]byte("l"),block.Serilalize())
	//	if err !=nil{
	//		log.Panic(err)
	//	}
	//
	//	return nil
	//})
	//if err != nil{
	//	log.Panic(err)
	//}
	//err = db.View(func(tx *bolt.Tx) error {
	//	b := tx.Bucket([]byte("blocks"))
	//	if b !=nil{
	//		data := b.Get([]byte("l"))
	//		//fmt.Printf("%s\n",data)//直接打印会乱码
	//		//反序列化
	//		block2:=BLC.DeserializeBlock(data)
	//		//fmt.Println(block2)
	//		fmt.Printf("%v\n",block2)
	//
	//	}
	//	return nil
	//})

	//7.测试创世区块存入数据库
	//blockchain:=BLC.CreateBlockChainWithGenesisBlock("Genesis Block..")
	//fmt.Println(blockchain)
	//defer blockchain.DB.Close()
	//8.测试新添加的区块
	//blockchain.AddBlockToBlockChain("Send 100RMB to wangergou")
	//blockchain.AddBlockToBlockChain("Send 100RMB to lixiaohua")
	//blockchain.AddBlockToBlockChain("Send 100RMB to rose")
	//fmt.Println(blockchain)
	//blockchain.PrintChains()

	//9.CLI操作
	cli:=BLC.CLI{}
	cli.Run()

	//outsMap := make(map[string]*BLC.TxOutputs)
	//outputs:=outsMap["aa"]
	//fmt.Println(outputs)
	//fmt.Println(len(outputs.UTXOS))


}