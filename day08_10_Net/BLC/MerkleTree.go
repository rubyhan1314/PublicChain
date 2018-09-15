package BLC

import (
	"crypto/sha256"
	"math"
)

/*
第一步：创建结构体对象，表示节点和树
 */
type MerkleNode struct {
	LeftNode  *MerkleNode
	RightNode *MerkleNode
	Data      []byte
}

type MerkleTree struct {
	RootNode *MerkleNode
}

/*
第二步：给一个左右节点，生成一个新的节点
 */
func NewMerkleNode(leftNode, rightNode *MerkleNode, txHash []byte) *MerkleNode {
	//1.创建当前的节点
	mNode := &MerkleNode{}

	//2.赋值
	if leftNode == nil && rightNode == nil {
		//mNode就是个叶子节点
		hash := sha256.Sum256(txHash)
		mNode.Data = hash[:]
	} else {
		//mNOde是非叶子节点
		prevHash := append(leftNode.Data, rightNode.Data...)
		hash := sha256.Sum256(prevHash)
		mNode.Data = hash[:]
	}
	mNode.LeftNode = leftNode
	mNode.RightNode = rightNode
	return mNode

}

/*
第三步：生成merkle
 */
func NewMerkleTree(txHashData [][]byte) *MerkleTree {
	/*
	Tx1,Tx2,Tx3
	{
		{tx1hash},
		{tx2hash},
		{tx3hash},
		{tx3hash}
	}
	 */

	//1.创建一个数组，用于存储node节点
	var nodes []*MerkleNode

	//2.判断交易量的奇偶性
	if len(txHashData)%2 != 0 {
		//奇数，复制最后一个
		txHashData = append(txHashData, txHashData[len(txHashData)-1])
	}
	//3.创建一排的叶子节点
	for _, datum := range txHashData {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, node)
	}

	//4.生成树其他的节点：6
	//for i:=0;i<len(txHashData)/2;i++{ //  i < 8
	//for {

	count := GetCircleCount(len(nodes))

	for i := 0; i < count; i++ {
		var newLevel []*MerkleNode

		for j := 0; j < len(nodes); j += 2 { //j=0  tx12 tx33
			node := NewMerkleNode(nodes[j], nodes[j+1], nil)
			newLevel = append(newLevel, node)

		}

		//if len(newLevel) == 1 {
		//	break
		//}
		//判断newLevel的长度的奇偶性
		if len(newLevel)%2 != 0 {
			newLevel = append(newLevel, newLevel[len(newLevel)-1])
		}

		nodes = newLevel // 3

	}

	mTree := &MerkleTree{nodes[0]}

	return mTree

}

//获取产生Merkle树根需要循环的次数
func GetCircleCount(len int) int {
	count := 0
	for {
		if int(math.Pow(2, float64(count))) >= len {
			return count
		}
		count++
	}
}
