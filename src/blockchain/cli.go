package blockchain

import (
	"flag"
	"os"
	"utils"
	"fmt"
	"strconv"
)

const (
	printChain = "printchain"//命令行 打印链表
	addBlock = "addblock"//命令行 新增区块
)

type Cli struct {
	Chain *Chain
}

func (cli *Cli) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet(addBlock, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(printChain, flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")

	var err error
	switch os.Args[1] {

	case addBlock:
		err = addBlockCmd.Parse(os.Args[2:])
	case printChain:
		err = printChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}
	utils.LogE(err)

	if addBlockCmd.Parsed() {
		if  *addBlockData == ""{
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

/**
	验证命令行参数
 */
func (cli *Cli) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

/**
	打印提示信息
 */
func (cli *Cli) printUsage() {
	fmt.Println("Useage")
	fmt.Println("    addblock -data BLOCK_DATA （生成一个区块）")
	fmt.Println("    printchain  (打印全部区块链)")
}
func (cli *Cli) addBlock(data string) {
	cli.Chain.AddBlock(data)
}

/**
	打印链
 */
func (cli *Cli) printChain() {
	iterator := cli.Chain.Iterator()
	for iterator.CurrentHash != nil {
		block := iterator.Next()
		fmt.Printf("Data：%s\n", block.Data)
		fmt.Println("Timestamp：", block.Timestamp)
		fmt.Printf("Hash：%x\n", block.Hash)
		fmt.Printf("PreBlockHash：%x\n", block.PreBlockHash)
		pow := NewPow(block)
		fmt.Printf("pow:%s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
