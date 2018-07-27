package wallet

import (
	"os"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
	"fmt"
)

type Wallets struct {
	Wallets map[string]*Wallet
}

/**
	创建钱包，从本地读取缓存
 */
func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFromFile()

	return &wallets, err
}

/**
	读取本地钱包信息
 */
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	ws.Wallets = wallets.Wallets

	return nil
}

/**
	保存钱包至本地
 */
func (ws Wallets) SaveToFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
/**
	创建一个新钱包
 */
func (ws *Wallets) CreateWallet() string{
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())
	ws.Wallets[address] = wallet
	return address
}

/**
	获取所有钱包地址
 */
func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) GetAddresses() []string{
	var addresses []string

	for key := range ws.Wallets {
		addresses = append(addresses,key)
	}
	return addresses
}
