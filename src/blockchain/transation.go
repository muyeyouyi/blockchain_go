package blockchain

import (
	"constants"
	"crypto/sha256"
)

/**
	交易输入
 */
type TxInput struct {
	TxId      []byte //对应out所在tx的id值
	OutIndex  int    //对应out所在tx的索引
	ScriptSig string //脚本，现版本存地址值
}

/**
	交易输出
 */
type TxOutput struct {
	ValueOut  int    //数量
	ScriptSig string //脚本，现版本存地址值
}

/**
	一笔交易模型
 */
type Transaction struct {
	Id      []byte     //自身的id值
	Inputs  []TxInput  //交易来源集合
	Outputs []TxOutput //交易输出集合
}

/**
	生成交易的ID
 */
func (transaction *Transaction) setId() {
	tx := Serialize(transaction)
	txHash := sha256.Sum256(tx)
	transaction.Id = txHash[:]


}

/**
	创建挖矿交易
 */
func NewCoinBaseTx (toAddress string) *Transaction{
	inputs := TxInput{[]byte{},-1,""}
	outputs := TxOutput{constants.Subsidy,toAddress}
	tx := Transaction{nil,[]TxInput{inputs},[]TxOutput{outputs}}
	tx.setId()
	return &tx
}
