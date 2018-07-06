package main

import (
	"fmt"
	"./blockchain"
)

func main() {
	bc := blockchain.CreateBlockChain()
	bc.AddBlock("第1个区块")
	bc.AddBlock("第2个区块")
	bc.AddBlock("第3个区块")

	for _, block := range bc.Blocks {
		fmt.Printf("Data：%s\n",block.Data)
		fmt.Println("Timestamp：",block.Timestamp)
		fmt.Printf("Hash：%x\n",block.Hash)
		fmt.Printf("PreBlockHash：%x\n\n",block.PreBlockHash)

	}
}



