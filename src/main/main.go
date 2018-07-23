package main

import (
	"blockchain"
	"crypto/ecdsa"
	"golang.org/x/crypto/ripemd160"
)

func main() {
	cli := &blockchain.Cli{}
	cli.Run()
	defer func() {
		if cli.Chain != nil {
			cli.Chain.Db.Close()
		}
	}()

	ripemd160.New()

}

type name struct {
	ecdsa.PublicKey
}