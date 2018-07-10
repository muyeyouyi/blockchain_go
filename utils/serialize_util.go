package utils

import (
	"../blockchain"
	"bytes"
	"encoding/gob"
)

/**
	序列化
 */
func Serialize(block *blockchain.Block) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	encoder.Encode(block)
	return result.Bytes()
}

/**
	反序列化
 */
func Deserialize(serialize []byte) *blockchain.Block {
	var block blockchain.Block
	decoder := gob.NewDecoder(bytes.NewReader(serialize))
	decoder.Decode(&block)
	return &block
}