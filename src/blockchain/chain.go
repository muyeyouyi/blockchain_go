package blockchain

import (
	"constants"
	"utils"
	"github.com/boltdb/bolt"
)

var lastHash = []byte("last_hash") //最后一个块的hash值

/**
创建区块链
*/
//func CreateBlockChain() *Chain {
//	return &Chain{[]*Block{NewGenesisBlock()}}
//}

/**
迭代器
*/
func (bc *Chain) Iterator() *ChainIterator {
	iterator := &ChainIterator{bc.Tip, bc.Db}
	return iterator
}

/**
区块链-增加区块
*/
func (bc *Chain) AddBlock(data string) {
	//数据库查询最后区块哈希值
	var preHash []byte
	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(constants.BlocksBucket))
		block := Deserialize(b.Get(bc.Tip))
		preHash = block.Hash
		return nil
	})
	//创建新区块
	newBlock := NewBlock(data, preHash)
	//新区快写入数据库
	bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(constants.BlocksBucket))
		b.Put(newBlock.Hash, Serialize(newBlock))
		b.Put(lastHash, newBlock.Hash)
		bc.Tip = newBlock.Hash
		return nil
	})

	utils.LogE(err)

	//preHash := bc.Blocks[len(bc.Blocks)-1].Hash//上个区块哈希值
	//newBlock := NewBlock(data, preHash)//创建新区快
	//bc.Blocks = append(bc.Blocks, newBlock)
}

/**
区块链
*/
type Chain struct {
	Tip []byte
	Db  *bolt.DB
}

/**
创建区块链
*/
func NewBlockChain() *Chain {
	utils.CreateFile(constants.DbFile)
	//打开数据库
	var tip []byte
	db, err := bolt.Open(constants.DbFile, 0600, nil)
	//数据库写入标准
	err = db.Update(func(tx *bolt.Tx) error {
		//检查表是否存在
		b := tx.Bucket([]byte(constants.BlocksBucket))
		//不存在，创建表
		if b == nil {
			b, err := tx.CreateBucket([]byte(constants.BlocksBucket))
			block := NewGenesisBlock()
			b.Put(block.Hash, Serialize(block))
			b.Put(lastHash, block.Hash)
			tip = block.Hash
			return err
			//存在，更新表
		} else {
			tip = b.Get(lastHash)
		}
		return nil
	})
	utils.LogE(err)

	chain := Chain{tip, db}
	return &chain
}
