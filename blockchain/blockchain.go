package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct { // simple strucutre for chian
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct { // simple strucutre for chian
	CurrentHash []byte
	Database    *badger.DB
}

// setup Genesis block
func InitBlock() *BlockChain { // take string and prev output pointer to blcok
	//return &BlockChain{[]*Block{Genesis()}}

	var lh []byte // last hash

	opts := badger.DefaultOptions(dbPath) // specifcy where to store
	opts.Dir = dbPath                     // dir stores keys and meta data
	opts.ValueDir = dbPath                // stores values  ( same loc in this case )

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		// init db so need to write db
		// txn would allow us use transactions
		// lh (last hash)
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
		//Handle(err)
		//return err
	})

	blockchain := BlockChain{lh, db}
	return &blockchain
}

// add block to chain
func (chain *BlockChain) AddBlock(data string) { // add block
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
	newB := CreateBlock(data, lh)

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
