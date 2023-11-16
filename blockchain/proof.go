package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// proof of work algo
// secure blockchain by doing work to check valid block
//steps

// take data from block
// create counter(nonce) starts 0
// create hash of data plus counter
// check hash to meets req
// retry till meet hash

// req: first few bytes 0s,
const difficulty = 12 // normally would increase to make harder over time

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork { // part target and create hash (similar to deriveHash)
	tar := big.NewInt(1)
	tar.Lsh(tar, uint(256-difficulty))

	pow := &ProofOfWork{b, tar}

	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte { // part target and create hash (similar to deriveHash)
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(difficulty)),
		}, []byte{},
	)
	return data
}

func ToHex(num int64) []byte { // part target and create hash (similar to deriveHash)
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func (pow *ProofOfWork) Run() (int, []byte) { // part target and create hash (similar to deriveHash)
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			// target lesss than number issue since signed numebr
			break
		} else {
			nonce++
		}

	}
	fmt.Println()

	return nonce, hash[:] // sign blcok at each step would need to work backwork
}

func (pow *ProofOfWork) Validate() bool { //use to prove hash is valid using nonce
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1 // confirm hash is correct value

}
