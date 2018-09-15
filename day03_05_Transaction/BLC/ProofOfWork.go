package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

// 0000 0000 0000 0000 1001 0001 0000  .... 0001
//256位Hash里面前面至少有16个零
const TargetBit = 16 // 20

type ProofOfWork struct {
	//要验证的区块
	Block *Block

	//大整数存储,目标哈希
	Target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	//1.创建一个big对象 0000000.....00001
	/*
	0000 0001
	0010 0000
	 */
	target := big.NewInt(1)

	//2.左移256-bits位
	target = target.Lsh(target, 256-TargetBit)

	return &ProofOfWork{block, target}
}

func (pow *ProofOfWork) Run() ([] byte, int64) {
	//1.将Block的属性拼接成字节数组
	//2.生成Hash
	//3.循环判断Hash的有效性，满足条件，跳出循环结束验证
	nonce := 0
	//var hashInt big.Int //用于存储新生成的hash
	hashInt := new(big.Int)
	var hash [32]byte
	for{
		//获取字节数组
		dataBytes := pow.prepareData(nonce)
		//生成hash
		hash = sha256.Sum256(dataBytes)
		//fmt.Printf("%d: %x\n",nonce,hash)
		fmt.Printf("\r%d: %x",nonce,hash)
		//将hash存储到hashInt
		hashInt.SetBytes(hash[:])
		//判断hashInt是否小于Block里的target
		/*
		Com compares x and y and returns:
		-1 if x < y
		0 if x == y
		1 if x > y
		 */
		 if pow.Target.Cmp(hashInt) == 1{
		 	break
		 }
		 nonce++
	}
	fmt.Println()
	return hash[:], int64(nonce)
}

func (pow *ProofOfWork) prepareData(nonce int)[]byte{
	data := bytes.Join(
		[][] byte{
			pow.Block.PrevBlockHash,
			pow.Block.HashTransactions(),//此行代码改变
			IntToHex(pow.Block.TimeStamp),
			IntToHex(int64(TargetBit)),
			IntToHex(int64(nonce)),
			IntToHex(int64(pow.Block.Height)),
		},
		[] byte{},
	)
	return data
}
