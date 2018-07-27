package blockchain

import (
	"crypto/sha256"
)


/**
	一笔交易模型
 */
type Transaction struct {
	Id      []byte     //自身的id值
	Inputs  []TxInput  //交易来源集合
	Outputs []TxOutput //交易输出集合
}

/**
	判断挖矿交易
 */
func (transaction *Transaction) IsCoinBase() bool {
	if transaction.Inputs[0].OutIndex == -1 {
		return true
	}
	return false
}

/**
	生成交易的ID
 */
func (transaction *Transaction) setId() {
	tx := Serialize(transaction)
	txHash := sha256.Sum256(tx)
	transaction.Id = txHash[:]

}

/*/
	创建一个output
 */
func NewTxOutput(value int, address string) *TxOutput {
	output := &TxOutput{value, nil}
	output.Lock([]byte(address))
	return output
}




