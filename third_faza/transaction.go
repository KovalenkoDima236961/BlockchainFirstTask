package third_faza

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"math"
	"time"
)

// Input represents a transaction input.
// It refers to a previous transaction's output that is being spent.
type Input struct {
	PrevTxHash        []byte
	OutputIndex       int
	Signature         []byte
	MultiSigSignature [][]byte
}

// NewInput creates a new Input instance with a copy of the previous transaction hash and a specified output index.
func NewInput(prevHash []byte, index int) *Input {
	in := &Input{
		PrevTxHash:  make([]byte, len(prevHash)),
		OutputIndex: index,
	}
	copy(in.PrevTxHash, prevHash)
	return in
}

// AddSignature assigns a signature to the input by copying the provided signature.
func (in *Input) AddSignature(sig []byte) {
	in.Signature = make([]byte, len(sig))
	copy(in.Signature, sig)
}

func (in *Input) AddMultiSignature(sig []byte) {
	in.MultiSigSignature = append(in.MultiSigSignature, sig)
}

func (in *Input) Equals(other *Input) bool {
	if other == nil {
		return false
	}

	if len(in.PrevTxHash) != len(other.PrevTxHash) {
		return false
	}
	for i := 0; i < len(in.PrevTxHash); i++ {
		if in.PrevTxHash[i] != other.PrevTxHash[i] {
			return false
		}
	}

	if in.OutputIndex != other.OutputIndex {
		return false
	}
	if len(in.Signature) != len(other.Signature) {
		return false
	}
	for i := 0; i < len(in.Signature); i++ {
		if in.Signature[i] != other.Signature[i] {
			return false
		}
	}
	return true
}

// Output represents a transaction output.
// It includes a value (in bitcoins) and a recipient's RSA public key (serving as an address).
type Output struct {
	Value             float64
	Address           *rsa.PublicKey
	MultiSigAddresses []*rsa.PublicKey
}

// NewOutput creates a new Output with the specified value and recipient address.
func NewOutput(value float64, address *rsa.PublicKey) *Output {
	return &Output{
		Value:   value,
		Address: address,
	}
}

func NewMultiSigOutput(value float64, addresses []*rsa.PublicKey) *Output {
	return &Output{
		Value:             value,
		MultiSigAddresses: addresses,
	}
}

// Equals checks if two outputs are identical by comparing both the value and the recipient's address.
func (out *Output) Equals(other *Output) bool {
	if other == nil {
		return false
	}

	if out.Value != other.Value {
		return false
	}
	if !out.Address.Equal(other.Address) {
		return false
	}
	return true
}

// Transaction represents a blockchain transaction.
// It contains a unique hash, a list of inputs, and a list of outputs.
type Transaction struct {
	Hash      []byte
	Inputs    []*Input
	Outputs   []*Output
	Coinbase  bool
	Timestamp int64
}

// NewTransaction creates a new transaction with empty slices for inputs and outputs.
func NewTransaction() *Transaction {
	return &Transaction{
		Inputs:    make([]*Input, 0),
		Outputs:   make([]*Output, 0),
		Coinbase:  false,
		Timestamp: time.Now().Unix(),
	}
}

// NewTransactionFromTransaction creates a deep copy of an existing transaction.
// It copies the hash, all inputs (including their signatures), and outputs.
func NewTransactionFromTransaction(tx *Transaction) *Transaction {
	newTx := &Transaction{
		Hash:      make([]byte, len(tx.Hash)),
		Inputs:    make([]*Input, len(tx.Inputs)),
		Outputs:   make([]*Output, len(tx.Outputs)),
		Coinbase:  false,
		Timestamp: tx.Timestamp,
	}
	copy(newTx.Hash, tx.Hash)

	for i, in := range tx.Inputs {
		newSig := make([]byte, len(in.Signature))
		copy(newSig, in.Signature)
		newTx.Inputs[i] = &Input{
			PrevTxHash:        append([]byte{}, in.PrevTxHash...),
			OutputIndex:       in.OutputIndex,
			Signature:         newSig,
			MultiSigSignature: in.MultiSigSignature,
		}
	}

	for i, op := range tx.Outputs {
		newTx.Outputs[i] = &Output{
			Value:             op.Value,
			Address:           op.Address,
			MultiSigAddresses: op.MultiSigAddresses,
		}
	}

	return newTx
}

