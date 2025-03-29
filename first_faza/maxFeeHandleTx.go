package first_faza

import "sort"

// MaxFeeHandler processes transactions by selecting only valid ones with positive fees.
// It returns a list of accepted transactions sorted by descending fee.
func MaxFeeHandler(transaction []*Transaction) []*Transaction {
	newPool := NewUTXOPoolWithPool(utxoPool)
	accepted := make([]*Transaction, 0)
	feeMap := make(map[string]float64)

	for _, tx := range transaction {
		fee := GetFee(tx, newPool)
		if TxIsValid(*tx, newPool) && fee > 0 {
			accepted = append(accepted, tx)
			feeMap[tx.Key()] = fee

			for _, input := range tx.GetInputs() {
				newPool.RemoveUTXO(UTXO{txHash: input.PrevTxHash, index: input.OutputIndex})
			}
			tx.Finalize()

			for i, output := range tx.GetOutputs() {
				newPool.AddUTXO(UTXO{txHash: tx.GetHash(), index: i}, *output)
			}
		}
	}

	utxoPool = newPool

	sort.Slice(accepted, func(i, j int) bool {
		return feeMap[accepted[i].Key()] > feeMap[accepted[j].Key()]
	})

	return accepted
}

// GetFee calculates the fee for a transaction as the difference
// between the total input value and total output value.
// Returns -1 if any input references an invalid UTXO.
func GetFee(tx *Transaction, pool *UTXOPool) float64 {
	totalInputValue := 0.0
	totalOutputValue := 0.0

	for _, input := range tx.Inputs {
		utxo := NewUTXO(input.PrevTxHash, input.OutputIndex)
		output := pool.GetTxOutput(*utxo)
		if output == nil {
			return -1
		}
		totalInputValue += output.Value
	}

	for _, output := range tx.Outputs {
		totalOutputValue += output.Value
	}
	return totalInputValue - totalOutputValue
}
