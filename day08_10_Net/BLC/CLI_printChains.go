package BLC

import (
	"fmt"
	"os"
)



func (cli *CLI)printChains(nodeID string){
	bc:=GetBlockchainObject(nodeID)
	if bc == nil{
		fmt.Println("没有区块可以打印。。")
		os.Exit(1)
	}
	defer bc.DB.Close()
	bc.PrintChains()
}