package third_faza

import (
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
)

// UTXOPool represents a UTXO pool that maps individual UTXOs to their corresponding transaction outputs
type UTXOPool struct {
	H map[string]Output
}

// NewUTXOPool creates a new empty UTXOPool.
func NewUTXOPool() *UTXOPool {
	return &UTXOPool{H: make(map[string]Output)}
}

// NewUTXOPoolWithPool creates a new UTXOPool that is a copy of the provided pool.
func NewUTXOPoolWithPool(pool *UTXOPool) *UTXOPool {
	newPool := &UTXOPool{H: make(map[string]Output)}
	for k, v := range pool.H {
		var multiSigCopy []*rsa.PublicKey
		if v.MultiSigAddresses != nil {
			multiSigCopy = append([]*rsa.PublicKey(nil), v.MultiSigAddresses...)
		}
		newPool.H[k] = Output{
			Value:             v.Value,
			Address:           v.Address,
			MultiSigAddresses: multiSigCopy,
		}
	}
	return newPool
}

// Put adds a mapping from UTXO utxo to the transaction output txOut in the pool.
func (utxoPool *UTXOPool) Put(utxo UTXO, txOut Output) {
	utxoPool.H[utxo.Key()] = txOut
}

// RemoveUTXO removes the UTXO utxo from the pool.
func (utxoPool *UTXOPool) RemoveUTXO(utxo UTXO) {
	delete(utxoPool.H, utxo.Key())
}

// GetTxOutput returns the transaction output corresponding to UTXO utxo, or null if the utxo is not in the pool.
func (utxoPool *UTXOPool) GetTxOutput(ut UTXO) *Output {
	if txOut, exists := utxoPool.H[ut.Key()]; exists {
		return &txOut
	}
	return nil
}

// Contains returns true if the UTXO utxo is in the pool, false otherwise.
func (utxoPool *UTXOPool) Contains(utxo UTXO) bool {
	_, exists := utxoPool.H[utxo.Key()]
	return exists
}

// GetAllUTXO returns a list of all UTXOs in the pool.
func (utxoPool *UTXOPool) GetAllUTXO() []*UTXO {
	utxos := make([]*UTXO, 0, len(utxoPool.H))
	for key := range utxoPool.H {
		if utxo, err := parseUTXOKey(key); err == nil {
			utxos = append(utxos, utxo)
		}
	}
	return utxos
}

// parseUTXOKey converts a string used as a map key to a UTXO.
func parseUTXOKey(key string) (*UTXO, error) {
	parts := strings.Split(key, ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid UTXO key format")
	}

	txHash, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

	index, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	return &UTXO{
		txHash: txHash,
		index:  index,
	}, nil
}
