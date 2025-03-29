package first_faza

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMaxFeeHandleTxs_withSimpleValidTransaction(t *testing.T) {
	utxoPool = NewUTXOPool()

	// Generate a key pair for funding/spending.
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

	// Create a funding transaction (tx0) with two outputs.
	tx0 := NewTransaction()
	tx0.AddOutput(10.0, pubKey1) // UTXO0 with 10.0 funds.
	tx0.AddOutput(20.0, pubKey2) // UTXO1 with 20.0 funds.
	tx0.Finalize()

	// Add both outputs to the global UTXO pool.
	utxo0a := NewUTXO(tx0.GetHash(), 0)
	utxo0b := NewUTXO(tx0.GetHash(), 1)
	utxoPool.AddUTXO(*utxo0a, *tx0.GetOutput(0))
	utxoPool.AddUTXO(*utxo0b, *tx0.GetOutput(1))

	// Create tx1: Spending UTXO0.
	// tx1 spends 10.0 and outputs 7.0, yielding a fee of 3.0.
	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddOutput(7.0, pubKey2)
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

	// Create tx2: Spending UTXO1.
	// tx2 spends 20.0 and outputs 16.0, yielding a fee of 4.0.
	tx2 := NewTransaction()
	tx2.AddInput(tx0.GetHash(), 1)
	tx2.AddOutput(16.0, pubKey2)
	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)
	sig2, err := rsa.SignPKCS1v15(rand.Reader, privKey2, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(sig2, 0)
	tx2.Finalize()

	utxo2a := NewUTXO(tx2.GetHash(), 0)
	utxoPool.AddUTXO(*utxo2a, *tx2.GetOutput(0))

	tx3 := NewTransaction()
	tx3.AddInput(tx2.GetHash(), 0)
	tx3.AddOutput(10.0, pubKey1)
	tx3.AddOutput(1.0, pubKey2)
	dataToSign3 := tx3.GetDataToSign(0)
	hashData3 := sha256.Sum256(dataToSign3)
	sig3, err := rsa.SignPKCS1v15(rand.Reader, privKey2, crypto.SHA256, hashData3[:])
	if err != nil {
		t.Fatal(err)
	}
	tx3.AddSignature(sig3, 0)
	tx3.Finalize()

	// Process the transactions using MaxFeeHandler.
	possibleTxs := []*Transaction{tx1, tx2, tx3}
	accepted := MaxFeeHandler(possibleTxs)

	// Expect both transactions to be accepted.
	assert.NotNil(t, accepted, "Accepted transactions should not be nil.")
	assert.Equal(t, 3, len(accepted), "Expected 2 transactions to be accepted.")

	// The fee for tx1 is 3.0 and for tx2 is 4.0.
	// Since MaxFeeHandler sorts accepted transactions by fee descending,
	// tx2 should appear before tx1.
	assert.Equal(t, tx3.Key(), accepted[0].Key(), "Transaction with higher fee (tx3) should be first.")
	assert.Equal(t, tx2.Key(), accepted[1].Key(), "Transaction with average fee (tx2) should be second.")
	assert.Equal(t, tx1.Key(), accepted[2].Key(), "Transaction with lower fee (tx1) should be third.")
}

