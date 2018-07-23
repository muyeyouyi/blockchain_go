package blockchain

import (
	"bytes"
	"encoding/gob"
)

/**
	序列化
*/
func Serialize(data interface{}) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	encoder.Encode(data)
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
