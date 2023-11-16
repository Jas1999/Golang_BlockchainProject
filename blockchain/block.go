package blockchain

import (
	"bytes"
	"crypto/sha256"
)

type BlockChain struct { // simple strucutre for chian
	Blocks []*Block // B in blocks capital makes it public
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
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlockk(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}

// Genesis block creations
func Genesis() *Block { // take string and prev output pointer to blcok
	return CreateBlockk("Gensis", []byte{})
}

// setup Genesis block
func InitBlock() *BlockChain { // take string and prev output pointer to blcok
	return &BlockChain{[]*Block{Genesis()}}
}
