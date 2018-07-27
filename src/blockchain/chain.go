package blockchain

import (
	"constants"
	"utils"
	"github.com/boltdb/bolt"
	"encoding/hex"
	"bytes"
	"errors"
	"crypto/ecdsa"
)

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
func (bc *Chain) AddBlock(txs []*Transaction) {
	for _, tx := range txs {
		if bc.VerifyTransaction(tx) != true {
			utils.LogD("交易验证失败！")
			return
		}
	}

	//数据库查询最后区块哈希值
	var preHash []byte
	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(constants.BlocksBucket))
		block := Deserialize(b.Get(bc.Tip))
		preHash = block.Hash
		return nil
	})
	//创建新区块
	newBlock := NewBlock(txs, preHash)
	//新区快写入数据库
	bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(constants.BlocksBucket))
		b.Put(newBlock.Hash, Serialize(newBlock))
		b.Put(constants.LastHash, newBlock.Hash)
		bc.Tip = newBlock.Hash
		return nil
	})

	utils.LogE(err)

	//preHash := bc.Blocks[len(bc.Blocks)-1].Hash//上个区块哈希值
	//newBlock := NewBlock(data, preHash)//创建新区快
	//bc.Blocks = append(bc.Blocks, newBlock)
}

/**
	链
*/
type Chain struct {
	Tip []byte //最后一个区块的哈希
	Db  *bolt.DB
}

/**
	创建区块链
*/
func NewBlockChain(transaction *Transaction) *Chain {
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
			block := NewGenesisBlock(transaction)
			b.Put(block.Hash, Serialize(block))
			b.Put(constants.LastHash, block.Hash)
			tip = block.Hash
			return err
			//存在，更新表
		} else {
			tip = b.Get(constants.LastHash)
		}
		return nil
	})
	utils.LogE(err)

	chain := Chain{tip, db}
	return &chain
}

func (bc *Chain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	//utxo
	var unspentTxs []Transaction
	//已花费的out索引集合{{tXid1:[index1,index2]},{tXid2:[index1,index2]}}
	spentOutputIndex := make(map[string][]int)

	iterator := bc.Iterator()
	//一层循环拿到区块
	for {
		block := iterator.Next()
		//二层循环拿到区块中的tx
		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.Id)
		Output:
		//三层循环之一，拿到output
			for outIndex, output := range tx.Outputs {
				//找到对应in的数组，查找其中相同项
				if spentOutputIndex[txId] != nil {
					for _, vOut := range spentOutputIndex[txId] {
						if vOut == outIndex {
							continue Output
						}
					}
				}
				//未找到对应in，该out存到utxo
				if output.CanBeUnlockedWith(pubKeyHash) {
					unspentTxs = append(unspentTxs, *tx)
				}

			}
			//三层循环之二，拿到input
			if !tx.IsCoinBase() {
				for _, input := range tx.Inputs {
					//所有in的数据写到stx集合
					if input.CanUnlockOutputWith(pubKeyHash) {
						inTxId := hex.EncodeToString(input.TxId)
						spentOutputIndex[inTxId] = append(spentOutputIndex[inTxId], input.outIndex)
					}
				}
			}

		}

		if len(iterator.CurrentHash) == 0 {
			break
		}
	}
	return unspentTxs
}

/**
	找到对应地址utxo的集合
 */
func (bc *Chain) FindUTXOs(pubKeyHash []byte) []TxOutput {
	var utxos []TxOutput
	txs := bc.FindUnspentTransactions(pubKeyHash)
	for _, tx := range txs {
		for _, output := range tx.Outputs {
			if output.CanBeUnlockedWith(pubKeyHash) {
				utxos = append(utxos, output)
			}
		}
	}
	return utxos
}

/**
	找到这一笔交易够用的output
 */
func (bc *Chain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	var vaildOutputs = make(map[string][]int)
	accumulated := 0
	txs := bc.FindUnspentTransactions(pubKeyHash)

Work:
	for _, tx := range txs {
		txId := hex.EncodeToString(tx.Id)
		for index, output := range tx.Outputs {
			if output.CanBeUnlockedWith(pubKeyHash) {
				accumulated += output.ValueOut
				vaildOutputs[txId] = append(vaildOutputs[txId], index)
				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, vaildOutputs

}

/**
	根据id在链上找到一笔交易
 */
func (bc *Chain) FindTransaction(ID []byte) (Transaction, error) {
	iterator := bc.Iterator()
	for{
		block := iterator.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(ID,tx.Id) == 0 {
				return *tx,nil
			}
		}


		if len(block.PreBlockHash) == 0 {
			break
		}
	}
	return Transaction{},errors.New("没有找到这笔交易")
}

/**
	对交易签名
 */
func (bc *Chain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey) {
	preventTxs := make(map[string]Transaction)
	for _, input := range tx.Inputs {
		preTx, e := bc.FindTransaction(input.TxId)
		utils.LogE(e)
		preventTxs[hex.EncodeToString(preTx.Id)] = preTx
	}
	tx.Sign(privateKey,preventTxs)
}

/**
	验证交易签名
 */
func (bc *Chain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Inputs {
		prevTX, err := bc.FindTransaction(vin.TxId)
		prevTXs[hex.EncodeToString(prevTX.Id)] = prevTX
		utils.LogE(err)
	}

	return tx.Verify(prevTXs)
}
