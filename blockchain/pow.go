package blockchain

import (
	"math/big"
	"bytes"
	"../utils"
	"math"
	"crypto/sha256"
	"fmt"
)

const targetBits = 10//目标难度的指数

type ProofOfWork struct {
	block  *Block
	target *big.Int //目标难度
}

/**
	拼装区块内所有数据
 */
func (pow *ProofOfWork) PrepareData(nonce int) []byte {
	//拼接哈希数据
	data := bytes.Join(
		[][]byte{
			pow.block.Data,
			pow.block.PreBlockHash,
			utils.IntToBytes(pow.block.Timestamp),
			utils.IntToBytes(int64(nonce)),//随机值
		},
		[]byte{'-'})
	return data
}

/**
	构造方法
 */
func NewPow(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, 256-targetBits)
	pow := &ProofOfWork{block, target}
	return pow
}

/**
	挖矿
 */
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	fmt.Printf("开始挖矿，目标难度:%d\n", int(math.Pow(2, targetBits)))
	for ; nonce < math.MaxInt64; nonce++ {
		data := pow.PrepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("hash:%x\n", hash)
			fmt.Printf("挖矿次数:%d\n\n", nonce)
			break
		}
	}
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.PrepareData(pow.block.nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid  := hashInt.Cmp(pow.target) == -1
	return isValid
}
