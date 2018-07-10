package blockchain

import (
	"bytes"
	"encoding/gob"
)

/**
	序列化
 */
func Serialize(block *Block) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	encoder.Encode(block)
	return result.Bytes()
}

/**
	反序列化
 */
func Deserialize(serialize []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(serialize))
	decoder.Decode(&block)
	return &block
}