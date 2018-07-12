package main

import (
	"./blockchain"
	"bytes"
	"strconv"
	"fmt"
)

func main() {
	run()
}
func run() {
	//创建区块链，生成4个区块
	//blockchain.NewBlockChain()
	bc := blockchain.NewBlockChain()
	for i := 1; i < 10; i++ {
		buffer := bytes.Buffer{}
		buffer.WriteString("第")
		buffer.WriteString(strconv.Itoa(i))
		buffer.WriteString("个区块")
		bc.AddBlock(buffer.String())
	}
	iterator := bc.Iterator()
	for iterator.CurrentHash != nil {
		block := iterator.Next()
		fmt.Printf("Data：%s\n", block.Data)
		fmt.Println("Timestamp：", block.Timestamp)
		fmt.Printf("Hash：%x\n", block.Hash)
		fmt.Printf("PreBlockHash：%x\n", block.PreBlockHash)
		pow := blockchain.NewPow(block)
		fmt.Printf("pow:%s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
	//
	//bc.AddBlock("第2个区块")
	//bc.AddBlock("第3个区块")

	////遍历区块链打印
	//for _, block := range bc.Blocks {
	//	fmt.Printf("Data：%s\n", block.Data)
	//	fmt.Println("Timestamp：", block.Timestamp)
	//	fmt.Printf("Hash：%x\n", block.Hash)
	//	fmt.Printf("PreBlockHash：%x\n", block.PreBlockHash)
	//	pow := blockchain.NewPow(block)
	//	fmt.Printf("pow:%s\n", strconv.FormatBool(pow.Validate()))
	//	fmt.Println()
	//}
}
