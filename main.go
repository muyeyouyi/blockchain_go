package main

import (
	"fmt"
	"./blockchain"
	"strconv"
)

func main() {
	testBc()
}
func testBc() {
	bc := blockchain.CreateBlockChain()
	bc.AddBlock("第1个区块")
	bc.AddBlock("第2个区块")
	bc.AddBlock("第3个区块")
	for _, block := range bc.Blocks {
		fmt.Printf("Data：%s\n", block.Data)
		fmt.Println("Timestamp：", block.Timestamp)
		fmt.Printf("Hash：%x\n", block.Hash)
		fmt.Printf("PreBlockHash：%x\n", block.PreBlockHash)
		pow := blockchain.NewPow(block)
		fmt.Printf("pow:%s\n",strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}