func TestMaxFeeHandleTxs_withTwoTransactionsClaimingSameOutput(t *testing.T) {
	utxoPool = NewUTXOPool()

	// Generate keys.
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

	// Create a funding transaction (tx0) with one output.
	tx0 := NewTransaction()
	tx0.AddOutput(10.0, pubKey1) // Output of 10.0 funds.
	tx0.AddOutput(10.0, pubKey2) // Output of 10.0 funds.
	tx0.Finalize()

	utxo0a := NewUTXO(tx0.GetHash(), 0)
	utxo0b := NewUTXO(tx0.GetHash(), 1)
	utxoPool.AddUTXO(*utxo0a, *tx0.GetOutput(0))
	utxoPool.AddUTXO(*utxo0b, *tx0.GetOutput(1))

	// tx1: spends tx0's output, outputs 7.0 funds (fee = 10.0 - 7.0 = 3.0).
	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddOutput(7.0, pubKey1)
	dataToSign1 := tx1.GetDataToSign(0)
	hashData1 := sha256.Sum256(dataToSign1)
	sig1, err := rsa.SignPKCS1v15(rand.Reader, privKey1, crypto.SHA256, hashData1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(sig1, 0)
	tx1.Finalize()

	// tx2: spends the same UTXO, outputs 6.5 funds (fee = 10.0 - 6.5 = 3.5).
	tx2 := NewTransaction()
	tx2.AddInput(tx0.GetHash(), 0)
	tx2.AddOutput(6.5, pubKey1)
	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)
	sig2, err := rsa.SignPKCS1v15(rand.Reader, privKey1, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(sig2, 0)
	tx2.Finalize()

	tx3 := NewTransaction()
	tx3.AddInput(tx0.GetHash(), 1)
	tx3.AddOutput(2.0, pubKey1)
	dataToSign3 := tx3.GetDataToSign(0)
	hashData3 := sha256.Sum256(dataToSign3)
	sig3, err := rsa.SignPKCS1v15(rand.Reader, privKey2, crypto.SHA256, hashData3[:])
	if err != nil {
		t.Fatal(err)
	}
	tx3.AddSignature(sig3, 0)
	tx3.Finalize()

	utxo3a := NewUTXO(tx3.GetHash(), 0)
	utxoPool.AddUTXO(*utxo3a, *tx3.GetOutput(0))

	// Both tx1 and tx2 claim the same UTXO.
	// We want the one with the higher fee (tx2) to be accepted.
	// To simulate this, we pass the transactions in an order where tx2 comes first.
	possibleTxs := []*Transaction{tx2, tx1, tx3}
	accepted := MaxFeeHandler(possibleTxs)

	// Expect only one transaction to be accepted.
	assert.NotNil(t, accepted, "Accepted transactions should not be nil.")
	assert.Equal(t, 2, len(accepted), "Only two transaction should be accepted because two of 3 double spend the same UTXO.")
	// The accepted transaction should be tx2, which has the higher fee.
	assert.Equal(t, tx3.Key(), accepted[0].Key(), "The transaction with the higher fee (tx3) should be accepted.")
	assert.Equal(t, tx2.Key(), accepted[1].Key(), "The transaction with the lower fee (tx2) should be accepted.")
}

func TestMaxFeeHandleTxs_withComplexTransaction_someValid_someInvalid(t *testing.T) {
	// Reset the global UTXO pool.
	utxoPool = NewUTXOPool()

	// Generate three key pairs.
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

	// Funding transaction (tx0): Create three outputs.
	// UTXO0: 10.0 funds to pubKey1.
	// UTXO1: 20.0 funds to pubKey1.
	// UTXO2: 15.0 funds to pubKey2.
	tx0 := NewTransaction()
	tx0.AddOutput(10.0, pubKey1)
	tx0.AddOutput(20.0, pubKey1)
	tx0.AddOutput(15.0, pubKey2)
	tx0.Finalize()

	utxo0 := NewUTXO(tx0.GetHash(), 0)
	utxo1 := NewUTXO(tx0.GetHash(), 1)
	utxo2 := NewUTXO(tx0.GetHash(), 2)
	utxoPool.AddUTXO(*utxo0, *tx0.GetOutput(0))
	utxoPool.AddUTXO(*utxo1, *tx0.GetOutput(1))
	utxoPool.AddUTXO(*utxo2, *tx0.GetOutput(2))

	// ----------------------------------
	// tx1: Valid independent transaction spending UTXO0.
	// Spends 10.0 and outputs 8.0 (fee = 2.0) to pubKey2.
	tx1 := NewTransaction()
	tx1.AddInput(tx0.GetHash(), 0)
	tx1.AddOutput(7.0, pubKey2)
	dataToSign1 := tx1.GetDataToSign(0)
	hashData1 := sha256.Sum256(dataToSign1)
	sig1, err := rsa.SignPKCS1v15(rand.Reader, key1, crypto.SHA256, hashData1[:])
	if err != nil {
		t.Fatal(err)
	}
	tx1.AddSignature(sig1, 0)
	tx1.Finalize()

	uxto1a := NewUTXO(tx1.GetHash(), 0)
	utxoPool.AddUTXO(*uxto1a, *tx1.GetOutput(0))

	// tx2: Invalid transaction: references a non-existing UTXO (fake hash).
	tx2 := NewTransaction()
	fakeHash := make([]byte, 32)
	for i := range fakeHash {
		fakeHash[i] = 42
	}
	tx2.AddInput(fakeHash, 0)
	tx2.AddOutput(5.0, pubKey2)
	dataToSign2 := tx2.GetDataToSign(0)
	hashData2 := sha256.Sum256(dataToSign2)
	sig2, err := rsa.SignPKCS1v15(rand.Reader, key1, crypto.SHA256, hashData2[:])
	if err != nil {
		t.Fatal(err)
	}
	tx2.AddSignature(sig2, 0)
	tx2.Finalize()

	// tx3: Invalid double spend: spends UTXO0 again.
	// Spends 10.0 and outputs 7.0.
	tx3 := NewTransaction()
	tx3.AddInput(tx0.GetHash(), 0)
	tx3.AddOutput(7.0, pubKey3)
	dataToSign3 := tx3.GetDataToSign(0)
	hashData3 := sha256.Sum256(dataToSign3)
	sig3, err := rsa.SignPKCS1v15(rand.Reader, key1, crypto.SHA256, hashData3[:])
	if err != nil {
		t.Fatal(err)
	}
	tx3.AddSignature(sig3, 0)
	tx3.Finalize()

	// tx4: Dependent valid transaction: spends output from tx1.
	// In isolation, tx4 is valid: it spends tx1's output (8.0) and outputs 7.0.
	tx4 := NewTransaction()
	tx4.AddInput(tx1.GetHash(), 0)
	tx4.AddOutput(6.0, pubKey2)
	dataToSign4 := tx4.GetDataToSign(0)
	hashData4 := sha256.Sum256(dataToSign4)
	sig4, err := rsa.SignPKCS1v15(rand.Reader, key2, crypto.SHA256, hashData4[:])
	if err != nil {
		t.Fatal(err)
	}
	tx4.AddSignature(sig4, 0)
	tx4.Finalize()

	utxo4a := NewUTXO(tx4.GetHash(), 0)
	utxoPool.AddUTXO(*utxo4a, *tx4.GetOutput(0))

	// tx5: Valid independent transaction spending UTXO2.
	// Change output so fee is higher: spends 15.0 and outputs 11.0 (fee = 4.0) to pubKey2.
	tx5 := NewTransaction()
	tx5.AddInput(tx0.GetHash(), 2)
	tx5.AddOutput(11.0, pubKey3)
	dataToSign5 := tx5.GetDataToSign(0)
	hashData5 := sha256.Sum256(dataToSign5)
	sig5, err := rsa.SignPKCS1v15(rand.Reader, key2, crypto.SHA256, hashData5[:])
	if err != nil {
		t.Fatal(err)
	}
	tx5.AddSignature(sig5, 0)
	tx5.Finalize()

	utxo5a := NewUTXO(tx5.GetHash(), 0)
	utxoPool.AddUTXO(*utxo5a, *tx5.GetOutput(0))

	// tx6: Invalid transaction: spending UTXO1 with output exceeding input.
	// Spends 20.0 but outputs 21.0.
	tx6 := NewTransaction()
	tx6.AddInput(tx0.GetHash(), 1)
	tx6.AddOutput(21.0, pubKey2)
	dataToSign6 := tx6.GetDataToSign(0)
	hashData6 := sha256.Sum256(dataToSign6)
	sig6, err := rsa.SignPKCS1v15(rand.Reader, key1, crypto.SHA256, hashData6[:])
	if err != nil {
		t.Fatal(err)
	}
	tx6.AddSignature(sig6, 0)
	tx6.Finalize()

	tx7 := NewTransaction()        // 2
	tx7.AddInput(tx5.GetHash(), 0) // pubkey3    --- 11
	tx7.AddOutput(7.0, pubKey3)
	tx7.AddOutput(2.0, pubKey2)
	dataToSign7 := tx7.GetDataToSign(0)
	hashData7 := sha256.Sum256(dataToSign7)
	sig7, err := rsa.SignPKCS1v15(rand.Reader, key3, crypto.SHA256, hashData7[:])
	if err != nil {
		t.Fatal(err)
	}
	tx7.AddSignature(sig7, 0)
	tx7.Finalize()

	utxo7a := NewUTXO(tx7.GetHash(), 0)
	utxo7b := NewUTXO(tx7.GetHash(), 0)
	utxoPool.AddUTXO(*utxo7a, *tx7.GetOutput(0))
	utxoPool.AddUTXO(*utxo7b, *tx7.GetOutput(0))

	// Construct the list of possible transactions.
	// We have: tx1 and tx5 as valid independent transactions,
	// tx2, tx3, and tx6 as invalid, and tx4 as dependent.
	possibleTxs := []*Transaction{tx1, tx2, tx3, tx4, tx5, tx6, tx7}
	handler := MaxFeeHandler(possibleTxs)

	// Expected outcome:
	// Only the independent valid transactions spending outputs from the original funding (tx0) are accepted.
	// These should be tx1 (spending UTXO0) and tx5 (spending UTXO2).
	// tx2, tx3, and tx6 are invalid.
	assert.NotNil(t, handler, "Handler result should not be nil.")
	assert.Equal(t, 4, len(handler), "Exactly 3 transactions should be processed immediately.")

	acceptedKeys := map[string]bool{
		tx1.Key(): true,
		tx4.Key(): true,
		tx5.Key(): true,
		tx7.Key(): true,
	}
	for _, tx := range handler {
		assert.True(t, acceptedKeys[tx.Key()], "Accepted transaction should be either tx1 or tx5 or tx4 or tx7.")
	}

	assert.Equal(t, tx5.Key(), handler[0].Key(), "Transaction with higher fee (tx5) should be first.")
	assert.Equal(t, tx1.Key(), handler[1].Key(), "Transaction with lower fee (tx1) should be second.")
	assert.Equal(t, tx7.Key(), handler[2].Key(), "Transaction with lower fee (tx7) should be second.")
	assert.Equal(t, tx4.Key(), handler[3].Key(), "Transaction with lower fee (tx4) should be second.")

}
