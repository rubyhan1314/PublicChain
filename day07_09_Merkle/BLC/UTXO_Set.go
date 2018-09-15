package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
	"os"
	"bytes"
)

type UTXOSet struct {
	BlockChain *BlockChain
}

const utxoTableName = "utxoTable"

//重置数据库表
func (utxoSet *UTXOSet) ResetUTXOSet() {
	err := utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			err := tx.DeleteBucket([]byte(utxoTableName))
			if err != nil {
				log.Panic("重置中，删除表失败。。")
			}

		}
		b, err := tx.CreateBucket([]byte(utxoTableName))
		if err != nil {
			log.Panic("重置中，创建新表失败。。")
		}
		if b != nil {
			txOutputMap := utxoSet.BlockChain.FindUnSpentOutputMap()
			//fmt.Println("未花费outputmap：",txOutputMap)
			for txIDStr, outputs := range txOutputMap {
				txID, _ := hex.DecodeString(txIDStr)
				b.Put(txID, outputs.Serilalize())
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (utxoSet *UTXOSet) GetBalance(address string) int64 {
	utxos := utxoSet.FindUnspentOutputsForAddress(address)
	var amount int64
	for _, utxo := range utxos {
		amount += utxo.Output.Value
		fmt.Println(address,"余额：",utxo.Output.Value)
	}
	//fmt.Printf("%s账户，有%d个Token\n",address,amount)
	return amount

}

func (utxoSet *UTXOSet) FindUnspentOutputsForAddress(address string) []*UTXO {
	var utxos [] *UTXO
	//查询数据，遍历所有的未花费
	err := utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				//fmt.Printf("key=%s,value=%v\n", k, v)
				txOutputs := DeserializeTXOutputs(v)
				for _, utxo := range txOutputs.UTXOS {
					if utxo.Output.UnLockWithAddress(address) {
						utxos = append(utxos, utxo)
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return utxos
}

//添加一个方法：用于查询给定地址下的，要转账使用的可以使用的utxo
func (utxoSet *UTXOSet) FindSpendableUTXOs(from string, amount int64, txs [] *Transaction) (int64, map[string][]int) {
	spentableUTXO := make(map[string][]int)
	var total int64 = 0
	//1.找出未打包的Transaction中未花费的
	unPackageUTXOs := utxoSet.FindUnPackageSpentableUTXOs(from, txs)
	for _, utxo := range unPackageUTXOs {
		total += utxo.Output.Value
		txIDStr := hex.EncodeToString(utxo.TxID)
		spentableUTXO[txIDStr] = append(spentableUTXO[txIDStr], utxo.Index)
		fmt.Println(amount,",未打包，转账花费：",utxo.Output.Value)
		if total >= amount {
			return total, spentableUTXO
		}
	}

	//钱不够
	//2.找出已经存在数据库中的未花费的
	err := utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			c := b.Cursor()
		dbLoop:
			for k, v := c.First(); k != nil; k, v = c.Next() {
				txOutpus := DeserializeTXOutputs(v)
				for _, utxo := range txOutpus.UTXOS {
					if utxo.Output.UnLockWithAddress(from) {
						total += utxo.Output.Value
						txIDStr := hex.EncodeToString(utxo.TxID)
						spentableUTXO[txIDStr] = append(spentableUTXO[txIDStr], utxo.Index)
						fmt.Println(amount,",数据库，转账花费：",utxo.Output.Value)
						if total >= amount {
							break dbLoop
						}
					}

				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	if total < amount {
		fmt.Printf("%s,账户余额不足，不能转账。。", from)
		os.Exit(1)
	}
	return total, spentableUTXO

}

func (utxoSet *UTXOSet) FindUnPackageSpentableUTXOs(from string, txs []*Transaction) []*UTXO {
	var unUTXOs []*UTXO
	//存储已经花费
	spentTxOutput := make(map[string][]int)

	for i := len(txs) - 1; i >= 0; i-- {
		unUTXOs = caculate(txs[i], from, spentTxOutput, unUTXOs)
	}
	return unUTXOs
}

//每次创建区块后，更新未花费的表
func (utxoSet *UTXOSet) Update() {
	/*
	每当创建新区块后，都会花掉一些原来的utxo，产生新的utxo。
	删除已经花费的，增加新产生的未花费
	表中存储的数据结构：
	key：交易ID
	value：TxInputs
		TxInputs里是UTXO数组

	 */

	//1.获取最新的区块，由于该block的产生，
	newBlock := utxoSet.BlockChain.Iterator().Next()

	//2.遍历该区块的交易
	inputs := []*TXInput{}
	//spentUTXOMap := make(map[string]*TxOutputs)
	//未花费
	outsMap := make(map[string]*TxOutputs)

	//获取已经花费的
	for _, tx := range newBlock.Txs {
		if tx.IsCoinbaseTransaction(){
			continue
		}
		for _, in := range tx.Vins {
			inputs = append(inputs, in)
		}
	}
	fmt.Println("inputs的长度：", len(inputs), inputs)
	//以上是找出新添加的区块中的所有的Input

	//以下是找到新添加的区块中的未花费了的Output
	for _, tx := range newBlock.Txs {

		utoxs := []*UTXO{}

	outLoop:
		for index, out := range tx.Vouts {
			isSpent := false
			for _, in := range inputs {
				if bytes.Compare(in.TxID, tx.TxID) == 0 && in.Vout == index && bytes.Compare(out.PubKeyHash, PubKeyHash(in.PublicKey)) == 0 {
					isSpent = true
					continue outLoop
				}
			}
			if isSpent == false {
				utxo := &UTXO{tx.TxID, index, out}
				utoxs = append(utoxs, utxo)
				fmt.Println("outsMaps,",out.Value)
			}
		}
		if len(utoxs) > 0 {
			txIDStr := hex.EncodeToString(tx.TxID)
			outsMap[txIDStr] = &TxOutputs{utoxs}
		}

	}
	fmt.Println("outsMap的长度：", len(outsMap), outsMap)
	/*
	outsMap[7788] = utxo
	outsMap[5555] = utxo,utxo
	outsMap[4444] = utxo

	 */

	//删除已经花费了的
	err := utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			//删除 ins中
			for i:=0;i<len(inputs);i++{

				in:=inputs[i]
				fmt.Println(i,"=========================")

				//txIDStr:=hex.EncodeToString(in.TxID)
				//outputs:=outsMap[txIDStr]
				//if len(outputs.UTXOS) > 0{
				//	utxos:=[]*UTXO{}
				//	for _,utxo:=range outputs.UTXOS{
				//		if bytes.Compare(utxo.Output.PubKeyHash, PubKeyHash(in.PublicKey)) == 0 && in.Vout ==utxo.Index {
				//
				//		}else{
				//			utxos = append(utxos,utxo)
				//		}
				//	}
				//
				//	outsMap[txIDStr] = &TxOutputs{utxos}
				//
				//
				//}else{
				txOutputsBytes := b.Get(in.TxID)
				if len(txOutputsBytes) == 0 {
					fmt.Println("break",i)
					continue
				}

				txOutputs := DeserializeTXOutputs(txOutputsBytes)

				// 判断是否需要
				isNeedDelete := false

				utxos := []*UTXO{} //存储未花费
				//根据IxID，如果该txOutputs中已经有output被新区块花掉了，那么将未花掉的添加到utxos里，并标记该txouputs要删除
				for _, utxo := range txOutputs.UTXOS {
					if bytes.Compare(utxo.Output.PubKeyHash, PubKeyHash(in.PublicKey)) == 0 && in.Vout == utxo.Index {
						isNeedDelete = true
					} else {
						utxos = append(utxos, utxo)
					}
				}
				fmt.Println(len(utxos))
				if isNeedDelete == true {
					b.Delete(in.TxID)
					if len(utxos) > 0 {
						//该删除的txoutputs中还有其他未花费的
						//preTXOutputs := outsMap[hex.EncodeToString(in.TxID)]
						//preTXOutputs.UTXOS = append(preTXOutputs.UTXOS,utxos...)
						//outsMap[hex.EncodeToString(in.TxID)] = preTXOutputs

						txOutputs := &TxOutputs{utxos}
						b.Put(in.TxID, txOutputs.Serilalize())
						fmt.Println("删除时：map：", len(outsMap), outsMap)
					}

				}

				//}

			}

			//增加
			for keyID, outPuts := range outsMap {
				keyHashBytes, _ := hex.DecodeString(keyID)
				b.Put(keyHashBytes, outPuts.Serilalize())
			}

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}
