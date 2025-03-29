package first_faza

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
)

var (
	utxoPool *UTXOPool
)

// HandleTxs
// Creates a public ledger whose current UTXO pool (collection of unspent
// transaction outputs) is {@code utxoPool}. It should create a secure copy
// of utxoPool using the UTXOPool constructor (UTXOPool uPool).
func HandleTxs(utxoPool2 *UTXOPool) {
	if utxoPool2 == nil {
		panic("utxoPool2 is nil")
	}
	utxoPool = NewUTXOPoolWithPool(utxoPool2)
}

// UTXOPoolGet
// Returns the current UTXO pool.
// If no current UTXO pool is found, returns an empty (not nil) {@code UTXOPool} object.
func UTXOPoolGet() *UTXOPool {
	if utxoPool == nil {
		return NewUTXOPool()
	}
	return utxoPool
}

// TxIsValid
// Returns true if:
// (1) all outputs claimed by {@code tx} exist in the current UTXO pool,
// (2) the signatures on each {@code tx} input are valid,
// (3) no UTXO is claimed multiple times,
// (4) all output values of {@code tx} are non-negative, and
// (5) the sum of input values of {@code tx} is greater than or equal to the sum of its output values;
// otherwise returns false.
func TxIsValid(tx Transaction, pool *UTXOPool) bool {
	sumOfInputs := 0.0
	claimedUTXOs := make(map[string]bool)

	for i, input := range tx.Inputs {
		utxo := NewUTXO(input.PrevTxHash, input.OutputIndex)
		if _, ok := pool.H[utxo.Key()]; !ok {
			return false
		}
		output := pool.GetTxOutput(*utxo)
		data := tx.GetDataToSign(i)
		signature := input.Signature

		if !VerifySignature(data, signature, output.Address) {
			return false
		}
		if claimedUTXOs[utxo.Key()] {
			return false
		}
		claimedUTXOs[utxo.Key()] = true

		sumOfInputs += output.Value
	}

	sumOfOutputs := 0.0
	for _, output := range tx.Outputs {
		if output.Value < 0 {
			return false
		}
		sumOfOutputs += output.Value
	}

	return sumOfInputs >= sumOfOutputs
}

// VerifySignature checks whether the given signature is valid for the given message and RSA public key.
func VerifySignature(message []byte, signature []byte, address *rsa.PublicKey) bool {
	hash := sha256.Sum256(message)
	err := rsa.VerifyPKCS1v15(address, crypto.SHA256, hash[:], signature)
	return err == nil
}

// Handler
// Processes each epoch (iteration) by taking an unordered list of proposed
// transactions, verifying the correctness of each transaction, returning an array
// of mutually valid accepted transactions, and updating the current UTXO pool accordingly.
func Handler(possibleTxs []Transaction) []*Transaction {
	originalPool := NewUTXOPoolWithPool(utxoPool)
	validTxs := make([]*Transaction, 0)

	for i := range possibleTxs {
		tx := &possibleTxs[i]
		if TxIsValid(*tx, originalPool) {
			validTxs = append(validTxs, tx)

			for _, input := range tx.GetInputs() {
				originalPool.RemoveUTXO(UTXO{txHash: input.PrevTxHash, index: input.OutputIndex})
			}
			for j, output := range tx.GetOutputs() {
				originalPool.AddUTXO(UTXO{txHash: tx.GetHash(), index: j}, *output)
			}
		}
	}

	utxoPool = originalPool
	return validTxs
}
