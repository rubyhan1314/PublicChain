package BLC

import "bytes"

type TXInput struct {
	//1.交易的ID
	TxID [] byte
	//2.存储Txoutput的vout里面的索引
	Vout int
	//3.用户名
	//ScriptSiq string
	Signature [] byte //数字签名
	PublicKey [] byte //公钥，钱包里面
}

//判断当前txInput消费，和指定的address是否一致
func (txInput *TXInput) UnLockWithAddress(pubKeyHash []byte) bool{
	//return txInput.ScriptSiq == address
	publicKey:=PubKeyHash(txInput.PublicKey)
	return bytes.Compare(pubKeyHash,publicKey) == 0
}


