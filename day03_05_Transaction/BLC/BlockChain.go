package BLC

import (
	"github.com/boltdb/bolt"
	"os"
	"fmt"
	"log"
	"math/big"
	"time"
	"strconv"
	"encoding/hex"
)

type BlockChain struct {
	//Blocks []*Block //存储有序的区块
	Tip [] byte  // 最近的取快递Hash值
	DB  *bolt.DB //数据库对象
}

//修改该方法
/*
1.仅仅用来创建区块链
如果数据库存在，证明区块链存在，直接结束该方法
否则进行创建创世区块，并存入数据库中
 */
func CreateBlockChainWithGenesisBlock(address string) {
	if dbExists() {
		fmt.Println("数据库已经存在。。。")
		return
	}

	//
	fmt.Println("创建创世区块：")
	//2.数据库不存在，说明第一次创建，然后存入到数据库中
	fmt.Println("数据库不存在。。")
	//A：创建创世区块
	//创建创世区块
	//先创建coinbase交易
	txCoinBase := NewCoinBaseTransaction(address)
	genesisBlock := CreateGenesisBlock([]*Transaction{txCoinBase})
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
	//return &BlockChain{genesisBlock.Hash, db}
}

func (bc *BlockChain) AddBlockToBlockChain(txs [] *Transaction) {
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
			newBlock := NewBlock(txs, lastBlock.Hash, lastBlock.Height+1)
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

//判断数据库是否存在
func dbExists() bool {
	if _, err := os.Stat(DBNAME); os.IsNotExist(err) {
		return false
	}
	return true
}

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

	//2.循环迭代
	for {
		block := bcIterator.Next()
		fmt.Printf("第%d个区块的信息：\n", block.Height+1)
		//获取当前hash对应的数据，并进行反序列化
		fmt.Printf("\t高度：%d\n", block.Height)
		fmt.Printf("\t上一个区块的hash：%x\n", block.PrevBlockHash)
		fmt.Printf("\t当前的hash：%x\n", block.Hash)
		//fmt.Printf("\t数据：%v\n", block.Txs)
		fmt.Println("\t交易：")
		for _, tx := range block.Txs {
			fmt.Printf("\t\t交易ID：%x\n", tx.TxID)
			fmt.Println("\t\tVins:")
			for _, in := range tx.Vins {
				fmt.Printf("\t\t\tTxID:%x\n", in.TxID)
				fmt.Printf("\t\t\tVout:%d\n", in.Vout)
				fmt.Printf("\t\t\tScriptSiq:%s\n", in.ScriptSiq)
			}
			fmt.Println("\t\tVouts:")
			for _, out := range tx.Vouts {
				fmt.Printf("\t\t\tvalue:%d\n", out.Value)
				fmt.Printf("\t\t\tScriptPubKey:%s\n", out.ScriptPubKey)
			}
		}

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

//新增方法，用于获取区块链
func GetBlockchainObject() *BlockChain {
	/*
	1.如果数据库不存在，直接返回nil
	2.读取数据库
	 */
	if !dbExists() {
		fmt.Println("数据库不存在，无法获取区块链。。")
		return nil
	}

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

//挖掘新的区块
func (bc *BlockChain) MineNewBlock(from, to, amount []string) {
	/*
	./bc send -from '["wangergou"]' -to '["lixiaohua"]' -amount '["4"]'
["wangergou"]
["lixiaohua"]
["4"]

	 */
	//fmt.Println(from)
	//fmt.Println(to)
	//fmt.Println(amount)
	//1.新建交易
	//2.新建区块
	//3.将区块存入到数据库
	var txs []*Transaction
	for i := 0; i < len(from); i++ {

		amountInt, _ := strconv.ParseInt(amount[i], 10, 64)
		tx := NewSimpleTransaction(from[i], to[i], amountInt, bc, txs)

		txs = append(txs, tx)
	}

	//amountInt, _ := strconv.ParseInt(amount[0], 10, 64)
	//
	//tx := NewSimpleTransaction(from[0], to[0], amountInt, bc)
	//
	//txs = append(txs, tx)

	var block *Block    //数据库中的最后一个block
	var newBlock *Block //要创建的新的block
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			hash := b.Get([] byte("l"))
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes) //数据库中的最后一个block
		}
		return nil
	})

	newBlock = NewBlock(txs, block.Hash, block.Height+1)

	bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			b.Put(newBlock.Hash, newBlock.Serilalize())
			b.Put([]byte("l"), newBlock.Hash)
			bc.Tip = newBlock.Hash
		}
		return nil
	})

}

