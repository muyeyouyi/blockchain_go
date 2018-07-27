package blockchain

import (
	"bytes"
	"utils"
	"wallet"
)

/**
	交易输出
 */
type TxOutput struct {
	ValueOut   int    //数量
	PubKeyHash []byte //公钥哈希值
}

/**
	之前交易得到的币能否被该公钥解锁
 */
func (out *TxOutput) CanBeUnlockedWith(pubKeyHash []byte) bool {
	return bytes.Compare(pubKeyHash,out.PubKeyHash) == 0
}

/**
	用地址锁定输出
 */
func (out *TxOutput) Lock(address []byte){
	fullHash := utils.Base58Decode(address)
	pubKeyHash := fullHash[1:len(fullHash)-wallet.AddressChecksumLen]
	out.PubKeyHash = pubKeyHash
}