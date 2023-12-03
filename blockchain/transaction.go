package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct { // made of inputs and outputs ( use inputs and outputs to derive values for transactiopns since everything is open)
	ID      []byte // hash
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	Value  int
	PubKey string // key to unlock token
}

type TxInput struct { // references to previous outs puts
	ID  []byte // transaction output relates to
	Out int    // index
	Sig string // used to get data
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func CoinbaseTx(to, data string) *Transaction { // create init transactions
	if data == "" {
		data = fmt.Sprintf("Coin to %s", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{100, to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return &tx

}

func (tx *Transaction) IsCoinbase() bool { // confirm coinbase transactions : check length, and confirm
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// both true confirm user owns data
func (in *TxInput) CanUnlock(data string) bool { // confirm sig
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}

func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput
	acc, validOutputs := chain.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("Error: not enough funds")
	}

	for txid, outs := range validOutputs { // go through valid outputs aand use tx id
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs { // create new input for unspent outputs
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to}) // output address to and the amount

	if acc > amount { // user has more funds so left over tokens in account
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs} // next txn based on in and outs
	tx.SetID()

	return &tx
}
