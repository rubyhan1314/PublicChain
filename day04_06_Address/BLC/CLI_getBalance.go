package BLC

import (
	"fmt"
	"os"
)

//查询余额
func (cli *CLI)getBalance(address string){
	fmt.Println("查询余额：",address)
	bc := GetBlockchainObject()

	if bc == nil{
		fmt.Println("数据库不存在，无法查询。。")
		os.Exit(1)
	}
	defer bc.DB.Close()
	//txOutputs:= bc.UnUTXOs(address)
	//for i,out:=range txOutputs{
	//	fmt.Println(i,"---->",out)
	//}
	balance:=bc.GetBalance(address,[]*Transaction{})
	fmt.Printf("%s,一共有%d个Token\n",address,balance)
}
