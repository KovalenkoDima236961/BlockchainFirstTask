package third_faza

import (
	"crypto/rsa"
	"crypto/sha256"
)

const COINBASE = 3.125

type Block struct {
	hash          []byte
	prevBlockHash []byte
	coinbase      *Transaction
	txs           []*Transaction
}

func NewBlock(prevHash []byte, address *rsa.PublicKey) *Block {
	coinbase := NewCoinbaseTransaction(COINBASE, address)
	newBlock := &Block{
		prevBlockHash: append([]byte{}, prevHash...),
		coinbase:      coinbase,
		txs:           []*Transaction{coinbase},
	}
	return newBlock
}

func (block *Block) GetCoinbase() *Transaction {
	return block.coinbase
}

func (block *Block) GetHash() []byte {
	return block.hash
}

func (block *Block) GetPrevBlockHash() []byte {
	return block.prevBlockHash
}

func (block *Block) GetTransactions() []*Transaction {
	copyTxs := make([]*Transaction, len(block.txs))
	copy(copyTxs, block.txs)
	return copyTxs
}

func (block *Block) GetTransaction(index int) *Transaction {
	return block.txs[index]
}

func (block *Block) TransactionAdd(tx *Transaction) {
	block.txs = append(block.txs, tx)
}

func (block *Block) GetBlock() []byte {
	rawBlock := make([]byte, 0)

	if block.prevBlockHash != nil {
		rawBlock = append(rawBlock, block.prevBlockHash...)
	}
	for _, tx := range block.txs {
		rawBlock = append(rawBlock, tx.GetTx()...)
	}

	return rawBlock
}

func (block *Block) Finalizee() {
	hash := sha256.Sum256(block.GetBlock())
	block.hash = hash[:]
}
