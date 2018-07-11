package blockchain

import (
	//"strconv"
	//"bytes"
	//"crypto/sha256"
	"time"
)

/**
	区块
 */
type Block struct {
	Timestamp int64 //时间戳
	PreBlockHash []byte	//上一个区块哈希值
	Hash []byte	//当前区块哈希值
	Data []byte	//本区块信息
	nonce int //随机值
}

///**
//	区块-计算自身哈希值
// */
//func (b *Block)setHash(){
//	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))//获取当前时间戳
//	headers := bytes.Join([][]byte{b.PreBlockHash, b.Data, timestamp}, []byte{})//拼接上区块哈希、本区块数据、时间戳，分隔符为空
//	hash := sha256.Sum256(headers)//生成哈希
//	b.Hash = hash[:]//给bean的hash赋值 byte[low_index:high_index] 数组切片，可以low和high省略
//
//}
/**
	创建区块
 */
func NewBlock(data string, prevBlockHash []byte) *Block {
	//生成时间戳
	timestamp := time.Now().Unix()
	//创建区块
	block := &Block{timestamp, prevBlockHash, []byte{}, []byte(data),0}
	pow := NewPow(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.nonce = nonce
	return block
}

/**
	创建创世区块
 */
func NewGenesisBlock() *Block {
	return NewBlock("我是创世0区块",[]byte{})
}