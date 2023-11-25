package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

/*
// moved to blockchain.go
type BlockChain struct { // simple strucutre for chian
	Blocks []*Block // B in blocks capital makes it public
}
*/
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

/*
// old version done by proof of work now
func (b *Block) DeriveHash() { // method to create hash based on data and previous
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{}) // 2d of data and prev, and empty bytpes
	hash := sha256.Sum256(info)                                // change data changes hash otherwise same
	b.Hash = hash[:]
}
*/
func CreateBlock(data string, prevHash []byte) *Block { // take string and prev output pointer to blcok
	block := &Block{[]byte{}, []byte(data), prevHash, 0} // new block

	// block.DeriveHash()

	// new derive using proof of work
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

/*
// moved to blockchain.go
// add block to chain
func (chain *BlockChain) AddBlock(data string) { // add block
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlockk(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}
*/
// Genesis block creations
func Genesis() *Block { // take string and prev output pointer to blcok
	return CreateBlock("Gensis", []byte{})
}

/*
// moved to blockchain.go
// setup Genesis block
func InitBlock() *BlockChain { // take string and prev output pointer to blcok
	return &BlockChain{[]*Block{Genesis()}}
}
*/

func (b *Block) Serializer() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	Handle(err)

	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
