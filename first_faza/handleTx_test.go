package first_faza

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleTxIsValid_withValidTransaction(t *testing.T) {
	utxoPool = NewUTXOPool()

	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey := &privKey.PublicKey

	tx0 := NewTransaction()
	tx0.AddOutput(10.0, pubKey)
	tx0.Finalize()

	utxo0 := NewUTXO(tx0.GetHash(), 0)
	utxoPool.AddUTXO(*utxo0, *tx0.GetOutput(0))

	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddOutput(5.0, pubKey)
	tx1.AddOutput(4.0, pubKey)

	dataToSign := tx1.GetDataToSign(0)
	hashData := sha256.Sum256(dataToSign)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData[:])
	if err != nil {
		t.Fatal(err)
	}

	tx1.AddSignature(signature, 0)
	tx1.Finalize()

	HandleTxs(utxoPool)
	valid := TxIsValid(*tx1, utxoPool)
	assert.True(t, valid, "Transaction should be valid")
}

func TestHandleTxIsValid_withIncorrectSignature(t *testing.T) {
	utxoPool = NewUTXOPool()

	privKey1, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey1 := &privKey1.PublicKey

	privKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey2 := &privKey2.PublicKey

	tx0 := NewTransaction()
	tx0.AddOutput(10.0, pubKey1)
	tx0.Finalize()

	utxo0 := NewUTXO(tx0.GetHash(), 0)
	utxoPool.AddUTXO(*utxo0, *tx0.GetOutput(0))

	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddOutput(5.0, pubKey1)
	tx1.AddOutput(5.0, pubKey2)
	dataToSign1 := tx1.GetDataToSign(0)
	hashData1 := sha256.Sum256(dataToSign1)
	validSignature1, err := rsa.SignPKCS1v15(rand.Reader, privKey1, crypto.SHA256, hashData1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(validSignature1, 0)

	HandleTxs(utxoPool)
	valid := TxIsValid(*tx1, utxoPool)
	assert.True(t, valid, "Transaction should be valid.")

	utxo1a := NewUTXO(tx1.GetHash(), 0)
	utxo1b := NewUTXO(tx1.GetHash(), 1)
	utxoPool.AddUTXO(*utxo1a, *tx1.GetOutput(0))
	utxoPool.AddUTXO(*utxo1b, *tx1.GetOutput(1))

	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 0)
	tx2.AddOutput(5.0, pubKey1)

	dataToSign := tx2.GetDataToSign(0)
	hashData := sha256.Sum256(dataToSign)
	badSignature, err := rsa.SignPKCS1v15(rand.Reader, privKey2, crypto.SHA256, hashData[:])
	if err != nil {
		t.Fatal(err)
	}

	for i := range badSignature {
		badSignature[i] = badSignature[i] & 1
	}

	tx2.AddSignature(badSignature, 0)
	tx2.Finalize()

	HandleTxs(utxoPool)
	valid = TxIsValid(*tx2, utxoPool)
	assert.False(t, valid, "Transaction should be invalid due to an incorrect signature.")
}

