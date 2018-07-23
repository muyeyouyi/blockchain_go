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
	address     = "address"     //命令行 地址
	from        = "from"        //命令行 币发送方
	to          = "to"          //命令行 币接收方
	send        = "send"        //命令行 转账
	amount      = "amount"      //命令行 数量
)

type Cli struct {
	Chain *Chain
}

/**
	运行命令行
 */
func (cli *Cli) Run() {

	//输出提示信息
	cli.validateArgs()

	//cmd创建链和创世区块
	createChainCmd := flag.NewFlagSet(createChain, flag.ExitOnError)
	createChainAddressData := createChainCmd.String(address, "", "在-address 后输入地址")

	//cmd创建链和创世区块
	sendCmd := flag.NewFlagSet(send, flag.ExitOnError)
	sendFromData := createChainCmd.String(from, "", "在-from 后输入地址")
	sendToData := createChainCmd.String(to, "", "在-to 后输入地址")
	sendAmountData := createChainCmd.Int(amount, 0, "在-amount 后输入币的数量")

	//cmd打印链
	printChainCmd := flag.NewFlagSet(printChain, flag.ExitOnError)

	//截取命令行内容
	var err error
	switch os.Args[1] {
	case printChain:
		err = printChainCmd.Parse(os.Args[2:])
	case createChain:
		err = createChainCmd.Parse(os.Args[2:])
	case send:
		err = createChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}
	utils.LogE(err)

	//发送交易，创建区块
	if sendCmd.Parsed() {
		if *sendFromData == "" || *sendToData == "" || *sendAmountData <= 0 {
			os.Exit(1)
		}
		cli.addBlock(*sendFromData, *sendToData, *sendAmountData)
	}

	//获取创世区块的string 创建区块
	if createChainCmd.Parsed() {
		if *createChainAddressData == "" {
			createChainCmd.Usage()
			os.Exit(1)
		}
		cli.NewBlockChain(*createChainAddressData)
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
	//fmt.Println("    ", addBlock, " -data BLOCK_DATA （生成一个区块）")
	fmt.Println("    ", printChain, "                 (打印全部区块)")
	fmt.Println("    ", createChain, " -", address, " dfz (生成创世区块)")
}

/**
	添加区块
 */
func (cli *Cli) addBlock(send, to string, amount int) {
	if cli.Chain == nil {
		cli.linkDb()
	}

	if cli.Chain == nil {
		fmt.Println("错误：创世区块尚未创建")
		cli.printUsage()
	} else {
		cli.Chain.AddBlock(newTransaction(send, to, amount))
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
			//fmt.Printf("Data：%s\n", block.Transactions)
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
func (cli *Cli) NewBlockChain(address string) {
	tx := NewCoinBaseTx(address)
	chain := NewBlockChain(tx)
	cli.Chain = chain
}

func newTransaction(send, to string, amount int) []*Transaction {
	//todo
	return nil
}

/**
	打印数据库为空的信息
 */
func (cli *Cli) printDbEmpty() {
	fmt.Println("未发现区块链，请先创建创世区块")
	cli.printUsage()
}
