package main

import (
	"flag"
	"fmt"
	"log"
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
	//fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	//fmt.Println(" print - Prints the blocks in the chain")

	fmt.Println(" getbalance -address ADDRESS : get the balance for an address")
	fmt.Println(" createblockchain -address ADDRESS : creates a blockchain and sends genesis reward to address")
	fmt.Println(" printchain - Prints : the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount integer - : Send amount of coins")
}

func (cli *CMD_Line) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit() // shut downs golang ( proper shutdown so no data corrupted)
	}
}

//func (cli *CMD_Line) addBlock(data string) {
//	cli.blockchain.AddBlock(data)
//	fmt.Println("Added Block!")
//}

func (cli *CMD_Line) printChain() {
	//iter := cli.blockchain.Iterator()
	chain := blockchain.ContinueBlockChain("") // iter.Next() // get next each loop
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		// fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 { // back to gensis block which has no prev hash
			break
		}
	}
}

func (cli *CMD_Line) createBlockChain(addr string) {
	chain := blockchain.InitBlockChain(addr)
	chain.Database.Close()
	fmt.Println("Finished!")
}

func (cli *CMD_Line) getBalance(addr string) {
	chain := blockchain.ContinueBlockChain(addr)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(addr)

	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf(" Balance of %s : %d \n", addr, balance)

}

func (cli *CMD_Line) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CMD_Line) run() {
	cli.validateArgs()

	/*
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
	*/
	// EP 4

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
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
	//defer os.Exit(0) // another safety to close db properly
	//chain := blockchain.InitBlock()
	//defer chain.Database.Close() // wait for os.exit

	//cli := CMD_Line{chain}
	//cli.run()

	defer os.Exit(0)
	cli := CMD_Line{}
	cli.run()
}
