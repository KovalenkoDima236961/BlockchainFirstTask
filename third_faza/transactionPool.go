package third_faza

import (
	"encoding/hex"
)

// TransactionPool represents a pool of transactions, keyed by the hex-encoded transaction hash.
type TransactionPool struct {
	H map[string]*Transaction
}

// NewTransactionPool creates a new empty TransactionPool.
func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		H: make(map[string]*Transaction),
	}
}

// NewTransactionPoolFromPool creates a new TransactionPool that is a copy of an existing one.
func NewTransactionPoolFromPool(tp *TransactionPool) *TransactionPool {
	newPool := &TransactionPool{
		H: make(map[string]*Transaction),
	}
	for k, tx := range tp.H {
		newPool.H[k] = tx
	}
	return newPool
}

func keyFor(txHash []byte) string {
	wrapper := NewByteArrayWrapper(txHash)
	return hex.EncodeToString(wrapper.contents)
}

// AddTransaction adds the given transaction to the pool, using its hash as the key.
func (tp *TransactionPool) AddTransaction(tx *Transaction) {
	key := keyFor(tx.GetHash())
	tp.H[key] = tx
}

// RemoveTransaction removes the transaction with the given hash from the pool.
func (tp *TransactionPool) RemoveTransaction(txHash []byte) {
	key := keyFor(txHash)
	delete(tp.H, key)
}

// GetTransaction returns the transaction associated with the given hash, or nil if not found.
func (tp *TransactionPool) GetTransaction(txHash []byte) *Transaction {
	key := keyFor(txHash)
	return tp.H[key]
}

// GetTransactions returns a slice containing all transactions in the pool.
func (tp *TransactionPool) GetTransactions() []*Transaction {
	txs := make([]*Transaction, 0, len(tp.H))
	for _, tx := range tp.H {
		txs = append(txs, tx)
	}
	return txs
}
