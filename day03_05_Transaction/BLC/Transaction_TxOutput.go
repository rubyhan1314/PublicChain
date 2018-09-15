package BLC
//交易的输出，就是币实际存储的地方
type TXOuput struct {
	Value        int64
	//一个锁定脚本(ScriptPubKey)，要花这笔钱，必须要解锁该脚本。
	ScriptPubKey string //公钥：先理解为，用户名
}

//判断当前txOutput消费，和指定的address是否一致
func (txOutput *TXOuput) UnLockWithAddress(address string) bool{
	return txOutput.ScriptPubKey == address
}