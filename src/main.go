package main

import (
	"./blockchain"
	"bytes"
	"fmt"
	"strconv"

)

func main() {
	//if len(os.Args) != 0 {
	//	fmt.Println(os.Args[0]) // args 第一个片 是文件路径
	//}
	//fmt.Println(os.Args[1])  // 第二个参数是， 用户输入的参数 例如 go run osdemo01.go 123
	//len := len(os.Args)
	//for i := 0; len ; i++ {
	//	fmt.Println(os.Args[i])
	//}
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
