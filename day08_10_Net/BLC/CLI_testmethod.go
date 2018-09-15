package BLC

import "fmt"


func (cli *CLI) TestMethod(nodeID string){
	blockchain:=GetBlockchainObject(nodeID)
	//defer blockchain.DB.Close()

	unSpentOutputMap:=blockchain.FindUnSpentOutputMap()
	fmt.Println(unSpentOutputMap)
	for key,value:=range unSpentOutputMap{
		fmt.Println(key)
		for _,utxo:=range value.UTXOS{
			fmt.Println("金额：",utxo.Output.Value)
			fmt.Printf("地址：%v\n",utxo.Output.PubKeyHash)
			fmt.Println("---------------------")
		}
	}

	utxoSet:=&UTXOSet{blockchain}
	utxoSet.ResetUTXOSet()
}





