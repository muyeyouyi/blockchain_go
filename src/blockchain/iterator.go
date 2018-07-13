package blockchain

import (
	"../constants"
	"github.com/boltdb/bolt"
)

/**
迭代器
*/
type ChainIterator struct {
	CurrentHash []byte
	db          *bolt.DB
}

/**
获得当前区块，并将上一个区块的hash存到结构里
*/
func (iterator *ChainIterator) Next() *Block {
	var block *Block
	iterator.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(constants.BlocksBucket))
		block = Deserialize(b.Get(iterator.CurrentHash))
		iterator.CurrentHash = block.PreBlockHash
		return nil
	})
	return block
}
