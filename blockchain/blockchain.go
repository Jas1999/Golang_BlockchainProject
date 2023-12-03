package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger/v3"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST" // file to confirm db started
	genesisData = "1st transaction from gensis8"
)

type BlockChain struct { // simple strucutre for chian
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct { // simple strucutre for chian
	CurrentHash []byte
	Database    *badger.DB
}

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// setup Genesis block
func InitBlockChain(addr string) *BlockChain { // take string and prev output pointer to blcok
	//return &BlockChain{[]*Block{Genesis()}}

	var lh []byte // last hash

	if DBexists() {
		fmt.Println("Blockchain exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath) // specifcy where to store
	opts.Dir = dbPath                     // dir stores keys and meta data
	opts.ValueDir = dbPath                // stores values  ( same loc in this case )

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(addr, genesisData) // addr mine block and be rewarded
		genesis := Genesis(cbtx)
		fmt.Println("created Genesis block")
		err = txn.Set(genesis.Hash, genesis.Serializer())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)

		lh = genesis.Hash

		return err

		// init db so need to write db
		// txn would allow us use transactions
		// lh (last hash)
		/*
			// pre transactions code
			if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
				fmt.Println(" no exisitng block found ")
				gen := Genesis()
				fmt.Println("gensis proved")
				err = txn.Set(gen.Hash, gen.Serializer())
				Handle(err)
				err = txn.Set([]byte("lh"), gen.Hash)
				lh = gen.Hash
				return err
			} else { // key exists
				item, err := txn.Get([]byte("lh"))
				Handle(err)
				//lh, err = item.Value()
				err = item.Value(func(val []byte) error {
					lh = append([]byte{}, val...)
					return nil
				})

				Handle(err)
				return err
			}
		*/
		//Handle(err)
		//return err
	})
	Handle(err)
	blockchain := BlockChain{lh, db}
	return &blockchain
}

func ContinueBlockChain(address string) *BlockChain { // use existing block chain
	if DBexists() == false {
		fmt.Println("No existing blockchain: create one!")
		runtime.Goexit()
	}

	var lh []byte                         // last hash
	opts := badger.DefaultOptions(dbPath) // specifcy where to store
	opts.Dir = dbPath                     // dir stores keys and meta data
	opts.ValueDir = dbPath                // stores values  ( same loc in this case )

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lh = append([]byte{}, val...)
			return nil
		})
		Handle(err)
		return err
	})
	Handle(err)
	chain := BlockChain{lh, db}
	return &chain
}

// add block to chain
func (chain *BlockChain) AddBlock(transactions []*Transaction) { // add block
	var lh []byte // last hash

	err := chain.Database.View(func(txn *badger.Txn) error {
		// init db so need to write db
		// txn would allow us use transactions
		// lh (last hash)
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		//lh, err = item.Value()
		err = item.Value(func(val []byte) error {
			lh = append([]byte{}, val...)
			return nil
		})

		return err
	})

	Handle(err)
	// create new block with lh
	newB := CreateBlock(transactions, lh) //(data, lh)

	err = chain.Database.Update(func(tx *badger.Txn) error {

		err := tx.Set(newB.Hash, newB.Serializer()) // update db
		Handle(err)
		err = tx.Set([]byte("lh"), newB.Hash)
		chain.LastHash = newB.Hash // update last hash to current for next item
		return err
	})
	Handle(err)

}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}
	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		// encodedBlock, err := item.Value()
		err = item.Value(func(val []byte) error {
			encodedBlock := append([]byte{}, val...)
			block = Deserialize(encodedBlock)
			return nil
		})
		return err
	})
	Handle(err)
	iter.CurrentHash = block.PrevHash

	return block
}

// find unspent ( outputs not in inputs)
// tokens left unspend (assigne dto user)

func (chain *BlockChain) FindUnspentTxns(addr string) []Transaction {

	var unspentTxns []Transaction

	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs: // label so breaks out specfic
			for outIdx, out := range tx.Outputs {
				for _, spentOut := range spentTXOs[txID] {
					if spentOut == outIdx {
						continue Outputs
					}

					if out.CanBeUnlocked(addr) {
						unspentTxns = append(unspentTxns, *tx) // add all transactions that can unlocked at address
					}

				}
				if tx.IsCoinbase() == false {
					for _, in := range tx.Inputs { // go through outputs and see which inputs we can add
						if in.CanUnlock(addr) {
							inTxID := hex.EncodeToString(in.ID)
							spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
						}
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}

	}
	return unspentTxns //
}

// find unspend transactions outputs
func (chain *BlockChain) FindUTXO(addr string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTxns(addr)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(addr) { // add unspent outputsb
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

// find unspent and check if amount exists, return accumlated amount and unspent outputs
func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxns := chain.FindUnspentTxns(address)
	accumulated := 0 //

Work:
	for _, tx := range unspentTxns { // got through unspent transactinos
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work // when accumlated has enought
				}
			}
		}
	}

	return accumulated, unspentOuts
}