func NewCoinbaseTransaction(coin float64, address *rsa.PublicKey) *Transaction {
	newTx := &Transaction{
		Inputs:    make([]*Input, 0),
		Outputs:   make([]*Output, 0),
		Coinbase:  true,
		Timestamp: time.Now().Unix(),
	}
	newTx.AddOutput(coin, address)
	newTx.Finalize()

	return newTx
}

func (tx *Transaction) IsCoinbase() bool {
	return tx.Coinbase
}

// AddInput appends a new input to the transaction based on the previous transaction hash and output index.
func (tx *Transaction) AddInput(prevTxHash []byte, outputIndex int) {
	tx.Inputs = append(tx.Inputs, NewInput(prevTxHash, outputIndex))
}

// AddOutput appends a new output to the transaction with the given value and recipient address.
func (tx *Transaction) AddOutput(value float64, address *rsa.PublicKey) {
	tx.Outputs = append(tx.Outputs, &Output{Value: value, Address: address})
}

func (tx *Transaction) AddMultisigOutput(multisig *Output) {
	tx.Outputs = append(tx.Outputs, &Output{Value: multisig.Value, MultiSigAddresses: multisig.MultiSigAddresses})
}

// RemoveInput removes the input at the specified index, if the index is valid.
func (tx *Transaction) RemoveInput(index int) {
	if index >= 0 && index < len(tx.Inputs) {
		tx.Inputs = append(tx.Inputs[:index], tx.Inputs[index+1:]...)
	}
}

// RemoveInputFromUTXO searches for an input corresponding to the provided UTXO and removes it.
// Note: UTXO and NewUTXO are assumed to be defined elsewhere.
func (tx *Transaction) RemoveInputFromUTXO(ut UTXO) {
	for i, in := range tx.Inputs {
		u := NewUTXO(in.PrevTxHash, in.OutputIndex)
		if u.Equals(&ut) {
			tx.RemoveInput(i)
			return
		}
	}
}

// GetDataToSign returns a byte slice containing the data needed for signing an input.
// It includes the previous transaction hash, the output index of that input,
// and the details of all outputs (value and recipient public key information).
func (tx *Transaction) GetDataToSign(index int) []byte {
	if index >= len(tx.Inputs) {
		return nil
	}

	in := tx.Inputs[index]
	data := make([]byte, 0)

	timestampBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timestampBytes, uint64(tx.Timestamp))
	data = append(data, timestampBytes...)

	data = append(data, in.PrevTxHash...)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(in.OutputIndex))
	data = append(data, buf...)

	for _, op := range tx.Outputs {
		valBuf := make([]byte, 8)
		binary.BigEndian.PutUint64(valBuf, math.Float64bits(op.Value))
		data = append(data, valBuf...)

		if op.MultiSigAddresses != nil && len(op.MultiSigAddresses) > 0 {
			data = append(data, byte(1))
			countBuf := make([]byte, 4)
			binary.BigEndian.PutUint32(countBuf, uint32(len(op.MultiSigAddresses)))
			data = append(data, countBuf...)

			for _, pubKey := range op.MultiSigAddresses {
				expBuf := make([]byte, 4)
				binary.BigEndian.PutUint32(expBuf, uint32(pubKey.E))
				data = append(data, expBuf...)
				data = append(data, pubKey.N.Bytes()...)
			}
		} else if op.Address != nil {
			expBuf := make([]byte, 4)
			binary.BigEndian.PutUint32(expBuf, uint32(op.Address.E))
			data = append(data, expBuf...)
			data = append(data, op.Address.N.Bytes()...)
		} else {
			continue
		}
	}

	return data
}

// AddSignature attaches the given signature to the input at the specified index.
func (tx *Transaction) AddSignature(signature []byte, index int) {
	if index >= 0 && index < len(tx.Inputs) {
		tx.Inputs[index].AddSignature(signature)
	}
}

