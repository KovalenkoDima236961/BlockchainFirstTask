package third_faza

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
)

// UTXO represents an Unspent Transaction Output, uniquely identified by a transaction hash and an output index.
type UTXO struct {
	txHash []byte
	index  int
}

// NewUTXO creates a new UTXO instance with a copy of the given transaction hash and the specified index.
func NewUTXO(txHash []byte, index int) *UTXO {
	// Copy the transaction hash to avoid external modifications.
	return &UTXO{
		txHash: append([]byte{}, txHash...),
		index:  index,
	}
}

// GetTxHash returns the transaction hash associated with this UTXO.
func (u *UTXO) GetTxHash() []byte {
	return u.txHash
}

// GetIndex returns the output index of this UTXO.
func (u *UTXO) GetIndex() int {
	return u.index
}

// Equals checks if the given UTXO is equal to the current one.
// Two UTXOs are equal if their transaction hashes and indices are identical.
func (u *UTXO) Equals(other *UTXO) bool {
	if other == nil {
		return false
	}

	otherHash := other.GetTxHash()
	otherIndex := other.GetIndex()

	if otherIndex != u.index || len(otherHash) != len(u.txHash) {
		return false
	}

	for i := range u.txHash {
		if u.txHash[i] != otherHash[i] {
			return false
		}
	}
	return true
}

// CompareTo compares this UTXO with another.
// Returns -1 if this UTXO is "less" than the other, 1 if "greater", or 0 if they are equal.
// The comparison is done first by index, and if equal, by lexicographical order of the transaction hash.
func (u *UTXO) CompareTo(other *UTXO) int {
	if other == nil {
		return 1
	}

	otherHash := other.GetTxHash()
	otherIndex := other.GetIndex()

	if u.index != otherIndex {
		if u.index < otherIndex {
			return -1
		}
		return 1
	}

	// If indices are equal, compare the transaction hashes.
	if len(u.txHash) != len(otherHash) {
		if len(u.txHash) < len(otherHash) {
			return -1
		}
		return 1
	}

	for i := range u.txHash {
		if u.txHash[i] < otherHash[i] {
			return -1
		} else if u.txHash[i] > otherHash[i] {
			return 1
		}
	}
	return 0
}

// HashCode computes a hash code for the UTXO combining its index and a FNV hash of the transaction hash.
func (u *UTXO) HashCode() int {
	// Start with a seed value and combine with the index.
	hash := 1
	hash = hash*17 + u.index

	fnvHasher := fnv.New32a()
	fnvHasher.Write(u.txHash)
	hash = hash*31 + int(fnvHasher.Sum32())

	return hash
}

// Key returns a unique string representation of the UTXO, combining the hex-encoded transaction hash and the index.
func (u *UTXO) Key() string {
	return fmt.Sprintf("%s:%d", hex.EncodeToString(u.txHash), u.index)
}
