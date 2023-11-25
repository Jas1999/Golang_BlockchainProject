package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/Jas1999/Golang_BlockchainProject/blockchain"
)

type CMD_Line struct {
	blockchain *blockchain.BlockChain
}

func (cli *CMD_Line) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *CMD_Line) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit() // shut downs golang ( proper shutdown so no data corrupted)
	}
}

func (cli *CMD_Line) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

func (cli *CMD_Line) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next() // get next each loop

		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 { // back to gensis block which has no prev hash
			break
		}
	}
}

func (cli *CMD_Line) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data") // data passed as string

	switch os.Args[1] { // first argument
	case "add": // go run main.go add -block "1st block"
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "print": // go run main.go print
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func main() {

	/*
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
	*/
	defer os.Exit(0) // another safety to close db properly
	chain := blockchain.InitBlock()
	defer chain.Database.Close() // wait for os.exit

	cli := CMD_Line{chain}
	cli.run()
}
