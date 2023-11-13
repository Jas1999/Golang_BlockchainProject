package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type BlockChain struct { // simple strucutre for chian
	blocks []*Block
}
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

func (b *Block) DeriveHash() { // method to create hash based on data and previous
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{}) // 2d of data and prev, and empty bytpes
	hash := sha256.Sum256(info)                                // change data changes hash otherwise same
	b.Hash = hash[:]
}

func CreateBlockk(data string, prevHash []byte) *Block { // take string and prev output pointer to blcok
	block := &Block{[]byte{}, []byte(data), prevHash} // new block
	block.DeriveHash()
	return block
}

// add block to chain
func (chain *BlockChain) AddBlock(data string) { // add block
	prevBlock := chain.blocks[len(chain.blocks)-1]
	new := CreateBlockk(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, new)
}

// Genesis block creations
func Genesis() *Block { // take string and prev output pointer to blcok
	return CreateBlockk("Gensis", []byte{})
}

// setup Genesis block
func InitBlock() *BlockChain { // take string and prev output pointer to blcok
	return &BlockChain{[]*Block{Genesis()}}
}

func main() {

	chain := InitBlock()
	chain.AddBlock("1st")
	chain.AddBlock("2st")
	chain.AddBlock("3st")
	chain.AddBlock("4st")

	for _, block := range chain.blocks {
		fmt.Printf("prev hash : %x\n", block.PrevHash)
		fmt.Printf("data : %s\n", block.Data)
		fmt.Printf("hash : %x\n", block.Hash)
	}
}
