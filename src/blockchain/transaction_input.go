package blockchain

import (
	"bytes"
	"wallet"
)

/**
	交易输入
 */
type TxInput struct {
	TxId      []byte //对应out所在tx的id值
	outIndex  int    //对应out所在tx的索引
	Signature []byte //签名字节数组
	PubKey    []byte //公钥
}

/**
	对比in和out的公钥是否匹配
 */
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	return bytes.Compare(pubKeyHash, wallet.HashPubKey(in.PubKey)) == 0
}

/**
	当前交易能否用该公钥解锁
 */
func (in *TxInput) CanUnlockOutputWith(pubKeyHash []byte) bool {
	return bytes.Compare(wallet.HashPubKey(in.PubKey),pubKeyHash) == 0
}