func TestHandleTxIsValid_withSignatureFromInvalidPrivateKey(t *testing.T) {
	utxoPool = NewUTXOPool()

	// Generate a valid key pair.
	validPrivKey1, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	validPubKey1 := &validPrivKey1.PublicKey

	validPrivKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	validPubKey2 := &validPrivKey2.PublicKey

	tx0 := NewTransaction()
	tx0.AddOutput(10.0, validPubKey1)
	tx0.Finalize()

	utxo0 := NewUTXO(tx0.GetHash(), 0)
	utxoPool.AddUTXO(*utxo0, *tx0.GetOutput(0))

	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddOutput(5.0, validPubKey1)
	tx1.AddOutput(5.0, validPubKey2)

	dataToSign := tx1.GetDataToSign(0)
	hashData := sha256.Sum256(dataToSign)
	validSignature, err := rsa.SignPKCS1v15(rand.Reader, validPrivKey1, crypto.SHA256, hashData[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(validSignature, 0)
	tx1.Finalize()

	HandleTxs(utxoPool)
	valid := TxIsValid(*tx1, utxoPool)
	assert.True(t, valid, "Transaction 1 should be valid.")

	utxo1a := NewUTXO(tx1.GetHash(), 0)
	utxo1b := NewUTXO(tx1.GetHash(), 1)
	utxoPool.AddUTXO(*utxo1a, *tx1.GetOutput(0))
	utxoPool.AddUTXO(*utxo1b, *tx1.GetOutput(1))

	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 1)
	tx2.AddOutput(5.0, validPubKey1)
	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)
	badSignature, err := rsa.SignPKCS1v15(rand.Reader, validPrivKey1, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(badSignature, 0)
	tx2.Finalize()

	HandleTxs(utxoPool)
	valid = TxIsValid(*tx2, utxoPool)
	assert.False(t, valid, "Transaction 2 should be invalid when signed with an incorrect private key.")
}

func TestHandleTxIsValid_withOutputExceedingInputs(t *testing.T) {
	utxoPool = NewUTXOPool()

	privKey1, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey1 := &privKey1.PublicKey

	privKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey2 := &privKey2.PublicKey

	tx0 := NewTransaction()
	tx0.AddOutput(10.0, pubKey1)
	tx0.Finalize()

	utxo0 := NewUTXO(tx0.GetHash(), 0)
	utxoPool.AddUTXO(*utxo0, *tx0.GetOutput(0))

	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddOutput(5.0, pubKey1)
	tx1.AddOutput(5.0, pubKey2)

	dataToSign := tx1.GetDataToSign(0)
	hashData := sha256.Sum256(dataToSign)
	validSignature, err := rsa.SignPKCS1v15(rand.Reader, privKey1, crypto.SHA256, hashData[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(validSignature, 0)
	tx1.Finalize()

	HandleTxs(utxoPool)
	valid := TxIsValid(*tx1, utxoPool)
	assert.True(t, valid, "Transaction 1 should be valid.")

	utxo1a := NewUTXO(tx1.GetHash(), 0)
	utxo1b := NewUTXO(tx1.GetHash(), 1)
	utxoPool.AddUTXO(*utxo1a, *tx1.GetOutput(0))
	utxoPool.AddUTXO(*utxo1b, *tx1.GetOutput(1))

	tx2 := NewTransaction()
	tx2.AddInput(tx0.GetHash(), 1)
	tx2.AddOutput(10.0, pubKey1)
	tx2.AddOutput(14.0, pubKey2)

	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)
	validSignature2, err := rsa.SignPKCS1v15(rand.Reader, privKey2, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(validSignature2, 0)
	tx2.Finalize()

	valid = TxIsValid(*tx2, utxoPool)
	assert.False(t, valid, "Transaction 2 should be invalid because transaction with output exceeding inputs")
}

func TestHandleTxIsValid_withInputsOutsideUTXOPool(t *testing.T) {
	utxoPool = NewUTXOPool()

	privKey1, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey1 := &privKey1.PublicKey

	privKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey2 := &privKey2.PublicKey

	tx0 := NewTransaction()
	tx0.AddOutput(10.0, pubKey1)
	tx0.Finalize()

	utxo0 := NewUTXO(tx0.GetHash(), 0)
	utxoPool.AddUTXO(*utxo0, *tx0.GetOutput(0))

	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)

	tx1.AddOutput(5.0, pubKey1)
	tx1.AddOutput(4.0, pubKey2)

	dataToSign := tx1.GetDataToSign(0)
	hashData := sha256.Sum256(dataToSign)
	validSignature, err := rsa.SignPKCS1v15(rand.Reader, privKey1, crypto.SHA256, hashData[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(validSignature, 0)
	tx1.Finalize()

	HandleTxs(utxoPool)
	valid := TxIsValid(*tx1, utxoPool)
	assert.True(t, valid, "Transaction should be valid.")

	utxo1a := NewUTXO(tx1.GetHash(), 0)
	utxo1b := NewUTXO(tx1.GetHash(), 1)
	utxoPool.AddUTXO(*utxo1a, *tx1.GetOutput(0))
	utxoPool.AddUTXO(*utxo1b, *tx1.GetOutput(1))

	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 1)

	tx2.AddOutput(2.0, pubKey1)
	tx2.AddOutput(2.0, pubKey2)

	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)
	validSignature2, err := rsa.SignPKCS1v15(rand.Reader, privKey2, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(validSignature2, 0)
	tx2.Finalize()

	HandleTxs(utxoPool)
	valid = TxIsValid(*tx2, utxoPool)
	assert.True(t, valid, "Transaction should be valid.")

	utxo2a := NewUTXO(tx2.GetHash(), 0)
	utxo2b := NewUTXO(tx2.GetHash(), 1)
	utxoPool.AddUTXO(*utxo2a, *tx2.GetOutput(0))
	utxoPool.AddUTXO(*utxo2b, *tx2.GetOutput(1))

	// Create a fake hash that does not exist in the UTXO pool.
	fakeHash := make([]byte, 32)
	for i := range fakeHash {
		fakeHash[i] = byte(i + 1)
	}

	tx3 := NewTransaction()
	tx3.AddInput(fakeHash, 0)

	tx3.AddOutput(1.0, pubKey1)
	tx3.AddOutput(1.0, pubKey2)

	dataToSign3 := tx3.GetDataToSign(0)
	hashData3 := sha256.Sum256(dataToSign3)
	validSignature3, err := rsa.SignPKCS1v15(rand.Reader, privKey1, crypto.SHA256, hashData3[:])
	if err != nil {
		t.Fatal(err)
	}
	tx3.AddSignature(validSignature3, 0)
	tx3.Finalize()

	HandleTxs(utxoPool)
	valid = TxIsValid(*tx3, utxoPool)
	assert.False(t, valid, "Transaction should be invalid because it references an input outside the current UTXOPool.")
}

func TestHandleTxIsValid_withDuplicateUTXO(t *testing.T) {
	utxoPool = NewUTXOPool()

	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey := &privKey.PublicKey

	tx0 := NewTransaction()
	tx0.AddOutput(10.0, pubKey)
	tx0.Finalize()

	utxo0 := NewUTXO(tx0.GetHash(), 0)
	utxoPool.AddUTXO(*utxo0, *tx0.GetOutput(0))

	tx1 := NewTransaction()
	// Add the same UTXO input twice.
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddInput(tx0.GetHash(), 0)

	tx1.AddOutput(5.0, pubKey)
	tx1.AddOutput(4.0, pubKey)

	dataToSign0 := tx1.GetDataToSign(0)
	hashData0 := sha256.Sum256(dataToSign0)
	validSignature0, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData0[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(validSignature0, 0)

	dataToSign1 := tx1.GetDataToSign(1)
	hashData1 := sha256.Sum256(dataToSign1)
	validSignature1, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(validSignature1, 1)

	tx1.Finalize()

	HandleTxs(utxoPool)
	valid := TxIsValid(*tx1, utxoPool)
	assert.False(t, valid, "Transaction should be invalid if the same UTXO is declared more than once.")

	tx2 := NewTransaction()
	// Add the same UTXO input twice.
	tx2.AddInput(tx0.GetHash(), 0)
	tx2.AddInput(tx0.GetHash(), 0)

	tx2.AddOutput(5.0, pubKey)
	tx2.AddOutput(4.0, pubKey)

	dataToSign2a := tx2.GetDataToSign(0)
	hashData2a := sha256.Sum256(dataToSign2a)
	validSignature2a, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData2a[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(validSignature2a, 0)

	dataToSign2b := tx2.GetDataToSign(1)
	hashData2b := sha256.Sum256(dataToSign2b)
	validSignature2b, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData2b[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(validSignature2b, 1)

	tx2.Finalize()

	HandleTxs(utxoPool)
	valid = TxIsValid(*tx2, utxoPool)
	assert.False(t, valid, "Transaction should be invalid if the same UTXO is declared more than once.")
}

func TestHandleTxIsValid_withNegativeOutputValue(t *testing.T) {
	utxoPool = NewUTXOPool()

	privKey1, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey1 := &privKey1.PublicKey

	privKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey2 := &privKey2.PublicKey

	tx0 := NewTransaction()
	tx0.AddOutput(10.0, pubKey2)
	tx0.Finalize()

	utxo0 := NewUTXO(tx0.GetHash(), 0)
	utxoPool.AddUTXO(*utxo0, *tx0.GetOutput(0))

	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddOutput(5.0, pubKey1)
	tx1.AddOutput(4.0, pubKey2)

	dataToSign := tx1.GetDataToSign(0)
	hashData := sha256.Sum256(dataToSign)
	validSignature, err := rsa.SignPKCS1v15(rand.Reader, privKey2, crypto.SHA256, hashData[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(validSignature, 0)
	tx1.Finalize()

	HandleTxs(utxoPool)
	valid := TxIsValid(*tx1, utxoPool)
	assert.True(t, valid, "Transaction should be valid.")

	utxo1a := NewUTXO(tx1.GetHash(), 0)
	utxo1b := NewUTXO(tx1.GetHash(), 1)
	utxoPool.AddUTXO(*utxo1a, *tx1.GetOutput(0))
	utxoPool.AddUTXO(*utxo1b, *tx1.GetOutput(1))

	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 0)
	tx2.AddOutput(2.5, pubKey1)
	tx2.AddOutput(-1.5, pubKey2)

	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)
	validSignature2, err := rsa.SignPKCS1v15(rand.Reader, privKey1, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(validSignature2, 0)
	tx2.Finalize()

	HandleTxs(utxoPool)
	valid = TxIsValid(*tx2, utxoPool)

	assert.False(t, valid, "Transaction should be invalid because it contains a negative output value.")
}

func TestHandleTxs_withSimpleValidTransaction(t *testing.T) {
	utxoPool = NewUTXOPool()

	// Generate keys for signing.
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey := &privKey.PublicKey

	// Generate a second key pair for sending funds to a different address.
	privKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey2 := &privKey2.PublicKey

	// Create an initial transaction (tx0).
	tx0 := NewTransaction()
	tx0.AddOutput(10.0, pubKey)
	tx0.Finalize()

	utxo0 := NewUTXO(tx0.GetHash(), 0)
	utxoPool.AddUTXO(*utxo0, *tx0.GetOutput(0))

	// Create a valid transaction (tx1) spending tx0's output.
	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddOutput(5.0, pubKey)
	tx1.AddOutput(4.0, pubKey2)

	dataToSign1 := tx1.GetDataToSign(0)
	hashData1 := sha256.Sum256(dataToSign1)
	sig1, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(sig1, 0)
	tx1.Finalize()

	utxo1a := NewUTXO(tx1.GetHash(), 0)
	utxo1b := NewUTXO(tx1.GetHash(), 1)
	utxoPool.AddUTXO(*utxo1a, *tx1.GetOutput(0))
	utxoPool.AddUTXO(*utxo1b, *tx1.GetOutput(1))

	// Create another valid transaction (tx2) spending tx1's second output.
	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 1) // Uses the second output from tx1.
	tx2.AddOutput(3.0, pubKey2)
	tx2.AddOutput(1.0, pubKey)

	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)
	sig2, err := rsa.SignPKCS1v15(rand.Reader, privKey2, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(sig2, 0)
	tx2.Finalize()

	// Process the transactions.
	HandleTxs(utxoPool)
	possibleTxs := []Transaction{*tx1, *tx2}
	handler := Handler(possibleTxs)

	// Check if only one transaction is processed immediately.
	assert.NotNil(t, handler, "Handler result should not be nil.")
	assert.Equal(t, 2, len(handler), "Two transaction should be processed.")

	// Verify that tx1 and tx2 was processed.
	containsTx1 := false
	containsTx2 := false
	for _, tx := range handler {
		if tx.Key() == tx1.Key() {
			containsTx1 = true
		} else if tx.Key() == tx2.Key() {
			containsTx2 = true
		}
	}
	assert.True(t, containsTx1, "Transaction tx1 should be processed.")
	assert.True(t, containsTx2, "Transaction tx2 should be processed.")
}

func TestHandleTxs_withSomeInvalidTransactionDueToInvalidSignatures(t *testing.T) {
	utxoPool = NewUTXOPool()

	// Generate keys for signing.
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey := &privKey.PublicKey

	// Generate a second key pair for sending funds to a different address.
	privKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey2 := &privKey2.PublicKey

	privKey3, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey3 := &privKey3.PublicKey

	tx0 := NewTransaction()
	tx0.AddOutput(10.0, pubKey)
	tx0.AddOutput(15.0, pubKey)
	tx0.Finalize()

	utxo0a := NewUTXO(tx0.GetHash(), 0)
	utxo0b := NewUTXO(tx0.GetHash(), 1)
	utxoPool.AddUTXO(*utxo0a, *tx0.GetOutput(0))
	utxoPool.AddUTXO(*utxo0b, *tx0.GetOutput(1))

	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)

	tx1.AddOutput(8.0, pubKey2)
	tx1.AddOutput(2.0, pubKey)
	dataToSign1 := tx1.GetDataToSign(0)
	hashData1 := sha256.Sum256(dataToSign1)
	sig1, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(sig1, 0)
	tx1.Finalize()

	utxo1a := NewUTXO(tx1.GetHash(), 0)
	utxo1b := NewUTXO(tx1.GetHash(), 1)
	utxoPool.AddUTXO(*utxo1a, *tx1.GetOutput(0))
	utxoPool.AddUTXO(*utxo1b, *tx1.GetOutput(1))

	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 1)
	tx2.AddOutput(2.0, pubKey2)
	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)

	sig2, err := rsa.SignPKCS1v15(rand.Reader, privKey3, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}

	tx2.AddSignature(sig2, 0)
	tx2.Finalize()

	tx3 := NewTransaction()
	tx3.AddInput(tx0.GetHash(), 1)
	tx3.AddOutput(14.0, pubKey2)
	tx3.AddOutput(1.0, pubKey3)
	dataToSign3 := tx3.GetDataToSign(0)
	hashData3 := sha256.Sum256(dataToSign3)

	sig3, err := rsa.SignPKCS1v15(rand.Reader, privKey2, crypto.SHA256, hashData3[:])
	if err != nil {
		t.Fatal(err)
	}

	tx3.AddSignature(sig3, 0)
	tx3.Finalize()

	HandleTxs(utxoPool)
	possibleTxs := []Transaction{*tx1, *tx2, *tx3}
	handler := Handler(possibleTxs)

	assert.NotNil(t, handler, "Handler result should not be nil.")
	assert.Equal(t, 1, len(handler), "One transaction should be processed.")

	containsTx1 := false
	for _, tx := range handler {
		if tx.Key() == tx1.Key() {
			containsTx1 = true
		}
	}
	assert.True(t, containsTx1, "Transaction tx1 should be immediately processed.")
}

func TestHandleTxs_withInvalidTransactionsDueToInputLessThanOutput(t *testing.T) {
	utxoPool = NewUTXOPool()

	// Generate keys.
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey := &privKey.PublicKey

	privKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey2 := &privKey2.PublicKey

	privKey3, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey3 := &privKey3.PublicKey

	// Create an initial funding transaction (tx0) with two outputs.
	tx0 := NewTransaction()
	tx0.AddOutput(10.0, pubKey) // Output 0: 10.0
	tx0.AddOutput(15.0, pubKey) // Output 1: 15.0
	tx0.Finalize()

	utxo0a := NewUTXO(tx0.GetHash(), 0)
	utxo0b := NewUTXO(tx0.GetHash(), 1)
	utxoPool.AddUTXO(*utxo0a, *tx0.GetOutput(0))
	utxoPool.AddUTXO(*utxo0b, *tx0.GetOutput(1))

	// tx1: Invalid because it spends 10.0 but outputs 11.0.
	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddOutput(11.0, pubKey2) // Output of 11.0 > input of 10.0.
	dataToSign1 := tx1.GetDataToSign(0)
	hashData1 := sha256.Sum256(dataToSign1)
	sig1, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(sig1, 0)
	tx1.Finalize()

	// tx2: Valid because it spends 15.0 and outputs 12.0 (fee 3.0).
	tx2 := NewTransaction()
	tx2.AddInput(tx0.GetHash(), 1)
	tx2.AddOutput(12.0, pubKey2)
	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)
	sig2, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(sig2, 0)
	tx2.Finalize()

	utxo2a := NewUTXO(tx2.GetHash(), 0)
	utxoPool.AddUTXO(*utxo2a, *tx2.GetOutput(0))

	// tx3: Invalid because it spends both outputs (10.0 + 15.0 = 25.0)
	// but outputs sum to 30.0 (20.0 + 10.0).
	tx3 := NewTransaction()
	tx3.AddInput(tx0.GetHash(), 0)
	tx3.AddInput(tx0.GetHash(), 1)
	tx3.AddOutput(20.0, pubKey2)
	tx3.AddOutput(10.0, pubKey3)
	// Sign each input.
	dataToSign3_0 := tx3.GetDataToSign(0)
	hashData3_0 := sha256.Sum256(dataToSign3_0)
	sig3_0, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData3_0[:])
	if err != nil {
		t.Fatal(err)
	}
	tx3.AddSignature(sig3_0, 0)
	dataToSign3_1 := tx3.GetDataToSign(1)
	hashData3_1 := sha256.Sum256(dataToSign3_1)
	sig3_1, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData3_1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx3.AddSignature(sig3_1, 1)
	tx3.Finalize()

	tx4 := NewTransaction()
	tx4.AddInput(tx2.GetHash(), 0)
	tx4.AddOutput(13.0, pubKey2)
	dataToSign4 := tx4.GetDataToSign(0)
	hashData4 := sha256.Sum256(dataToSign4)
	sig4, err := rsa.SignPKCS1v15(rand.Reader, privKey2, crypto.SHA256, hashData4[:])
	if err != nil {
		t.Fatal(err)
	}
	tx4.AddSignature(sig4, 0)
	tx4.Finalize()

	// Process the transactions.
	// Only tx2 is valid (its input covers its output), while tx1 and tx3 are invalid.
	HandleTxs(utxoPool)
	possibleTxs := []Transaction{*tx1, *tx2, *tx3, *tx4}
	handler := Handler(possibleTxs)

	// Expect only tx2 to be processed immediately.
	assert.NotNil(t, handler, "Handler result should not be nil.")
	assert.Equal(t, 1, len(handler), "Only one transaction should be processed.")

	containsTx2 := false
	for _, tx := range handler {
		if tx.Key() == tx2.Key() {
			containsTx2 = true
		}
	}
	assert.True(t, containsTx2, "Transaction tx2 should be processed.")
}

func TestHandleTxs_withDoubleSpends(t *testing.T) {
	utxoPool = NewUTXOPool()

	// Generate four key pairs.
	privKey1, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey1 := &privKey1.PublicKey

	privKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey2 := &privKey2.PublicKey

	privKey3, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey3 := &privKey3.PublicKey

	privKey4, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey4 := &privKey4.PublicKey

	// Create funding transaction (tx0) with two outputs:
	// - Output 0: 10.0 funds to pubKey1.
	// - Output 1: 15.0 funds to pubKey1.
	tx0 := NewTransaction()
	tx0.AddOutput(10.0, pubKey1) // UTXO0
	tx0.AddOutput(15.0, pubKey1) // UTXO1
	tx0.Finalize()

	utxo0a := NewUTXO(tx0.GetHash(), 0)
	utxo0b := NewUTXO(tx0.GetHash(), 1)
	utxoPool.AddUTXO(*utxo0a, *tx0.GetOutput(0))
	utxoPool.AddUTXO(*utxo0b, *tx0.GetOutput(1))

	// tx1: Valid transaction spending output 0 (10.0).
	// It sends 9.0 funds to pubKey2.
	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddOutput(9.0, pubKey2)
	dataToSign1 := tx1.GetDataToSign(0)
	hashData1 := sha256.Sum256(dataToSign1)
	sig1, err := rsa.SignPKCS1v15(rand.Reader, privKey1, crypto.SHA256, hashData1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(sig1, 0)
	tx1.Finalize()

	utxo1a := NewUTXO(tx1.GetHash(), 0)
	utxoPool.AddUTXO(*utxo1a, *tx1.GetOutput(0))

	// tx2: Another transaction attempting to spend the same UTXO (output 0).
	// It is almost identical but signed with an incorrect key (privKey3) so it should be invalid.
	tx2 := NewTransaction()
	tx2.AddInput(tx0.GetHash(), 0)
	tx2.AddOutput(8.5, pubKey4)
	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)
	sig2, err := rsa.SignPKCS1v15(rand.Reader, privKey1, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(sig2, 0)
	tx2.Finalize()

	// tx3: Valid transaction spending output 1 (15.0).
	// It sends 12.0 funds to pubKey2.
	tx3 := NewTransaction()
	tx3.AddInput(tx0.GetHash(), 1)
	tx3.AddOutput(12.0, pubKey2)
	tx3.AddOutput(3.0, pubKey3)
	dataToSign3 := tx3.GetDataToSign(0)
	hashData3 := sha256.Sum256(dataToSign3)
	sig3, err := rsa.SignPKCS1v15(rand.Reader, privKey1, crypto.SHA256, hashData3[:])
	if err != nil {
		t.Fatal(err)
	}
	tx3.AddSignature(sig3, 0)
	tx3.Finalize()

	utxo3a := NewUTXO(tx3.GetHash(), 0)
	utxo3b := NewUTXO(tx3.GetHash(), 1)
	utxoPool.AddUTXO(*utxo3a, *tx3.GetOutput(0))
	utxoPool.AddUTXO(*utxo3b, *tx3.GetOutput(1))

	// Full ok
	tx4 := NewTransaction()
	tx4.AddInput(tx3.GetHash(), 0) // priv2
	tx4.AddOutput(9.0, pubKey3)
	tx4.AddOutput(3.0, pubKey4)
	dataToSign4 := tx4.GetDataToSign(0)
	hashData4 := sha256.Sum256(dataToSign4)
	sig4, err := rsa.SignPKCS1v15(rand.Reader, privKey2, crypto.SHA256, hashData4[:])
	if err != nil {
		t.Fatal(err)
	}
	tx4.AddSignature(sig4, 0)
	tx4.Finalize()

	utxo4a := NewUTXO(tx4.GetHash(), 0)
	utxo4b := NewUTXO(tx4.GetHash(), 1)
	utxoPool.AddUTXO(*utxo4a, *tx4.GetOutput(0))
	utxoPool.AddUTXO(*utxo4b, *tx4.GetOutput(1))

	tx5 := NewTransaction()
	tx5.AddInput(tx3.GetHash(), 0)
	tx5.AddOutput(9.0, pubKey3)
	tx5.AddOutput(3.0, pubKey4)
	dataToSign5 := tx5.GetDataToSign(0)
	hashData5 := sha256.Sum256(dataToSign5)
	sig5, err := rsa.SignPKCS1v15(rand.Reader, privKey2, crypto.SHA256, hashData5[:])
	if err != nil {
		t.Fatal(err)
	}
	tx5.AddSignature(sig5, 0)
	tx5.Finalize()

	// Prepare the list of possible transactions.
	// tx1 and tx2 both try to spend UTXO0; tx3 spends UTXO1.
	possibleTxs := []Transaction{*tx1, *tx2, *tx3, *tx4, *tx5}
	handler := Handler(possibleTxs)

	// Expect 2 transactions to be processed immediately:
	// one for UTXO0 (either tx1 or tx2, but tx1 is valid) and tx3 for UTXO1.
	assert.NotNil(t, handler, "Handler result should not be nil.")
	assert.Equal(t, 3, len(handler), "Three transactions should be processed.")

	// Check that the accepted transactions are tx1 and tx3.
	acceptedKeys := map[string]bool{
		tx1.Key(): true,
		tx3.Key(): true,
		tx4.Key(): true,
	}
	for _, tx := range handler {
		assert.True(t, acceptedKeys[tx.Key()], "Accepted transaction should be either tx1 or tx3 or tx4.")
	}
}

func TestHandleTxs_withDependentAndSimpleValidTransaction(t *testing.T) {
	// Reset global UTXO pool.
	utxoPool = NewUTXOPool()

	// Generate a key pair for funding/spending.
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey := &privKey.PublicKey

	// Generate a second key pair (could be used for sending funds to another address).
	privKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey2 := &privKey2.PublicKey

	// Create an initial funding transaction (tx0) with two outputs.
	// UTXO0: 10.0 funds to pubKey.
	// UTXO1: 15.0 funds to pubKey.
	tx0 := NewTransaction()
	tx0.AddOutput(10.0, pubKey)
	tx0.AddOutput(15.0, pubKey)
	tx0.Finalize()

	utxo0 := NewUTXO(tx0.GetHash(), 0)
	utxo1 := NewUTXO(tx0.GetHash(), 1)
	utxoPool.AddUTXO(*utxo0, *tx0.GetOutput(0))
	utxoPool.AddUTXO(*utxo1, *tx0.GetOutput(1))

	// tx1: Independent valid transaction spending UTXO0.
	// It sends 9.0 funds to pubKey2.
	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddOutput(9.0, pubKey2)
	dataToSign1 := tx1.GetDataToSign(0)
	hashData1 := sha256.Sum256(dataToSign1)
	sig1, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(sig1, 0)
	tx1.Finalize()

	// tx2: Dependent transaction that spends tx1's output.
	// It sends, say, 8.0 funds to pubKey2.
	// Because tx1's output isn't in the original UTXO snapshot,
	// tx2 is not valid immediately.
	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 0) // Depends on tx1's output.
	tx2.AddOutput(8.0, pubKey2)
	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)
	sig2, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(sig2, 0)
	tx2.Finalize()

	// tx3: Independent valid transaction spending UTXO1.
	// It sends 12.0 funds to pubKey2.
	tx3 := NewTransaction()
	tx3.AddInput(tx0.GetHash(), 1)
	tx3.AddOutput(12.0, pubKey2)
	dataToSign3 := tx3.GetDataToSign(0)
	hashData3 := sha256.Sum256(dataToSign3)
	sig3, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData3[:])
	if err != nil {
		t.Fatal(err)
	}
	tx3.AddSignature(sig3, 0)
	tx3.Finalize()

	// Submit all transactions together.
	// tx1 and tx3 are valid with respect to the original UTXO pool.
	// tx2, although individually valid, is dependent on tx1's output and should be deferred.
	possibleTxs := []Transaction{*tx1, *tx2, *tx3}
	handler := Handler(possibleTxs)

	// Expect 2 transactions to be processed immediately: tx1 (from UTXO0) and tx3 (from UTXO1).
	assert.NotNil(t, handler, "Handler result should not be nil.")
	assert.Equal(t, 2, len(handler), "Two transactions should be processed immediately.")

	// Verify that tx1 and tx3 were processed by comparing their keys.
	acceptedKeys := map[string]bool{
		tx1.Key(): true,
		tx3.Key(): true,
	}
	for _, tx := range handler {
		assert.True(t, acceptedKeys[tx.Key()], "Accepted transaction should be either tx1 or tx3.")
	}
}

