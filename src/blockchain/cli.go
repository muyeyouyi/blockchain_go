package blockchain

import (
	"flag"
	"os"
	"utils"
	"fmt"
	"constants"
	"github.com/boltdb/bolt"
	"strconv"
	"encoding/hex"
)

const (
	printChain  = "printchain"  //命令行 打印链表
	createChain = "createchain" //命令行 新增区块
	address     = "address"     //命令行 地址
	from        = "from"        //命令行 币发送方
	to          = "to"          //命令行 币接收方
	send        = "send"        //命令行 转账
	amount      = "amount"      //命令行 数量
	getBalance  = "getbalance"  //命令行 查询余额
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
	sendFromData := sendCmd.String(from, "", "在-from 后输入地址")
	sendToData := sendCmd.String(to, "", "在-to 后输入地址")
	sendAmountData := sendCmd.Int(amount, 0, "在-amount 后输入币的数量")

	//cmd创建链和创世区块
	getBalanceCmd := flag.NewFlagSet(getBalance, flag.ExitOnError)
	getBalanceAddressData := getBalanceCmd.String(address, "", "在-address 后输入地址")

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
		err = sendCmd.Parse(os.Args[2:])
	case getBalance:
		err = getBalanceCmd.Parse(os.Args[2:])
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

	//获取地址的余额
	if getBalanceCmd.Parsed() {
		if *getBalanceAddressData == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddressData)
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
	fmt.Println("    ", send, " -",from," abc -",to," xyz -",amount," 30 （发送一笔交易）")
	fmt.Println("    ", printChain, "                 (打印全部区块)")
	fmt.Println("    ", createChain, " -", address, " dfz (生成创世区块)")
}

/**
	创建区块和交易
 */
func (cli *Cli) addBlock(send, to string, amount int) {
	if cli.Chain == nil {
		cli.linkDb()
	}

	if cli.Chain == nil {
		fmt.Println("错误：创世区块尚未创建")
		cli.printUsage()
	} else {
		tx := newUTXOTransaction(send, to, amount, cli.Chain)
		if tx != nil {
			cli.Chain.AddBlock([]*Transaction{tx})
		}

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

/**
	创建一笔转账交易
 */
func newUTXOTransaction(send, to string, amount int, chain *Chain) *Transaction {
	count, validOutputs := chain.FindSpendableOutputs(send, amount)
	if count < amount {
		utils.LogD("余额不足，请重新确认")
		return nil
	}

	var inputs []TxInput
	var outputs []TxOutput

	//创建input
	for txId, outIndexs := range validOutputs {
		decodeId, error := hex.DecodeString(txId)
		utils.LogE(error)
		for _, index := range outIndexs {
			inputs = append(inputs,TxInput{decodeId,index,send})
		}
	}

	//创建output
	outputs = append(outputs,TxOutput{amount,to})
	if	count > amount{
		outputs = append(outputs,TxOutput{count-amount,send})
	}
	//生成tx
	tx := &Transaction{nil, inputs, outputs}
	tx.setId()

	return tx
}

/**
	打印数据库为空的信息
 */
func (cli *Cli) printDbEmpty() {
	fmt.Println("未发现区块链，请先创建创世区块")
	cli.printUsage()
}

/**
	查询余额
 */
func (cli *Cli) getBalance(address string) {
	if cli.Chain == nil {
		cli.linkDb()
	}
	balance := 0
	if cli.Chain != nil {

		utxos := cli.Chain.FindUTXOs(address)
		for _, output := range utxos {
			balance += output.ValueOut
		}
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
