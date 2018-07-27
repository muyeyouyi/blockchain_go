package main

import (
	"blockchain"
)

func main() {
	run()
}

func run() {
	cli := &blockchain.Cli{}
	cli.Run()
	defer func() {
		if cli.Chain != nil {
			cli.Chain.Db.Close()
		}
	}()
}