func TestHandleTxs_withNonExistingUTXOInputs(t *testing.T) {
	utxoPool = NewUTXOPool()

	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey := &privKey.PublicKey

	privKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey2 := &privKey2.PublicKey

	// Create a funding transaction (tx0) with two outputs.
	tx0 := NewTransaction()
	tx0.AddOutput(20.0, pubKey) // UTXO0: 20.0 funds.
	tx0.AddOutput(30.0, pubKey) // UTXO1: 30.0 funds.
	tx0.Finalize()

	utxo0 := NewUTXO(tx0.GetHash(), 0)
	utxo1 := NewUTXO(tx0.GetHash(), 1)
	utxoPool.AddUTXO(*utxo0, *tx0.GetOutput(0))
	utxoPool.AddUTXO(*utxo1, *tx0.GetOutput(1))

	// tx1: Valid transaction spending UTXO0.
	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	// For example, spending 20.0 to create an output of 18.0 (fee of 2.0).
	tx1.AddOutput(18.0, pubKey2)
	dataToSign1 := tx1.GetDataToSign(0)
	hashData1 := sha256.Sum256(dataToSign1)
	sig1, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(sig1, 0)
	tx1.Finalize()

	// tx2: Invalid transaction referencing a non-existing UTXO (fake hash).
	tx2 := NewTransaction()
	fakeHash := make([]byte, 32)
	for i := range fakeHash {
		fakeHash[i] = 99
	}
	tx2.AddInput(fakeHash, 0)
	tx2.AddOutput(15.0, pubKey2)
	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)
	sig2, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(sig2, 0)
	tx2.Finalize()

	// tx3: Transaction with two inputs: one valid from UTXO1 and one fake.
	tx3 := NewTransaction()
	tx3.AddInput(tx0.GetHash(), 1) // valid input from tx0 output1.
	tx3.AddInput(fakeHash, 0)      // fake input.
	tx3.AddOutput(25.0, pubKey2)   // Total output (25.0) is less than UTXO1's 30.0, but invalid due to fake input.
	// Sign each input separately.
	dataToSign3_0 := tx3.GetDataToSign(0)
	hashData3_0 := sha256.Sum256(dataToSign3_0)
	sig3_0, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData3_0[:])
	if err != nil {
		t.Fatal(err)
	}
	tx3.AddSignature(sig3_0, 0)
	dataToSign3_1 := tx3.GetDataToSign(1)
	hashData3_1 := sha256.Sum256(dataToSign3_1)
	sig3_1, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData3_1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx3.AddSignature(sig3_1, 1)
	tx3.Finalize()

	// tx4: Transaction referencing a valid tx0 hash but with a non-existing output index (e.g., index 2).
	tx4 := NewTransaction()
	tx4.AddInput(tx0.GetHash(), 2) // tx0 has outputs only at indices 0 and 1.
	tx4.AddOutput(10.0, pubKey2)
	dataToSign4 := tx4.GetDataToSign(0)
	hashData4 := sha256.Sum256(dataToSign4)
	sig4, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData4[:])
	if err != nil {
		t.Fatal(err)
	}
	tx4.AddSignature(sig4, 0)
	tx4.Finalize()

	// Submit all transactions together.
	possibleTxs := []Transaction{*tx1, *tx2, *tx3, *tx4}
	handler := Handler(possibleTxs)

	// Expect only tx1 to be processed immediately,
	// since tx2, tx3, and tx4 contain inputs that do not exist in the original UTXO pool snapshot.
	assert.NotNil(t, handler, "Handler result should not be nil.")
	assert.Equal(t, 1, len(handler), "Only one transaction should be processed immediately (tx1).")

	containsTx1 := false
	for _, tx := range handler {
		if tx.Key() == tx1.Key() {
			containsTx1 = true
		}
	}
	assert.True(t, containsTx1, "Transaction tx1 should be immediately processed.")
}

