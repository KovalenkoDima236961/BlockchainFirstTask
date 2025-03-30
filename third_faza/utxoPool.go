package third_faza

import (
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
)

// UTXOPool Táto štruktúra predstavuje pool UTXO, ktorý mapuje jednotlivé UTXO na ich zodpovedajúce transakčné výstupy
type UTXOPool struct {
	/**
	 * Aktuálna zbierka UTXO, pričom každé z nich je mapované na zodpovedajúci
	 * výstup transakcie
	 */
	H map[string]Output
}

// NewUTXOPool  Vytvorí nový prázdny UTXOPool
func NewUTXOPool() *UTXOPool {
	return &UTXOPool{H: make(map[string]Output)}
}

// NewUTXOPoolWithPool Vytvorí nový UTXOPool, ktorý je kópiou pool
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

// Put Pridá namapovanie z UTXO utxo do transackčného výstupu txOut * v poole
func (utxoPool *UTXOPool) Put(utxo UTXO, txOut Output) {
	utxoPool.H[utxo.Key()] = txOut
}

// RemoveUTXO Odstráni UTXO utxo z poolu
func (utxoPool *UTXOPool) RemoveUTXO(utxo UTXO) {
	delete(utxoPool.H, utxo.Key())
}

// GetTxOutput return výstup transakcie zodpovedajúci UTXO utxo alebo null, ak utxo nie je v poole.
func (utxoPool *UTXOPool) GetTxOutput(ut UTXO) *Output {
	if txOut, exists := utxoPool.H[ut.Key()]; exists {
		return &txOut
	}
	return nil
}

// Contains return true ak UTXO utxo je v poole a inak false
func (utxoPool *UTXOPool) Contains(utxo UTXO) bool {
	_, exists := utxoPool.H[utxo.Key()]
	return exists
}

// GetAllUTXO Vráti list všetkých UTXOs v poole
func (utxoPool *UTXOPool) GetAllUTXO() []*UTXO {
	utxos := make([]*UTXO, 0, len(utxoPool.H))
	for key := range utxoPool.H {
		if utxo, err := parseUTXOKey(key); err == nil {
			utxos = append(utxos, utxo)
		}
	}
	return utxos
}

// parseUTXOKey Konverzia stringu používaného ako mapový kľúč na UTXO
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