//找到所有未花费的交易输出
func (bc *BlockChain) UnUTXOs(address string, txs []*Transaction) []*UTXO {
	/*
	1.先遍历未打包的交易(参数txs)，找出未花费的Output。
	2.遍历数据库，获取每个块中的Transaction，找出未花费的Output。
	 */
	var unUTXOs []*UTXO                      //未花费
	spentTxOutputs := make(map[string][]int) //存储已经花费

	//1.添加先从txs遍历，查找未花费
	//for i, tx := range txs {
	for i:=len(txs)-1;i>=0;i--{
		unUTXOs = caculate(txs[i], address, spentTxOutputs, unUTXOs)
	}

	bcIterator := bc.Iterator()
	for {
		block := bcIterator.Next()
		//统计未花费
		//2.获取block中的每个Transaction
		for i := len(block.Txs) - 1; i >= 0; i-- {
			unUTXOs = caculate(block.Txs[i], address, spentTxOutputs, unUTXOs)
		}

		//结束迭代
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}
	}
	return unUTXOs
}

func (bc *BlockChain) GetBalance(address string, txs []*Transaction) int64 {
	//txOutputs:=bc.UnUTXOs(address)
	unUTXOs := bc.UnUTXOs(address, txs)
	//fmt.Println(address, unUTXOs)
	var amount int64
	for _, utxo := range unUTXOs {
		amount = amount + utxo.Output.Value
	}
	return amount

}

//转账时查获在可用的UTXO
func (bc *BlockChain) FindSpendableUTXOs(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	/*
	1.获取所有的UTXO
	2.遍历UTXO

	返回值：map[hash]{index}
	 */

	var balance int64
	utxos := bc.UnUTXOs(from, txs)
	//fmt.Println(from,utxos)
	spendableUTXO := make(map[string][]int)
	for _, utxo := range utxos {
		balance += utxo.Output.Value
		hash := hex.EncodeToString(utxo.TxID)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if balance >= amount {
			break
		}
	}
	if balance < amount {
		fmt.Printf("%s 余额不足。。总额：%d，需要：%d\n", from,balance,amount)
		os.Exit(1)
	}
	return balance, spendableUTXO

}

func caculate(tx *Transaction, address string, spentTxOutputs map[string][]int, unUTXOs []*UTXO) []*UTXO {
	//2.先遍历TxInputs，表示花费
	if !tx.IsCoinbaseTransaction() {
		for _, in := range tx.Vins {
			//如果解锁
			if in.UnLockWithAddress(address) {
				key := hex.EncodeToString(in.TxID)
				spentTxOutputs[key] = append(spentTxOutputs[key], in.Vout)
			}
		}
	}

	//fmt.Println("===>", spentTxOutputs)
	//3.遍历TxOutputs
outputs:
	for index, out := range tx.Vouts {
		if out.UnLockWithAddress(address) {
			//fmt.Println("height,", block.Height, ",index---", index, out, "map-->", spentTxOutputs, len(spentTxOutputs))
			//如果对应的花费容器中长度不为0,
			if len(spentTxOutputs) != 0 {
				var isSpentUTXO bool

				for txID, indexArray := range spentTxOutputs {
					for _, i := range indexArray {
						if i == index && txID == hex.EncodeToString(tx.TxID) {
							isSpentUTXO = true
							continue outputs
						}
					}
				}
				if !isSpentUTXO {
					utxo := &UTXO{tx.TxID, index, out}
					unUTXOs = append(unUTXOs, utxo)
					//unSpentTxOutputs = append(unSpentTxOutputs, out)
				}

			} else {
				utxo := &UTXO{tx.TxID, index, out}
				unUTXOs = append(unUTXOs, utxo)
				//unSpentTxOutputs = append(unSpentTxOutputs, out)
			}
			//fmt.Println(block.Height, "   ", index, "----....", unUTXOs)
		}
	}
	return unUTXOs
}