func TestHandleTxs_withComplexTransactions(t *testing.T) {
	// Reset the global UTXO pool.
	utxoPool = NewUTXOPool()

	// Generate several key pairs.
	key1, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey1 := &key1.PublicKey

	key2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey2 := &key2.PublicKey

	key3, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey3 := &key3.PublicKey

	key4, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey4 := &key4.PublicKey

	key5, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey5 := &key5.PublicKey

	// ----------------------------------
	// Create funding transaction (tx0) with three outputs.
	// O0: 20.0 funds to pubKey1.
	// O1: 30.0 funds to pubKey2.
	// O2: 50.0 funds to pubKey1.
	tx0 := NewTransaction()
	tx0.AddOutput(20.0, pubKey1) // UTXO0
	tx0.AddOutput(30.0, pubKey2) // UTXO1
	tx0.AddOutput(50.0, pubKey1) // UTXO2
	tx0.Finalize()

	utxo0 := NewUTXO(tx0.GetHash(), 0)
	utxo1 := NewUTXO(tx0.GetHash(), 1)
	utxo2 := NewUTXO(tx0.GetHash(), 2)
	utxoPool.AddUTXO(*utxo0, *tx0.GetOutput(0))
	utxoPool.AddUTXO(*utxo1, *tx0.GetOutput(1))
	utxoPool.AddUTXO(*utxo2, *tx0.GetOutput(2))

	// ----------------------------------
	// tx1: Valid independent transaction spending O0.
	// Spends 20.0, outputs 18.0 to pubKey3 (fee 2.0).
	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddOutput(18.0, pubKey3)
	dataToSign1 := tx1.GetDataToSign(0)
	hashData1 := sha256.Sum256(dataToSign1)
	sig1, err := rsa.SignPKCS1v15(rand.Reader, key1, crypto.SHA256, hashData1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(sig1, 0)
	tx1.Finalize()

	// ----------------------------------
	// tx2: Attempt to double spend O0.
	// Spends 20.0, outputs 15.0 to pubKey4.
	// Even if signature is valid, double spending makes it invalid.
	tx2 := NewTransaction()
	tx2.AddInput(tx0.GetHash(), 0)
	tx2.AddOutput(15.0, pubKey4)
	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)
	sig2, err := rsa.SignPKCS1v15(rand.Reader, key1, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(sig2, 0)
	tx2.Finalize()

	// ----------------------------------
	// tx3: Valid independent transaction spending O1.
	// Spends 30.0, outputs 28.0 to pubKey3 (fee 2.0).
	tx3 := NewTransaction()
	tx3.AddInput(tx0.GetHash(), 1)
	tx3.AddOutput(28.0, pubKey3)
	dataToSign3 := tx3.GetDataToSign(0)
	hashData3 := sha256.Sum256(dataToSign3)
	sig3, err := rsa.SignPKCS1v15(rand.Reader, key2, crypto.SHA256, hashData3[:])
	if err != nil {
		t.Fatal(err)
	}
	tx3.AddSignature(sig3, 0)
	tx3.Finalize()

	// ----------------------------------
	// tx4: Transaction with two inputs:
	// One valid from O2 (50.0 funds) and one fake input.
	// Outputs 40.0 to pubKey5.
	// This should be invalid due to the fake input.
	tx4 := NewTransaction()
	tx4.AddInput(tx0.GetHash(), 2) // valid input from O2.
	// Create a fake UTXO input.
	fakeHash := make([]byte, 32)
	for i := range fakeHash {
		fakeHash[i] = 77
	}
	tx4.AddInput(fakeHash, 0) // fake input.
	tx4.AddOutput(40.0, pubKey5)
	// Sign each input.
	dataToSign4_0 := tx4.GetDataToSign(0)
	hashData4_0 := sha256.Sum256(dataToSign4_0)
	sig4_0, err := rsa.SignPKCS1v15(rand.Reader, key1, crypto.SHA256, hashData4_0[:])
	if err != nil {
		t.Fatal(err)
	}
	tx4.AddSignature(sig4_0, 0)
	dataToSign4_1 := tx4.GetDataToSign(1)
	hashData4_1 := sha256.Sum256(dataToSign4_1)
	sig4_1, err := rsa.SignPKCS1v15(rand.Reader, key1, crypto.SHA256, hashData4_1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx4.AddSignature(sig4_1, 1)
	tx4.Finalize()

	// ----------------------------------
	// tx5: Dependent transaction spending tx1's output.
	// Spends tx1's output of 18.0, outputs 17.0 to pubKey3.
	// This transaction is valid in isolation but is dependent (tx1's output is not in original snapshot).
	tx5 := NewTransaction()
	tx5.AddInput(tx1.GetHash(), 0)
	tx5.AddOutput(17.0, pubKey3)
	dataToSign5 := tx5.GetDataToSign(0)
	hashData5 := sha256.Sum256(dataToSign5)
	sig5, err := rsa.SignPKCS1v15(rand.Reader, key3, crypto.SHA256, hashData5[:])
	if err != nil {
		t.Fatal(err)
	}
	tx5.AddSignature(sig5, 0)
	tx5.Finalize()

	// ----------------------------------
	// tx6: Invalid transaction spending O1 due to output exceeding input.
	// Spends 30.0 from O1, but outputs 31.0 to pubKey2.
	tx6 := NewTransaction()
	tx6.AddInput(tx0.GetHash(), 1)
	tx6.AddOutput(31.0, pubKey2)
	dataToSign6 := tx6.GetDataToSign(0)
	hashData6 := sha256.Sum256(dataToSign6)
	sig6, err := rsa.SignPKCS1v15(rand.Reader, key2, crypto.SHA256, hashData6[:])
	if err != nil {
		t.Fatal(err)
	}
	tx6.AddSignature(sig6, 0)
	tx6.Finalize()

	// ----------------------------------
	// tx7: Another transaction attempting to spend O1.
	// Spends 30.0, outputs 28.0 to pubKey5.
	// This is a double spend of O1 (tx3 already spends it).
	tx7 := NewTransaction()
	tx7.AddInput(tx0.GetHash(), 1)
	tx7.AddOutput(28.0, pubKey5)
	dataToSign7 := tx7.GetDataToSign(0)
	hashData7 := sha256.Sum256(dataToSign7)
	sig7, err := rsa.SignPKCS1v15(rand.Reader, key2, crypto.SHA256, hashData7[:])
	if err != nil {
		t.Fatal(err)
	}
	tx7.AddSignature(sig7, 0)
	tx7.Finalize()

	// ----------------------------------
	// Build the list of possible transactions.
	// This list includes independent valid transactions (tx1, tx3),
	// a dependent one (tx5), and several invalid ones (tx2, tx4, tx6, tx7).
	possibleTxs := []Transaction{*tx1, *tx2, *tx3, *tx4, *tx5, *tx6, *tx7}

	// Process transactions.
	handler := Handler(possibleTxs)
	println(handler)

	assert.NotNil(t, handler, "Handler result should not be nil.")
	assert.Equal(t, 3, len(handler), "Three transactions should be processed immediately.")

	acceptedKeys := map[string]bool{
		tx1.Key(): true,
		tx3.Key(): true,
		tx5.Key(): true,
	}
	for _, tx := range handler {
		assert.True(t, acceptedKeys[tx.Key()], "Accepted transaction should be either tx1 or tx3 or tx5.")
	}
}

func TestHandleTxs_repeatedInvocationsReflectPoolUpdates(t *testing.T) {
	// Reset the global UTXO pool.
	utxoPool = NewUTXOPool()

	// Generate a key pair for funding and spending.
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	pubKey := &privKey.PublicKey

	// Create a funding transaction (tx0) with two outputs.
	// Output0: 20.0 funds; Output1: 30.0 funds.
	tx0 := NewTransaction()
	tx0.AddOutput(20.0, pubKey) // UTXO0
	tx0.AddOutput(30.0, pubKey) // UTXO1
	tx0.Finalize()

	// Insert both outputs into the global UTXO pool.
	utxo0 := NewUTXO(tx0.GetHash(), 0)
	utxo1 := NewUTXO(tx0.GetHash(), 1)
	utxoPool.AddUTXO(*utxo0, *tx0.GetOutput(0))
	utxoPool.AddUTXO(*utxo1, *tx0.GetOutput(1))

	// --- First Batch ---
	// tx1: Spend UTXO0 (20.0 funds). Produce one output (e.g., 18.0 funds; fee 2.0).
	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddOutput(18.0, pubKey)
	dataToSign1 := tx1.GetDataToSign(0)
	hashData1 := sha256.Sum256(dataToSign1)
	sig1, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(sig1, 0)
	tx1.Finalize()

	possibleBatch1 := []Transaction{*tx1}
	acceptedBatch1 := Handler(possibleBatch1)
	// Expect tx1 to be processed.
	assert.Equal(t, 1, len(acceptedBatch1), "Expected 1 accepted transaction in first batch")
	assert.Equal(t, tx1.Key(), acceptedBatch1[0].Key(), "Accepted transaction should be tx1")

	// At this point, the global UTXO pool should have:
	// - UTXO1 from tx0 (30.0 funds), and
	// - A new UTXO from tx1's output (18.0 funds).
	// Total expected UTXOs: 2.
	allUTXOs := utxoPool.GetAllUTXO()
	assert.Equal(t, 2, len(allUTXOs), "UTXO pool should contain 2 UTXOs after first batch")

	// --- Second Batch ---
	// tx2: Spend the output from tx1. (tx1's output should now be in the pool.)
	// For example, tx2 spends 18.0 funds and produces an output of 16.0 (fee 2.0).
	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 0)
	tx2.AddOutput(16.0, pubKey)
	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)
	sig2, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(sig2, 0)
	tx2.Finalize()

	possibleBatch2 := []Transaction{*tx2}
	acceptedBatch2 := Handler(possibleBatch2)
	// Expect tx2 to be processed.
	assert.Equal(t, 1, len(acceptedBatch2), "Expected 1 accepted transaction in second batch")
	assert.Equal(t, tx2.Key(), acceptedBatch2[0].Key(), "Accepted transaction should be tx2")

	// Now, the global UTXO pool should reflect:
	// - UTXO1 from tx0 (30.0 funds) still, and
	// - tx2's output (16.0 funds) replacing tx1's output.
	allUTXOs = utxoPool.GetAllUTXO()
	assert.Equal(t, 2, len(allUTXOs), "UTXO pool should contain 2 UTXOs after second batch")
}
