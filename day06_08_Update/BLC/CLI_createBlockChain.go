package BLC

func (cli *CLI) createGenesisBlockchain(address string){
	//fmt.Println(data)
	CreateBlockChainWithGenesisBlock(address)


	bc:=GetBlockchainObject()
	defer bc.DB.Close()
	if bc != nil{
		utxoSet:=&UTXOSet{bc}
		utxoSet.ResetUTXOSet()
	}

}
