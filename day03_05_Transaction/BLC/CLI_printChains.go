package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI)printChains(){
	bc:=GetBlockchainObject()
	if bc == nil{
		fmt.Println("没有区块可以打印。。")
		os.Exit(1)
	}
	defer bc.DB.Close()
	bc.PrintChains()
}