package BLC

import (
	"github.com/boltdb/bolt"
	"os"
	"fmt"
	"log"
	"math/big"
	"time"
)

//step1：修改BlockChain的结构体
type BlockChain struct {
	//Blocks []*Block //存储有序的区块
	Tip [] byte  // 最近的取快递Hash值
	DB  *bolt.DB //数据库对象
}

//step2：修改该方法
func CreateBlockChainWithGenesisBlock(data string) *BlockChain {
	//1.先判断数据库是否存在，如果有，从数据库读取
	if dbExists() {
		fmt.Println("数据库已经存在。。")
		//A：打开数据库
		db, err := bolt.Open(DBNAME, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		//defer db.Close()
		var blockchain *BlockChain
		//B：读取数据库
		err = db.View(func(tx *bolt.Tx) error {
			//C：打开表
			b := tx.Bucket([]byte(BLOCKTABLENAME))
			if b != nil {
				//D：读取最后一个hash
				hash := b.Get([]byte("l"))
				//E：创建blockchain
				blockchain = &BlockChain{hash, db}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		return blockchain
	}

	//2.数据库不存在，说明第一次创建，然后存入到数据库中
	fmt.Println("数据库不存在。。")
	//A：创建创世区块
	//创建创世区块
	genesisBlock := CreateGenesisBlock(data)
	//B：打开数据库
	db, err := bolt.Open(DBNAME, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	//C：存入数据表
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(BLOCKTABLENAME))
		if err != nil {
			log.Panic(err)
		}
		if b != nil {
			err = b.Put(genesisBlock.Hash, genesisBlock.Serilalize())
			if err != nil {
				log.Panic("创世区块存储有误。。。")
			}
			//存储最新区块的hash
			b.Put([]byte("l"), genesisBlock.Hash)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	//返回区块链对象
	return &BlockChain{genesisBlock.Hash, db}
}

//step4：修改该方法
func (bc *BlockChain) AddBlockToBlockChain(data string) {
	//创建新区块
	//newBlock := NewBlock(data,prevHash,height)
	//添加到切片中
	//bc.Blocks = append(bc.Blocks,newBlock)
	//1.更新数据库
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		//2.打开表
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			//2.根据最新块的hash读取数据，并反序列化最后一个区块
			blockBytes := b.Get(bc.Tip)
			lastBlock := DeserializeBlock(blockBytes)
			//3.创建新的区块
			newBlock := NewBlock(data, lastBlock.Hash, lastBlock.Height+1)
			//4.将新的区块序列化并存储
			err := b.Put(newBlock.Hash, newBlock.Serilalize())
			if err != nil {
				log.Panic(err)
			}
			//5.更新最后一个哈希值，以及blockchain的tip
			b.Put([]byte("l"), newBlock.Hash)
			bc.Tip = newBlock.Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

}

//step3：
//判断数据库是否存在
func dbExists() bool {
	if _, err := os.Stat(DBNAME); os.IsNotExist(err) {
		return false
	}
	return true
}

//step5：新增方法，遍历数据库，打印输出所有的区块信息
/*
func (bc *BlockChain) PrintChains() {
	//1.根据bc的tip，获取最新的hash值，表示当前的hash
	var currentHash = bc.Tip
	//2.循环，根据当前hash读取数据，反序列化得到最后一个区块
	var count = 0
	block := new(Block) // var block *Block
	for {
		err := bc.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(BLOCKTABLENAME))

			if b != nil {
				count++
				fmt.Printf("第%d个区块的信息：\n", count)
				//获取当前hash对应的数据，并进行反序列化
				blockBytes := b.Get(currentHash)
				block = DeserializeBlock(blockBytes)
				fmt.Printf("\t高度：%d\n", block.Height)
				fmt.Printf("\t上一个区块的hash：%x\n", block.PrevBlockHash)
				fmt.Printf("\t当前的hash：%x\n", block.Hash)
				fmt.Printf("\t数据：%s\n", block.Data)
				//fmt.Printf("\t时间：%v\n", block.TimeStamp)
				fmt.Printf("\t时间：%s\n",time.Unix(block.TimeStamp,0).Format("2006-01-02 15:04:05"))
				fmt.Printf("\t次数：%d\n", block.Nonce)
			}

			return nil
		})
		if err != nil {
			log.Panic(err)
		}
		//3.直到父hash值为0
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}
		//4.更新当前区块的hash值
		currentHash = block.PrevBlockHash
	}
}
*/

//2.获取一个迭代器的方法
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.Tip, bc.DB}
}

func (bc *BlockChain) PrintChains() {
	//1.获取迭代器对象
	bcIterator := bc.Iterator()

	var count = 0
	//2.循环迭代
	for {
		block := bcIterator.Next()
		count++
		fmt.Printf("第%d个区块的信息：\n", count)
		//获取当前hash对应的数据，并进行反序列化
		fmt.Printf("\t高度：%d\n", block.Height)
		fmt.Printf("\t上一个区块的hash：%x\n", block.PrevBlockHash)
		fmt.Printf("\t当前的hash：%x\n", block.Hash)
		fmt.Printf("\t数据：%s\n", block.Data)
		//fmt.Printf("\t时间：%v\n", block.TimeStamp)
		fmt.Printf("\t时间：%s\n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("\t次数：%d\n", block.Nonce)

		//3.直到父hash值为0
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}
	}
}

