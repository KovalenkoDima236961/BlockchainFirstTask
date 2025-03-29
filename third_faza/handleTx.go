package third_faza

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
)

var (
	utxoPool *UTXOPool
)

/**
 * Vytvorí verejný ledger (účtovnú knihu), ktorého aktuálny UTXOPool (zbierka nevyčerpaných
 * transakčných výstupov) je {@code utxoPool}. Malo by to vytvoriť bezpečnú kópiu
 * utxoPool pomocou konštruktora UTXOPool (UTXOPool Pool).
 */
func HandleTxs(utxoPool2 *UTXOPool) {
	if utxoPool2 == nil {
		panic("utxoPool2 is nil")
	}
	utxoPool = NewUTXOPoolWithPool(utxoPool2)
}

/**
 * return aktuálny UTXO pool.
 * Ak nenájde žiadny aktuálny UTXO pool, tak vráti prázdny (nie nulový) objekt {@code UTXOPool}.
 */
func UTXOPoolGet() *UTXOPool {
	if utxoPool == nil {
		return NewUTXOPool()
	}
	return utxoPool
}

/**
 * return true, ak
 * (1) sú všetky nárokované výstupy {@code tx} v aktuálnom UTXO pool,
 * (2) podpisy na každom vstupe {@code tx} sú platné,
 * (3) žiadne UTXO nie je nárokované viackrát,
 * (4) všetky výstupné hodnoty {@code tx}s sú nezáporné a
 * (5) súčet vstupných hodnôt {@code tx}s je väčší alebo rovný súčtu jej
 *     výstupných hodnôt; a false inak.
 */
func TxIsValid(tx Transaction, pool *UTXOPool) bool {
	sumOfInputs := 0.0
	claimedUTXOs := make(map[string]bool)

	if tx.Coinbase {
		return true
	}

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

	// All outputs must be non-negative.
	sumOfOutputs := 0.0
	for _, output := range tx.Outputs {
		if output.Value < 0 {
			return false
		}
		sumOfOutputs += output.Value
	}

	return sumOfInputs >= sumOfOutputs
}

func VerifySignature(message []byte, signature []byte, address *rsa.PublicKey) bool {
	hash := sha256.Sum256(message)
	err := rsa.VerifyPKCS1v15(address, crypto.SHA256, hash[:], signature)
	return err == nil
}

/**
 * Spracováva každú epochu (iteráciu) prijímaním neusporiadaného radu navrhovaných
 * transakcií, kontroluje správnosť každej transakcie, vracia pole vzájomne
 * platných prijatých transakcií a aktualizuje aktuálny UTXO pool podľa potreby.
 */
func Handler(possibleTxs []*Transaction) []*Transaction {
	originalPool := NewUTXOPoolWithPool(utxoPool)
	validTxs := make([]*Transaction, 0)

	for i := range possibleTxs {
		tx := possibleTxs[i]
		if TxIsValid(*tx, originalPool) {
			validTxs = append(validTxs, tx)

			for _, input := range tx.GetInputs() {
				originalPool.RemoveUTXO(UTXO{txHash: input.PrevTxHash, index: input.OutputIndex})
			}
			for j, output := range tx.GetOutputs() {
				originalPool.Put(UTXO{txHash: tx.GetHash(), index: j}, *output)
			}
		}
	}

	utxoPool = originalPool
	return validTxs
}