// GetTx aggregates all transaction data into a single byte slice.
// It concatenates the data from all inputs and outputs, including signatures,
// which is later used to compute the transaction hash.
func (tx *Transaction) GetTx() []byte {
	// Add here timestamp
	data := make([]byte, 0)

	timestampBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timestampBytes, uint64(tx.Timestamp))
	data = append(data, timestampBytes...)

	for _, in := range tx.Inputs {
		data = append(data, in.PrevTxHash...)
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(in.OutputIndex))
		data = append(data, buf...)
		if in.Signature != nil {
			data = append(data, in.Signature...)
		}
	}

	for _, op := range tx.Outputs {
		valBuf := make([]byte, 8)
		binary.BigEndian.PutUint64(valBuf, math.Float64bits(op.Value))
		data = append(data, valBuf...)

		if op.MultiSigAddresses != nil && len(op.MultiSigAddresses) > 0 {
			countBuf := make([]byte, 4)
			binary.BigEndian.PutUint32(countBuf, uint32(len(op.MultiSigAddresses)))
			data = append(data, countBuf...)

			for _, pubKey := range op.MultiSigAddresses {
				expBuf := make([]byte, 4)
				binary.BigEndian.PutUint32(expBuf, uint32((pubKey.E)))
				data = append(data, expBuf...)
				data = append(data, pubKey.N.Bytes()...)
			}
		} else if op.Address != nil {
			expBuf := make([]byte, 4)
			binary.BigEndian.PutUint32(expBuf, uint32(op.Address.E))
			data = append(data, expBuf...)
			data = append(data, op.Address.N.Bytes()...)
		} else {
			continue
		}
	}

	return data
}

// Finalize calculates the SHA-256 hash of the transaction data (from GetTx)
// and sets this hash as the transaction's unique identifier.
func (tx *Transaction) Finalize() {
	hash := sha256.Sum256(tx.GetTx())
	tx.Hash = hash[:]
}

// SetHash assigns a given hash to the transaction after making a copy of it.
func (tx *Transaction) SetHash(h []byte) {
	tx.Hash = make([]byte, len(h))
	copy(tx.Hash, h)
}

// GetHash returns the transaction's hash.
func (tx *Transaction) GetHash() []byte {
	return tx.Hash
}

// GetInputs returns the list of inputs in the transaction.
func (tx *Transaction) GetInputs() []*Input {
	return tx.Inputs
}

// GetOutputs returns the list of outputs in the transaction.
func (tx *Transaction) GetOutputs() []*Output {
	return tx.Outputs
}

// GetInput returns the input at the specified index if it exists, or nil otherwise.
func (tx *Transaction) GetInput(index int) *Input {
	if index >= 0 && index < len(tx.Inputs) {
		return tx.Inputs[index]
	}
	return nil
}

// GetOutput returns the output at the specified index if it exists, or nil otherwise.
func (tx *Transaction) GetOutput(index int) *Output {
	if index >= 0 && index < len(tx.Outputs) {
		return tx.Outputs[index]
	}
	return nil
}

// NumInputs returns the number of inputs in the transaction.
func (tx *Transaction) NumInputs() int {
	return len(tx.Inputs)
}

// NumOutputs returns the number of outputs in the transaction.
func (tx *Transaction) NumOutputs() int {
	return len(tx.Outputs)
}

// Key returns a hexadecimal string representation of the transaction's hash.
// If the hash is not yet computed, it calls Finalize to compute it first.
func (transaction *Transaction) Key() string {
	if transaction.Hash == nil || len(transaction.Hash) == 0 {
		transaction.Finalize()
	}
	return hex.EncodeToString(transaction.Hash)
}

func (tx *Transaction) SignTx(sk *rsa.PrivateKey, input int) {
	dataToSign := tx.GetDataToSign(input)
	hashData1 := sha256.Sum256(dataToSign)
	sig1, err := rsa.SignPKCS1v15(rand.Reader, sk, crypto.SHA256, hashData1[:])
	if err != nil {
		panic(err)
	}
	tx.AddSignature(sig1, input)
	tx.Finalize()
}

func (tx *Transaction) SignMultiSigTx(privKey *rsa.PrivateKey, inputIndex int) {
	dataToSign := tx.GetDataToSign(inputIndex)
	hashData1 := sha256.Sum256(dataToSign)
	sig, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData1[:])
	if err != nil {
		panic(err)
	}
	tx.Inputs[inputIndex].AddMultiSignature(sig)
	tx.Finalize()
}
