package blockchain

import (
	"time"
	"bytes"
	"crypto/sha256"
	"constants"
)

/**
区块
*/
type Block struct {
	Timestamp    int64          //时间戳
	PreBlockHash []byte         //上一个区块哈希值
	Hash         []byte         //当前区块哈希值
	Transactions []*Transaction //本区块信息
	Nonce        int            //随机值
}

/**
	拼装tx的id，返回哈希
 */
func (block *Block) HashTransactions() []byte {
	var txIds [][]byte
	for _, tx := range block.Transactions {
		txIds = append(txIds,tx.Id)
	}
	ids := bytes.Join(txIds, []byte{})
	hash := sha256.Sum256(ids)
	return hash[:]
}

/**
创建区块
*/
func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	//生成时间戳
	timestamp := time.Now().Unix()
	//创建区块
	block := &Block{timestamp, prevBlockHash, []byte{}, txs, 0}
	pow := NewPow(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

/**
	创建挖矿交易
 */
func NewCoinBaseTx(toAddress string) *Transaction {
	inputs := TxInput{[]byte{}, -1, []byte{},[]byte("创世区块！！")}
	output := NewTxOutput(constants.Subsidy,toAddress)
	tx := Transaction{nil, []TxInput{inputs}, []TxOutput{*output}}
	tx.setId()
	return &tx
}


/**
创建创世区块
*/
func NewGenesisBlock(transaction *Transaction) *Block {
	return NewBlock([]*Transaction{transaction}, []byte{})
}
