package BLC

type TXInput struct {
	//1.交易的ID
	TxID [] byte
	//2.存储Txoutput的vout里面的索引
	Vout int
	//3.用户名
	ScriptSiq string
}

//判断当前txInput消费，和指定的address是否一致
func (txInput *TXInput) UnLockWithAddress(address string) bool{
	return txInput.ScriptSiq == address
}


