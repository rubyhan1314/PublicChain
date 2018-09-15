package BLC

import (
	"fmt"
	"os"
)

//转账
func (cli *CLI) send(from, to, amount [] string) {
	if !dbExists() {
		fmt.Println("数据库不存在。。。")
		os.Exit(1)
	}
	blockchain := GetBlockchainObject()

	blockchain.MineNewBlock(from, to, amount)
	defer blockchain.DB.Close()
}
