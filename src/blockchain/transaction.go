package blockchain

import (
	"crypto/sha256"
	"crypto/ecdsa"
	"encoding/hex"
	"crypto/rand"
	"utils"
	"crypto/elliptic"
	"math/big"
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
func (tx *Transaction) IsCoinBase() bool {
	if tx.Inputs[0].outIndex == -1 {
		return true
	}
	return false
}

/**
	生成交易的ID
 */
func (tx *Transaction) setId() {
	serTx := Serialize(tx)
	txHash := sha256.Sum256(serTx)
	tx.Id = txHash[:]
}

/**
	签名实现细节
 */
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, preventTxs map[string]Transaction) {
	//挖矿交易不签名
	if tx.IsCoinBase() {
		return
	}

	//复制一份tx样本
	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Inputs {
		//取到对应前一个out公钥哈希
		prevTx := preventTxs[hex.EncodeToString(vin.TxId)]
		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].PubKey = prevTx.Outputs[vin.outIndex].PubKeyHash
		//计算txid
		txCopy.Id = txCopy.Hash()
		//公钥置空
		txCopy.Inputs[inID].PubKey = nil

		//签名后写入真实tx
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.Id)
		signature := append(r.Bytes(), s.Bytes()...)
		utils.LogE(err)
		tx.Inputs[inID].Signature = signature
	}
}

/**
	拷贝一份tx用于签名
 */
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, vin := range tx.Inputs {
		inputs = append(inputs, TxInput{vin.TxId, vin.outIndex, nil, nil})
	}

	for _, vout := range tx.Outputs {
		outputs = append(outputs, TxOutput{vout.ValueOut, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.Id, inputs, outputs}
	return txCopy
}

/**
	验证tx签名
 */
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Inputs {
		prevTx := prevTXs[hex.EncodeToString(vin.TxId)]
		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].PubKey = prevTx.Outputs[vin.outIndex].PubKeyHash
		txCopy.Id = txCopy.Hash()
		txCopy.Inputs[inID].PubKey = nil

		//拆分签名文件
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		//拆分公钥
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])
		//还原为原始公钥
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		//公钥、签名文件、原始数据确认签名有效性
		if ecdsa.Verify(&rawPubKey, txCopy.Id, &r, &s) == false {
			return false
		}
	}

	return true
}

/**
	对交易哈希
 */
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.Id = []byte{}

	hash = sha256.Sum256(Serialize(txCopy))

	return hash[:]
}

/*/
	创建一个output
 */
func NewTxOutput(value int, address string) *TxOutput {
	output := &TxOutput{value, nil}
	output.Lock([]byte(address))
	return output
}




