package BLC

import "bytes"

//交易的输出，就是币实际存储的地方
type TXOuput struct {
	Value int64
	//一个锁定脚本(ScriptPubKey)，要花这笔钱，必须要解锁该脚本。
	//ScriptPubKey string //公钥：先理解为，用户名
	PubKeyHash [] byte // 公钥
}

//判断当前txOutput消费，和指定的address是否一致
func (txOutput *TXOuput) UnLockWithAddress(address string) bool {
	//return txOutput.ScriptPubKey == address
	fullPayloadHash := Base58Decode([]byte(address))
	pubKeyHash := fullPayloadHash[1:len(fullPayloadHash)-4]
	return bytes.Compare(txOutput.PubKeyHash, pubKeyHash) == 0
}


func NewTXOuput(value int64,address string) *TXOuput{
	txOutput := &TXOuput{value, nil}
	//设置Ripemd160Hash
	txOutput.Lock(address)
	return txOutput
}

func (txOutput *TXOuput) Lock(address string) {
	publicKeyHash := Base58Decode([] byte(address))
	txOutput.PubKeyHash = publicKeyHash[1:len(publicKeyHash)-4]
}