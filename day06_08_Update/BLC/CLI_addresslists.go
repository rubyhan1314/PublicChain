package BLC

import "fmt"

func (cli *CLI)addressLists(){
	fmt.Println("打印所有的钱包地址。。")
	//获取
	Wallets:=NewWallets()
	for address,_ := range Wallets.WalletsMap{
		fmt.Println("address:",address)
	}
}
