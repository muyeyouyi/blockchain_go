package blockchain

import (
	"flag"
	"os"
	"utils"
	"fmt"
	"constants"
	"github.com/boltdb/bolt"
	"strconv"
)

const (
	printChain  = "printchain"  //命令行 打印链表
	addBlock    = "addblock"    //命令行 新增区块
	createChain = "createchain" //命令行 新增区块
)

type Cli struct {
	Chain *Chain
}

/**
	运行命令行
 */
func (cli *Cli) Run() {
	//cli.printChain()
	//
	////cli.NewBlockChain("创世区块")
	//return

	//输出提示信息
	cli.validateArgs()

	//cmd创建区块
	addBlockCmd := flag.NewFlagSet(addBlock, flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "在-data 后输入区块的内容")

	//cmd创建链和创世区块
	createChainCmd := flag.NewFlagSet(createChain, flag.ExitOnError)
	createChainData := createChainCmd.String("data", "", "在-data 后输入区块的内容")

	//cmd打印链
	printChainCmd := flag.NewFlagSet(printChain, flag.ExitOnError)

	//截取命令行内容
	var err error
	switch os.Args[1] {
	case addBlock:
		err = addBlockCmd.Parse(os.Args[2:])
	case printChain:
		err = printChainCmd.Parse(os.Args[2:])
	case createChain:
		err = createChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}
	utils.LogE(err)

	//获取新增区块的string 创建区块
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	//获取创世区块的string 创建区块
	if createChainCmd.Parsed() {
		if *createChainData == "" {
			createChainCmd.Usage()
			os.Exit(1)
		}
		cli.NewBlockChain(*createChainData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
func (cli *Cli) linkDb() {
	var tip []byte
	utils.CreateFile(constants.DbFile)
	//打开数据库
	db, err := bolt.Open(constants.DbFile, 0600, nil)
	//数据库写入标准
	err = db.View(func(tx *bolt.Tx) error {
		//检查表是否存在
		b := tx.Bucket([]byte(constants.BlocksBucket))
		//
		if b != nil {
			tip = b.Get(constants.LastHash)
			chain := Chain{tip, db}
			cli.Chain = &chain
		} else {
			cli.printDbEmpty()
		}
		return nil
	})
	utils.LogE(err)
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
	fmt.Println("用法：")
	fmt.Println("    ", addBlock, " -data BLOCK_DATA （生成一个区块）")
	fmt.Println("    ", printChain, "                (打印全部区块)")
	fmt.Println("     ", createChain, "-data BLOCK_DATA                (生成创世区块)")
}
/**
	添加区块
 */
func (cli *Cli) addBlock(data string) {
	if cli.Chain == nil {
		cli.linkDb()
	}

	if cli.Chain == nil {
		fmt.Println("错误：创世区块尚未创建")
		cli.printUsage()
	}else{
		cli.Chain.AddBlock(data)
	}
}

/**
	打印链
 */
func (cli *Cli) printChain() {
	if cli.Chain == nil {
		cli.linkDb()
	}

	if cli.Chain != nil {
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

}

/**
	创建区块链
 */
func (cli *Cli) NewBlockChain(data string) {
	chain := NewBlockChain(data)
	cli.Chain = chain
}
/**
	打印数据库为空的信息
 */
func (cli *Cli) printDbEmpty() {
	fmt.Println("未发现区块链，请先创建创世区块")
	cli.printUsage()
}
