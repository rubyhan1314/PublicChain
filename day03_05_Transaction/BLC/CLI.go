package BLC

import (
	"os"
	"fmt"
	"flag"
	"log"
)

//step1:
//CLI结构体
type CLI struct {
	//Blockchain *BlockChain
}

//step2：添加Run方法
func (cli *CLI) Run(){
	//判断命令行参数的长度
	isValidArgs()

	//1.创建flagset标签对象
	sendBlockCmd := flag.NewFlagSet("send",flag.ExitOnError)
	//fmt.Printf("%T\n",addBlockCmd) //*FlagSet
	printChainCmd:=flag.NewFlagSet("printchain",flag.ExitOnError)
	createBlockChainCmd:=flag.NewFlagSet("createblockchain",flag.ExitOnError)
	getBalanceCmd:=flag.NewFlagSet("getbalance",flag.ExitOnError)

	//2.设置标签后的参数
	//flagAddBlockData:= addBlockCmd.String("data","helloworld..","交易数据")
	flagFromData:=sendBlockCmd.String("from","","转帐源地址")
	flagToData:=sendBlockCmd.String("to","","转帐目标地址")
	flagAmountData:=sendBlockCmd.String("amount","","转帐金额")
	flagCreateBlockChainData := createBlockChainCmd.String("address","","创世区块交易地址")
	flagGetBalanceData := getBalanceCmd.String("address","","要查询的某个账户的余额")


	//3.解析
	switch os.Args[1] {
	case "send":
		err:=sendBlockCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
		//fmt.Println("----",os.Args[2:])

	case "printchain":
		err :=printChainCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
		//fmt.Println("====",os.Args[2:])


	case "createblockchain":
		err :=createBlockChainCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "getbalance":
		err :=getBalanceCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}


	default:
		printUsage()
		os.Exit(1)//退出
	}

	if sendBlockCmd.Parsed(){
		if *flagFromData == "" || *flagToData =="" ||*flagAmountData == "" {
			printUsage()
			os.Exit(1)
		}
		//cli.addBlock([]*Transaction{})
		fmt.Println(*flagFromData)
		fmt.Println(*flagToData)
		fmt.Println(*flagAmountData)
		//fmt.Println(JSONToArray(*flagFrom))
		//fmt.Println(JSONToArray(*flagTo))
		//fmt.Println(JSONToArray(*flagAmount))
		from:=JSONToArray(*flagFromData)
		to:=JSONToArray(*flagToData)
		amount:=JSONToArray(*flagAmountData)

		cli.send(from,to,amount)
	}
	if printChainCmd.Parsed(){
		cli.printChains()
	}

	if createBlockChainCmd.Parsed(){
		if *flagCreateBlockChainData == ""{
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*flagCreateBlockChainData)
	}

	if getBalanceCmd.Parsed(){
		if *flagGetBalanceData == ""{
			fmt.Println("查询地址不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceData)

	}

}

func isValidArgs(){
	if len(os.Args) < 2{
		printUsage()
		os.Exit(1)
	}
}
func printUsage(){
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -address DATA -- 创建创世区块")
	fmt.Println("\tsend -from From -to To -amount Amount - 交易数据")
	fmt.Println("\tprintchain - 输出信息")
	fmt.Println("\tgetbalance -address DATA -- 查询账户余额")
}








