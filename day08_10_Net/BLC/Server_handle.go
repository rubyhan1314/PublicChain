package BLC

import (
	"bytes"
	"log"
	"encoding/gob"
	"fmt"
	"encoding/hex"
	"github.com/boltdb/bolt"
)

func handleVersion(request []byte, bc *BlockChain) {

	var buff bytes.Buffer
	var payload Version

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	//Version
	//1. Version
	//2. BestHeight
	//3. 节点地址

	bestHeight := bc.GetBestHeight()          //3 1
	foreignerBestHeight := payload.BestHeight // 1 3

	if bestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	} else if bestHeight < foreignerBestHeight {
		// 去向主节点要信息
		sendGetBlocks(payload.AddrFrom)
	}

	if !nodeIsKnown(payload.AddrFrom) {
		knowNodes = append(knowNodes, payload.AddrFrom)
	}

}

func handleAddr(request []byte, bc *BlockChain) {

}

func handleGetblocks(request []byte, bc *BlockChain) {

	var buff bytes.Buffer
	var payload GetBlocks

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()

	//txHash blockHash
	sendInv(payload.AddrFrom, BLOCK_TYPE, blocks)

}

func handleInv(request []byte, bc *BlockChain) {

	var buff bytes.Buffer
	var payload Inv

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	// Ivn 3000 block hashes [][]

	if payload.Type == BLOCK_TYPE {

		//tansactionArray = payload.Items

		//payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, BLOCK_TYPE, blockHash)

		if len(payload.Items) >= 1 {
			transactionArray = payload.Items[1:]
		}
	}

	if payload.Type == TX_TYPE {

		txHash := payload.Items[0]
		if memoryTxPool[hex.EncodeToString(txHash)] == nil {
			sendGetData(payload.AddrFrom, TX_TYPE, txHash)
		}

	}

}

func handleGetData(request []byte, bc *BlockChain) {

	var buff bytes.Buffer
	var payload GetData

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == BLOCK_TYPE {

		block, err := bc.GetBlock([]byte(payload.Hash))
		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, block)
	}

	if payload.Type == TX_TYPE {

		tx := memoryTxPool[hex.EncodeToString(payload.Hash)]

		sendTx(payload.AddrFrom, tx)

	}
}

func handleBlock(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payload BlockData

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockBytes := payload.Block

	block := DeserializeBlock(blockBytes)

	fmt.Println("Recevied a new block!")
	bc.AddBlock(block)
	UTXOSet := &UTXOSet{bc}
	UTXOSet.Update()

	fmt.Printf("Added block %x\n", block.Hash)

	if len(transactionArray) > 0 {
		blockHash := transactionArray[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		transactionArray = transactionArray[1:]
	} else {

		//fmt.Println("数据库重置......")
		//UTXOSet := &UTXOSet{bc}
		//UTXOSet.ResetUTXOSet()

	}

}

func handleTx(request []byte, bc *BlockChain) {

	var buff bytes.Buffer
	var payload Tx

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	//-----

	tx := payload.Tx
	memoryTxPool[hex.EncodeToString(tx.TxID)] = tx

	// 说明主节点自己
	if nodeAddress == knowNodes[0] {
		// 给矿工节点发送交易hash
		for _, nodeAddr := range knowNodes {

			if nodeAddr != nodeAddress && nodeAddr != payload.AddrFrom {
				sendInv(nodeAddr, TX_TYPE, [][]byte{tx.TxID})
			}

		}
	}

	// 矿工进行挖矿验证
	// "" | 1DVFvyCK8qTQkLBTZ5fkh5eDSbcZVoHAsj
	if len(memoryTxPool) >= 1 && len(minerAddress) > 0 {

	MineTransactions:

		utxoSet := &UTXOSet{bc}

		txs := []*Transaction{tx}

		//奖励
		coinbaseTx := NewCoinBaseTransaction(minerAddress)
		txs = append(txs, coinbaseTx)

		_txs := []*Transaction{}

		//fmt.Println("开始进行数字签名验证.....")

		for _, tx := range txs {

			//fmt.Printf("开始第%d次验证...\n",index)

			// 数字签名失败
			if bc.VerifyTransaction(tx, _txs) != true {
				log.Panic("ERROR: Invalid transaction")
			}

			//fmt.Printf("第%d次验证成功\n",index)
			_txs = append(_txs, tx)
		}

		//fmt.Println("数字签名验证成功.....")

		//1. 通过相关算法建立Transaction数组
		var block *Block

		bc.DB.View(func(tx *bolt.Tx) error {

			b := tx.Bucket([]byte(BLOCKTABLENAME))
			if b != nil {

				hash := b.Get([]byte("l"))

				blockBytes := b.Get(hash)

				block = DeserializeBlock(blockBytes)

			}

			return nil
		})

		//2. 建立新的区块
		block = NewBlock(txs, block.Hash, block.Height+1)

		//将新区块存储到数据库
		bc.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(BLOCKTABLENAME))
			if b != nil {

				b.Put(block.Hash, block.Serilalize())

				b.Put([]byte("l"), block.Hash)

				bc.Tip = block.Hash

			}
			return nil
		})
		utxoSet.Update()
		sendBlock(knowNodes[0], block.Serilalize())

		for _, tx := range txs {
			txID := hex.EncodeToString(tx.TxID)
			delete(memoryTxPool, txID)
		}
		for _, node := range knowNodes {
			if node != nodeAddress {
				sendInv(node, "block", [][]byte{block.Hash})
			}
		}

		if len(memoryTxPool) > 0 {
			goto MineTransactions
		}
	}
}
