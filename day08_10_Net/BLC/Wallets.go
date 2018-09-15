package BLC

import (
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
)

//1.创建钱包
type Wallets struct {
	WalletsMap map[string]*Wallet
}

//2.创建一个钱包集合
//创建钱包集合:文件中存在从文件中读取，否则新建一个
const walletFile = "Wallets_%s.dat"

func NewWallets(nodeID string) *Wallets {
	//wallets := &WalletsMap{}
	//wallets.WalletsMap = make(map[string]*Wallet)
	//return wallets

	walletFile := fmt.Sprintf(walletFile,nodeID)


	//1.判断钱包文件是否存在
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		fmt.Println("文件不存在")
		wallets := &Wallets{}
		wallets.WalletsMap = make(map[string]*Wallet)
		return wallets
	}
	//2.否则读取文件中的数据
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	return &wallets
}

//3.创建一个新钱包
func (ws *Wallets) CreateNewWallet(nodeID string) {
	wallet := NewWallet()
	fmt.Printf("创建钱包地址：%s\n", wallet.GetAddress())
	ws.WalletsMap[string(wallet.GetAddress())] = wallet

	//将钱包保存
	ws.SaveWallets(nodeID)
}

/*
要让数据对象能在网络上传输或存储，我们需要进行编码和解码。现在比较流行的编码方式有JSON,XML等。然而，Go在gob包中为我们提供了另一种方式，该方式编解码效率高于JSON。
gob是Golang包自带的一个数据结构序列化的编码/解码工具
 */
func (ws *Wallets) SaveWallets(nodeID string) {


	walletFile := fmt.Sprintf(walletFile,nodeID)


	var content bytes.Buffer
	//注册的目的，为了可以序列化任何类型，wallet结构体中有接口类型。将接口进行注册
	gob.Register(elliptic.P256()) //gob是Golang包自带的一个数据结构序列化的编码/解码工具
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	//将序列化后的数据写入到文件，原来的文件中的内容会被覆盖掉
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
