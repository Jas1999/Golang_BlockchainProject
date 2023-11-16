package main

import (
	"fmt"
	"strconv"

	"github.com/Jas1999/Golang_BlockchainProject/blockchain"
)

func main() {

	chain := blockchain.InitBlock()
	chain.AddBlock("1st")
	chain.AddBlock("2st")
	chain.AddBlock("3st")
	chain.AddBlock("4st")

	for _, block := range chain.Blocks {
		fmt.Printf("prev hash : %x\n", block.PrevHash)
		fmt.Printf("data : %s\n", block.Data)
		fmt.Printf("hash : %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s \n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
