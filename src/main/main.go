package main

import (
	"blockchain"
)

func main() {
	cli := &blockchain.Cli{}
	cli.Run()
	defer func() {
		if cli.Chain != nil {
			cli.Chain.Db.Close()
		}
	}()


}
