package BLC


func (cli *CLI) createGenesisBlockchain(address string,nodeID string){
	//fmt.Println(data)
	CreateBlockChainWithGenesisBlock(address,nodeID)


	bc:=GetBlockchainObject(nodeID)
	defer bc.DB.Close()
	if bc != nil{
		utxoSet:=&UTXOSet{bc}
		utxoSet.ResetUTXOSet()
	}

}


