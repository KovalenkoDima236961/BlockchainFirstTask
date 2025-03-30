package third_faza

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockchain_processBlock_Without_Transaction(t *testing.T) {

	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyAlice := &privateKeyAlice.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	block1 := NewBlock(genesisBlock.GetHash(), pubKeyAlice)
	block1.Finalizee()

	result := BlockProcess(block1)

	assert.True(t, result, "Block with no transactions should be accepted")
}

func TestBlockchain_processBlock_With_One_Transaction(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyAlice := &privateKeyAlice.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()
	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	block1 := NewBlock(genesisBlock.GetHash(), pubKeyAlice)
	tx1 := NewTransaction()
	tx1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)

	tx1.AddOutput(1, pubKeyBob)
	tx1.AddOutput(1, pubKeyAlice)
	tx1.AddOutput(1.120, pubKeyBob)

	tx1.SignTx(privateKeyBob, 0)

	block1.TransactionAdd(tx1)
	block1.Finalizee()

	result := BlockProcess(block1)

	assert.True(t, result, "Block with one transactions should be accepted")
}

func TestBlockchain_processBlock_With_A_Lot_Transaction(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyAlice := &privateKeyAlice.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()
	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	block1 := NewBlock(genesisBlock.GetHash(), pubKeyAlice)

	tx1 := NewTransaction()
	tx1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx1.AddOutput(1, pubKeyAlice)
	tx1.AddOutput(2.125, pubKeyBob)
	tx1.SignTx(privateKeyBob, 0)

	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 1)
	tx2.AddOutput(1, pubKeyBob)
	tx2.AddOutput(1.125, pubKeyAlice)
	tx2.SignTx(privateKeyBob, 0)

	block1.TransactionAdd(tx1)
	block1.TransactionAdd(tx2)
	block1.Finalizee()

	result := BlockProcess(block1)
	assert.True(t, result, "First Block with few transactions should be accepted")

	block2 := NewBlock(block1.GetHash(), pubKeyBob)
	tx1 = NewTransaction()
	tx1.AddInput(block1.GetCoinbase().GetHash(), 0)
	tx1.AddOutput(1, pubKeyBob)
	tx1.AddOutput(1.125, pubKeyBob)
	tx1.SignTx(privateKeyAlice, 0)

	block2.TransactionAdd(tx1)
	block2.Finalizee()

	result = BlockProcess(block2)
	assert.True(t, result, "Second Block with few transactions should be accepted")
}

func TestBlockchain_processBlock_With_Some_DoubleSpend(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyAlice := &privateKeyAlice.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()
	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	block1 := NewBlock(genesisBlock.GetHash(), pubKeyAlice)

	tx1 := NewTransaction()
	tx1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx1.AddOutput(1, pubKeyAlice)
	tx1.AddOutput(2.125, pubKeyBob)
	tx1.SignTx(privateKeyBob, 0)

	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 1)
	tx2.AddOutput(1, pubKeyBob)
	tx2.AddOutput(1.125, pubKeyAlice)
	tx2.SignTx(privateKeyBob, 0)

	block1.TransactionAdd(tx1)
	block1.TransactionAdd(tx2)
	block1.Finalizee()

	result := BlockProcess(block1)
	assert.True(t, result, "First Block with few valid transactions should be accepted")

	block2 := NewBlock(block1.GetHash(), pubKeyBob)
	tx1 = NewTransaction()
	tx1.AddInput(block1.GetCoinbase().GetHash(), 0)
	tx1.AddOutput(1, pubKeyBob)
	tx1.AddOutput(1.125, pubKeyBob)
	tx1.SignTx(privateKeyAlice, 0)
	block2.TransactionAdd(tx1)

	tx2 = NewTransaction()
	tx2.AddInput(block1.GetCoinbase().GetHash(), 0)
	tx2.AddOutput(1, pubKeyBob)
	tx2.SignTx(privateKeyAlice, 0)
	block2.TransactionAdd(tx2)

	block2.Finalizee()

	result = BlockProcess(block2)
	assert.False(t, result, "Second Block with few invalid transactions shouldn't be accepted")
}

func TestBlockchain_processBlock_With_New_GenesisBlock(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyAlice := &privateKeyAlice.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()
	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	block1 := NewBlock(nil, pubKeyAlice)
	block1.Finalizee()

	result := BlockProcess(block1)

	assert.False(t, result, "Second genesis block shouldn't be accepted")
}

func TestBlockchain_processBlock_With_Incorrect_PrevBlockHash(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyAlice := &privateKeyAlice.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()
	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	block1 := NewBlock(genesisBlock.GetHash(), pubKeyAlice)

	tx1 := NewTransaction()
	tx1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx1.AddOutput(1, pubKeyAlice)
	tx1.AddOutput(2.125, pubKeyBob)
	tx1.SignTx(privateKeyBob, 0)

	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 1)
	tx2.AddOutput(1, pubKeyBob)
	tx2.AddOutput(1.125, pubKeyAlice)
	tx2.SignTx(privateKeyBob, 0)

	block1.TransactionAdd(tx1)
	block1.TransactionAdd(tx2)
	block1.Finalizee()

	result := BlockProcess(block1)
	assert.True(t, result, "First Block with few transactions should be accepted")

	invalidHash := make([]byte, len(block1.GetHash()))
	copy(invalidHash, block1.GetHash())

	for i := range invalidHash {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		invalidHash[i] = byte(n.Int64())
	}

	println("Invalid hash: ", invalidHash)
	println("Block1 hash: ", block1.GetHash())

	block2 := NewBlock(invalidHash, pubKeyBob)
	tx1 = NewTransaction()
	tx1.AddInput(block1.GetCoinbase().GetHash(), 0)
	tx1.AddOutput(1, pubKeyBob)
	tx1.AddOutput(1.125, pubKeyBob)
	tx1.SignTx(privateKeyAlice, 0)

	block2.TransactionAdd(tx1)
	block2.Finalizee()

	result = BlockProcess(block2)
	assert.False(t, result, "Second Block with few transactions should be accepted")
}

func TestBlockchain_processBlock_With_DifferentTypeInvalidTransactions(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyAlice := &privateKeyAlice.PublicKey

	// Genesis block
	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()
	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	// Valid block (block1)
	block1 := NewBlock(genesisBlock.GetHash(), pubKeyAlice)

	tx1 := NewTransaction()
	tx1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx1.AddOutput(1, pubKeyAlice)
	tx1.AddOutput(2.125, pubKeyBob)
	tx1.SignTx(privateKeyBob, 0)

	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 1)
	tx2.AddOutput(1, pubKeyBob)
	tx2.AddOutput(1.125, pubKeyAlice)
	tx2.SignTx(privateKeyBob, 0)

	block1.TransactionAdd(tx1)
	block1.TransactionAdd(tx2)
	block1.Finalizee()

	result := BlockProcess(block1)
	assert.True(t, result, "First Block with few valid transactions should be accepted")

	// Create invalid block (block2) with multiple errors
	block2 := NewBlock(block1.GetHash(), pubKeyBob)

	// 1. Double spending the same input as tx1
	doubleSpendTx := NewTransaction()
	doubleSpendTx.AddInput(block1.GetCoinbase().GetHash(), 0)
	doubleSpendTx.AddOutput(1, pubKeyBob)
	doubleSpendTx.SignTx(privateKeyAlice, 0)
	block2.TransactionAdd(doubleSpendTx)

	doubleSpendTx2 := NewTransaction()
	doubleSpendTx2.AddInput(block1.GetCoinbase().GetHash(), 0) // same input reused
	doubleSpendTx2.AddOutput(1, pubKeyBob)
	doubleSpendTx2.SignTx(privateKeyAlice, 0)
	block2.TransactionAdd(doubleSpendTx2)

	// 2. Invalid signature
	invalidSigTx := NewTransaction()
	invalidSigTx.AddInput(tx1.GetHash(), 0)
	invalidSigTx.AddOutput(1, pubKeyBob)
	invalidSigTx.SignTx(privateKeyAlice, 0) // should be signed by Bob
	block2.TransactionAdd(invalidSigTx)

	// 3. Input < Output (overdraft)
	overdraftTx := NewTransaction()
	overdraftTx.AddInput(tx2.GetHash(), 0)
	overdraftTx.AddOutput(5, pubKeyAlice) // more than available
	overdraftTx.SignTx(privateKeyBob, 0)
	block2.TransactionAdd(overdraftTx)

	block2.Finalizee()

	result = BlockProcess(block2)
	assert.False(t, result, "Block with invalid transactions (double spend, bad signature, overdraft) should be rejected")

	// Only Input < Output (with normal tx, no coinbase)
	block2 = NewBlock(block1.GetHash(), pubKeyBob)
	wrongTx := NewTransaction()
	wrongTx.AddInput(tx2.GetHash(), 0)
	wrongTx.AddOutput(10, pubKeyBob)
	wrongTx.SignTx(privateKeyBob, 0)
	block2.TransactionAdd(wrongTx)
	block2.Finalizee()

	result = BlockProcess(block2)
	assert.False(t, result, "Block should be rejected due to input < output in transaction without using coinbase")
}

func TestBlockchain_processBlock_WithFewBlocksAboveGenesisBlock(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyAlice := &privateKeyAlice.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()
	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	block1 := NewBlock(genesisBlock.GetHash(), pubKeyAlice)

	tx1 := NewTransaction()
	tx1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx1.AddOutput(1, pubKeyAlice)
	tx1.SignTx(privateKeyBob, 0)

	block1.TransactionAdd(tx1)
	block1.Finalizee()

	result := BlockProcess(block1)
	assert.True(t, result, "First Block with few valid transactions should be accepted")

	block2 := NewBlock(block1.GetHash(), pubKeyBob)
	tx1 = NewTransaction()
	tx1.AddInput(block1.GetCoinbase().GetHash(), 0)
	tx1.AddOutput(1, pubKeyBob)
	tx1.AddOutput(1.125, pubKeyBob)
	tx1.SignTx(privateKeyAlice, 0)
	block2.TransactionAdd(tx1)

	tx2 := NewTransaction()
	tx2.AddInput(block1.GetCoinbase().GetHash(), 0)
	tx2.AddOutput(1, pubKeyBob)
	tx2.SignTx(privateKeyAlice, 0)
	block2.TransactionAdd(tx2)

	block2.Finalizee()

	result = BlockProcess(block2)
	assert.False(t, result, "Second Block with few invalid transactions shouldn't be accepted")

	block3 := NewBlock(genesisBlock.GetHash(), pubKeyAlice)

	tx3_1 := NewTransaction()
	tx3_1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx3_1.AddOutput(1, pubKeyAlice)
	tx3_1.SignTx(privateKeyBob, 0)

	block3.TransactionAdd(tx3_1)
	block3.Finalizee()

	result = BlockProcess(block3)
	assert.True(t, result, "Second Block above genesis block is accepted")

	block4 := NewBlock(genesisBlock.GetHash(), pubKeyAlice)

	tx4_1 := NewTransaction()
	tx4_1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx4_1.AddOutput(2, pubKeyAlice)
	tx4_1.SignTx(privateKeyBob, 0)

	block4.TransactionAdd(tx4_1)
	block4.Finalizee()

	result = BlockProcess(block4)
	assert.True(t, result, "Third Block above genesis block is accepted")

	number_of_nodes := len(blockchain.MaxHeightNode)
	assert.Equal(t, 3, number_of_nodes, "Should be 3 blocks above genesis block")
}

func TestBlockchain_processBlock_ThatClaimsUTXO_whichHasAlreadyBeenClaimedByTransactionInTheParentBlock(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyAlice := &privateKeyAlice.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()
	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	block1 := NewBlock(genesisBlock.GetHash(), pubKeyAlice)

	tx1 := NewTransaction()
	tx1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx1.AddOutput(1, pubKeyAlice)
	tx1.SignTx(privateKeyBob, 0)

	block1.TransactionAdd(tx1)
	block1.Finalizee()

	result := BlockProcess(block1)
	assert.True(t, result, "First Block with few valid transactions should be accepted")

	block2 := NewBlock(block1.GetHash(), pubKeyAlice)
	tx1 = NewTransaction()
	tx1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx1.AddOutput(2, pubKeyAlice)
	tx1.SignTx(privateKeyBob, 0)

	block2.TransactionAdd(tx1)
	block2.Finalizee()

	result = BlockProcess(block2)
	assert.False(t, result, "Second Block have transaction that used utxo which was already claimed by transaction in parent block")
}

func TestBlockchain_processBlock_WhichContaining_Transaction_ThatClaims_UTXO_From_Outside_Its_Branch(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyAlice := &privateKeyAlice.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()
	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	blockA := NewBlock(genesisBlock.GetHash(), pubKeyAlice)

	txA_1 := NewTransaction()
	txA_1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	txA_1.AddOutput(1, pubKeyAlice)
	txA_1.SignTx(privateKeyBob, 0)

	blockA.TransactionAdd(txA_1)
	blockA.Finalizee()

	result := BlockProcess(blockA)
	assert.True(t, result, "First Block with few valid transactions should be accepted")

	blockB_1 := NewBlock(blockA.GetHash(), pubKeyBob)
	txB_1_1 := NewTransaction()
	txB_1_1.AddInput(blockA.GetCoinbase().GetHash(), 0)
	txB_1_1.AddOutput(1, pubKeyBob)
	txB_1_1.AddOutput(1.125, pubKeyBob)
	txB_1_1.SignTx(privateKeyAlice, 0)
	blockB_1.TransactionAdd(txB_1_1)

	blockB_1.Finalizee()

	result = BlockProcess(blockB_1)
	assert.True(t, result, "Second Block with few invalid transactions shouldn't be accepted")

	blockB_2 := NewBlock(blockA.GetHash(), pubKeyAlice)

	txB_2_1 := NewTransaction()
	txB_2_1.AddInput(blockA.GetCoinbase().GetHash(), 0)
	txB_2_1.AddOutput(1, pubKeyAlice)
	txB_2_1.AddOutput(1, pubKeyBob)
	txB_2_1.SignTx(privateKeyAlice, 0)

	blockB_2.TransactionAdd(txB_2_1)
	blockB_2.Finalizee()

	result = BlockProcess(blockB_2)
	assert.True(t, result, "Second Block above genesis block is accepted")

	number_of_nodes := len(blockchain.MaxHeightNode)
	assert.Equal(t, 2, number_of_nodes, "Should be 2 blocks above genesis block")

	blockC1 := NewBlock(blockB_1.GetHash(), pubKeyAlice)

	txC1_1 := NewTransaction()
	txC1_1.AddInput(txB_2_1.GetHash(), 1)
	txC1_1.AddOutput(1, pubKeyAlice)
	txC1_1.SignTx(privateKeyBob, 0)

	blockC1.TransactionAdd(txC1_1)
	blockC1.Finalizee()

	result = BlockProcess(blockC1)
	assert.False(t, result, "Block should not be accepted because contain transaction that claims utxo from outside its branch")

}

func TestBlockchain_processBlock_WhichContaining_Transaction_ThatClaimsAnOlderUTXO_WithinTheSameBranch_ThatHasNotYetBeenClaimed(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyAlice := &privateKeyAlice.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()
	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	blockA := NewBlock(genesisBlock.GetHash(), pubKeyAlice)

	txA_1 := NewTransaction()
	txA_1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	txA_1.AddOutput(1, pubKeyAlice)
	txA_1.AddOutput(2.125, pubKeyBob)
	txA_1.SignTx(privateKeyBob, 0)

	blockA.TransactionAdd(txA_1)
	blockA.Finalizee()

	result := BlockProcess(blockA)
	assert.True(t, result, "First Block with few valid transactions should be accepted")

	blockB := NewBlock(blockA.GetHash(), pubKeyBob)
	txB_1 := NewTransaction()
	txB_1.AddInput(blockA.GetCoinbase().GetHash(), 0)
	txB_1.AddOutput(1, pubKeyBob)
	txB_1.AddOutput(2.125, pubKeyAlice)
	txB_1.SignTx(privateKeyAlice, 0)

	txB_2 := NewTransaction()
	txB_2.AddInput(txA_1.GetHash(), 1)
	txB_2.AddOutput(1, pubKeyBob)
	txB_2.AddOutput(1.125, pubKeyAlice)
	txB_2.SignTx(privateKeyBob, 0)

	blockB.TransactionAdd(txB_2)
	blockB.TransactionAdd(txB_1)

	blockB.Finalizee()

	result = BlockProcess(blockB)
	assert.True(t, result, "Second Block with few valid transactions should be accepted")

	blockC := NewBlock(blockB.GetHash(), pubKeyBob)
	txC_1 := NewTransaction()
	txC_1.AddInput(txA_1.GetHash(), 0)
	txC_1.AddOutput(1, pubKeyBob)
	txC_1.SignTx(privateKeyAlice, 0)

	blockC.TransactionAdd(txC_1)
	blockC.Finalizee()

	result = BlockProcess(blockC)
	assert.True(t, result, "Third Block should be accepted because contain transaction that claims an older utxo within the same branch")
}

func TestBlockchain_processBlock_LinearChainOfBlocks(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyAlice := &privateKeyAlice.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()
	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	blockA := NewBlock(genesisBlock.GetHash(), pubKeyAlice)

	txA_1 := NewTransaction()
	txA_1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	txA_1.AddOutput(1, pubKeyAlice)
	txA_1.AddOutput(2.125, pubKeyBob)
	txA_1.SignTx(privateKeyBob, 0)

	blockA.TransactionAdd(txA_1)
	blockA.Finalizee()

	result := BlockProcess(blockA)
	assert.True(t, result, "First Block with few valid transactions should be accepted")

	blockB := NewBlock(blockA.GetHash(), pubKeyBob)
	txB_1 := NewTransaction()
	txB_1.AddInput(blockA.GetCoinbase().GetHash(), 0)
	txB_1.AddOutput(1, pubKeyBob)
	txB_1.AddOutput(2.125, pubKeyAlice)
	txB_1.SignTx(privateKeyAlice, 0)

	txB_2 := NewTransaction()
	txB_2.AddInput(txA_1.GetHash(), 1)
	txB_2.AddOutput(1, pubKeyBob)
	txB_2.AddOutput(1.125, pubKeyAlice)
	txB_2.SignTx(privateKeyBob, 0)

	blockB.TransactionAdd(txB_2)
	blockB.TransactionAdd(txB_1)

	blockB.Finalizee()

	result = BlockProcess(blockB)
	assert.True(t, result, "Second Block with few valid transactions should be accepted")

	blockC := NewBlock(blockB.GetHash(), pubKeyBob)
	txC_1 := NewTransaction()
	txC_1.AddInput(txA_1.GetHash(), 0)
	txC_1.AddOutput(1, pubKeyBob)
	txC_1.SignTx(privateKeyAlice, 0)

	blockC.TransactionAdd(txC_1)
	blockC.Finalizee()

	result = BlockProcess(blockC)
	assert.True(t, result, "Third Block should be accepted because contain transaction that claims an older utxo within the same branch")

	blockD := NewBlock(blockC.GetHash(), pubKeyBob)
	txD_1 := NewTransaction()
	txD_1.AddInput(blockC.GetCoinbase().GetHash(), 0)
	txD_1.AddOutput(1, pubKeyBob)
	txD_1.AddOutput(2.125, pubKeyAlice)
	txD_1.SignTx(privateKeyBob, 0)

	txD_2 := NewTransaction()
	txD_2.AddInput(txB_2.GetHash(), 1)
	txD_2.AddOutput(1, pubKeyBob)
	txD_2.AddOutput(0.125, pubKeyAlice)
	txD_2.SignTx(privateKeyAlice, 0)

	blockD.TransactionAdd(txD_1)
	blockD.TransactionAdd(txD_2)
	blockD.Finalizee()

	result = BlockProcess(blockD)
	assert.True(t, result, "Forth Block should be accepted because contain transaction that claims an older utxo within the same branch")

	number_of_nodes := len(blockchain.BlockChain)
	assert.Equal(t, 5, number_of_nodes, "Blockchain should contain 5 nodes")
}

func TestBlockchain_processBlock_LinearChainOfLength_CUTOFFAGE_PlusOne_ThenBlockOnGenesis_ShouldFail(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyAlice := &privateKeyAlice.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	prevHash := make([]byte, len(genesisBlock.GetHash()))
	copy(prevHash, genesisBlock.GetHash())

	// Build a chain of length CUT_OFF_AGE (excluding Genesis)
	block := NewBlock(prevHash, pubKeyAlice)
	for i := 0; i <= CUT_OFF_AGE; i++ {
		block.Finalizee()

		result := BlockProcess(block)
		assert.True(t, result, fmt.Sprintf("Block #%d should be accepted", i+1))

		prevHash = make([]byte, len(block.GetHash()))
		copy(prevHash, block.GetHash())
		block = NewBlock(prevHash, pubKeyAlice)
	}

	// Now try to add a block on top of Genesis (should be rejected)
	blockOnGenesis := NewBlock(genesisBlock.GetHash(), pubKeyAlice)
	blockOnGenesis.Finalizee()

	result := BlockProcess(blockOnGenesis)
	println(result)
	assert.False(t, result, "Block on Genesis should be rejected because it is outside the CUT_OFF_AGE range")
}

func TestBlockchain_processBlock_LinearChainOfBlocksOfLength_CUT_OFF_AGE_Plus_One_ThenBlockOnTopOfTheGenesisBlock(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyAlice := &privateKeyAlice.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	prevHash := make([]byte, len(genesisBlock.GetHash()))
	copy(prevHash, genesisBlock.GetHash())

	for i := 0; i <= CUT_OFF_AGE+1; i++ {
		block := NewBlock(prevHash, pubKeyAlice)
		block.Finalizee()

		result := BlockProcess(block)
		assert.True(t, result, fmt.Sprintf("Block #%d should be accepted", i+1))

		prevHash = make([]byte, len(block.GetHash()))
		copy(prevHash, block.GetHash())
	}
}

func TestBlockchain_createBlock_WhenNoTransactionHaveBeenProcessed(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	block := BlockCreate(pubKeyBob)
	assert.NotNil(t, block, "Block with only coinbase and no transaction should be accepted")
}

func TestBlockchain_createBlock_AfterProcessingOneValidTransaction(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}

	pubKeyAlice := &privateKeyAlice.PublicKey
	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	tx := NewTransaction()
	tx.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx.AddOutput(1, pubKeyBob)
	tx.AddOutput(2.125, pubKeyAlice)
	tx.SignTx(privateKeyBob, 0)

	TxProcess(tx)

	block := BlockCreate(pubKeyBob)
	assert.NotNil(t, block, "Block with only one valid transaction should be accepted")
}

func TestBlockchain_createBlock_AfterProcessingValidTransaction_AndThenCreateSecondBlock(t *testing.T) {
	privateKeyBob, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}

	pubKeyAlice := &privateKeyAlice.PublicKey
	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	tx := NewTransaction()
	tx.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx.AddOutput(1, pubKeyBob)
	tx.AddOutput(2.125, pubKeyAlice)
	tx.SignTx(privateKeyBob, 0)

	TxProcess(tx)

	block := BlockCreate(pubKeyBob)
	assert.NotNil(t, block, "Block with valid transaction should be accepted")

	tx2 := NewTransaction()
	tx2.AddInput(tx.GetHash(), 1)
	tx2.AddOutput(1, pubKeyAlice)
	tx2.AddOutput(1.125, pubKeyBob)
	tx2.SignTx(privateKeyAlice, 0)

	TxProcess(tx2)

	tx3 := NewTransaction()
	tx3.AddInput(tx2.GetHash(), 1)
	tx3.AddOutput(0.125, pubKeyAlice)
	tx3.SignTx(privateKeyBob, 0)

	TxProcess(tx3)

	block2 := BlockCreate(pubKeyAlice)
	assert.NotNil(t, block2, "Block with only coinbase and no transaction should be accepted")
}

func TestBlockchain_createBlock_AfterProcessingValidTransaction_WhichIsAlreadyInABlockInTheLongestValidBranch(t *testing.T) {
	privateKeyBob, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyAlice := &privateKeyAlice.PublicKey

	// 1. Create Genesis block
	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	// 2. Create a valid transaction and process it into the global pool
	tx := NewTransaction()
	tx.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx.AddOutput(1, pubKeyBob)
	tx.AddOutput(2.125, pubKeyAlice)
	tx.SignTx(privateKeyBob, 0)

	TxProcess(tx)

	// 3. Create Block A (should include tx)
	blockA := BlockCreate(pubKeyAlice)
	assert.NotNil(t, blockA, "Block A should be created with valid transaction")
	assert.Equal(t, 2, len(blockA.GetTransactions()), "Block A should contain coinbase and the valid transaction")

	// 4. Try to process the same transaction again
	TxProcess(tx)

	// 5. Create Block B (should not include tx again)
	blockB := BlockCreate(pubKeyBob)

	// 6. Block B should contain only coinbase, because tx is already confirmed in Block A
	assert.NotNil(t, blockB, "Block B should be created")
	assert.Equal(t, 1, len(blockB.GetTransactions()), "Block B should contain only the coinbase transaction")
}

func TestBlockchain_createBlock_AfterProcessingValidTransaction_WhichUsesAUTXOThatHasAlreadyBeenClaimedByATransactionInTheLongestValidChain(t *testing.T) {
	privateKeyBob, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyAlice := &privateKeyAlice.PublicKey

	// 1. Create Genesis block
	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	// 2. Create a valid transaction and process it into the global pool
	tx := NewTransaction()
	tx.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx.AddOutput(1, pubKeyBob)
	tx.AddOutput(2.125, pubKeyAlice)
	tx.SignTx(privateKeyBob, 0)

	TxProcess(tx)

	// 3. Create Block A (should include tx)
	blockA := BlockCreate(pubKeyAlice)
	assert.NotNil(t, blockA, "Block A should be created with valid transaction")
	assert.Equal(t, 2, len(blockA.GetTransactions()), "Block A should contain coinbase and the valid transaction")

	// 4. Try to process the same transaction again
	tx2 := NewTransaction()
	tx2.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx2.AddOutput(1, pubKeyBob)
	tx2.AddOutput(2.125, pubKeyAlice)
	tx2.SignTx(privateKeyBob, 0)

	// 5. Create Block B (should not include tx again)
	blockB := BlockCreate(pubKeyBob)

	// 6. Block B should contain only coinbase, because tx is already confirmed in Block A
	assert.NotNil(t, blockB, "Block B should be created")
	assert.Equal(t, 1, len(blockB.GetTransactions()), "Block B should contain only the coinbase transaction")
}

func TestBlockchain_createBlock_AfterProcessingAValidTransactionThatIsNotADoubleSpendInTheLongestValidBranch_AndHasntBeenUsedInAnyOtherBlockYet(t *testing.T) {
	privateKeyBob, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyAlice := &privateKeyAlice.PublicKey

	// 1. Create Genesis block
	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	// 2. Create a valid transaction and process it into the global pool
	tx := NewTransaction()
	tx.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx.AddOutput(1, pubKeyBob)
	tx.AddOutput(2.125, pubKeyAlice)
	tx.SignTx(privateKeyBob, 0)

	TxProcess(tx)

	// 3. Create Block A (should include tx)
	blockA := BlockCreate(pubKeyAlice)
	assert.NotNil(t, blockA, "Block A should be created with valid transaction")
	assert.Equal(t, 2, len(blockA.GetTransactions()), "Block A should contain coinbase and the valid transaction")

	// 4. Try to process the same transaction again
	tx2 := NewTransaction()
	tx2.AddInput(blockA.GetCoinbase().GetHash(), 0)
	tx2.AddOutput(1, pubKeyBob)
	tx2.AddOutput(2.125, pubKeyAlice)
	tx2.SignTx(privateKeyAlice, 0)

	TxProcess(tx2)

	// 5. Create Block B (should not include tx again)
	blockB := BlockCreate(pubKeyBob)

	assert.NotNil(t, blockB, "Block B should be created")
	assert.Equal(t, 2, len(blockB.GetTransactions()), "Block B should contain coinbase and the valid transaction")

	tx3 := NewTransaction()
	tx3.AddInput(tx2.GetHash(), 1)
	tx3.AddOutput(1, pubKeyBob)
	tx3.AddOutput(1.125, pubKeyAlice)
	tx3.SignTx(privateKeyAlice, 0)

	TxProcess(tx3)

	tx4 := NewTransaction()
	tx4.AddInput(tx.GetHash(), 1)
	tx4.AddOutput(1, pubKeyBob)
	tx4.AddOutput(1.125, pubKeyAlice)
	tx4.SignTx(privateKeyAlice, 0)

	TxProcess(tx4)

	blockC := BlockCreate(pubKeyAlice)
	assert.NotNil(t, blockC, "Block C should be created")
	assert.Equal(t, 3, len(blockC.GetTransactions()), "Block C should contain only the coinbase transaction")
}

func TestBlockchain_createBlock_AfterProcessingOnlyInvalidTransactions(t *testing.T) {
	privateKeyBob, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyAlice := &privateKeyAlice.PublicKey

	// 1. Create Genesis block
	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	// 2. Create a valid transaction and process it into the global pool
	tx := NewTransaction()
	tx.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx.AddOutput(1, pubKeyBob)
	tx.AddOutput(2.125, pubKeyAlice)
	tx.SignTx(privateKeyBob, 0)

	TxProcess(tx)

	// 3. Create Block A (should include tx)
	blockA := BlockCreate(pubKeyAlice)
	assert.NotNil(t, blockA, "Block A should be created with valid transaction")
	assert.Equal(t, 2, len(blockA.GetTransactions()), "Block A should contain coinbase and the valid transaction")

	// 4. Try to process the same transaction again
	TxProcess(tx)

	tx3 := NewTransaction()

	fakeHash := make([]byte, 32)
	for i := range fakeHash {
		fakeHash[i] = 99
	}

	tx3.AddInput(fakeHash, 0)
	tx3.AddOutput(1, pubKeyBob)
	tx3.AddOutput(2.125, pubKeyAlice)
	tx3.SignTx(privateKeyAlice, 0)

	TxProcess(tx3)

	tx4 := NewTransaction()
	tx4.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx4.AddOutput(1, pubKeyBob)
	tx4.AddOutput(2.125, pubKeyAlice)
	tx4.SignTx(privateKeyBob, 0)

	TxProcess(tx4)

	// 5. Create Block B (should not include tx again)
	blockB := BlockCreate(pubKeyBob)

	assert.NotNil(t, blockB, "Block B should be created")
	assert.Equal(t, 1, len(blockB.GetTransactions()), "Block B should contain only coinbase")
}

func TestBlockchain_processTransaction_CreateBlock_ProcessAnotherTransaction_CreateBlock(t *testing.T) {
	privateKeyBob, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyAlice := &privateKeyAlice.PublicKey

	privateKeyCyril, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyCyril := &privateKeyCyril.PublicKey

	// 1. Create Genesis block
	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	tx := NewTransaction()
	tx.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx.AddOutput(2, pubKeyBob)
	tx.AddOutput(1.125, pubKeyAlice)
	tx.SignTx(privateKeyBob, 0)

	TxProcess(tx)

	blockA := BlockCreate(pubKeyAlice)
	assert.NotNil(t, blockA, "Block A should be created with valid transaction")
	assert.Equal(t, 2, len(blockA.GetTransactions()), "Block A should contain only coinbase")

	tx2 := NewTransaction()
	tx2.AddInput(tx.GetHash(), 0)
	tx2.AddOutput(1, pubKeyBob)
	tx2.AddOutput(1, pubKeyCyril)
	tx2.SignTx(privateKeyBob, 0)

	TxProcess(tx2)

	blockB := BlockCreate(pubKeyCyril)
	assert.NotNil(t, blockB, "Block B should be created with valid transaction")
	assert.Equal(t, 2, len(blockB.GetTransactions()), "Block B should contain only coinbase")

	tx3 := NewTransaction()
	tx3.AddInput(tx2.GetHash(), 1)
	tx3.AddOutput(1, pubKeyBob)
	tx3.SignTx(privateKeyCyril, 0)

	TxProcess(tx3)

	blockC := BlockCreate(pubKeyBob)
	assert.NotNil(t, blockC, "Block C should be created with valid transaction")
	assert.Equal(t, 2, len(blockC.GetTransactions()), "Block C should contain only coinbase")
}

func TestBlockchain_processTransaction_CreateBlock_ProcessBlockOnTopOfThatBlockWithATransactionClaimingTheUTXOFromThatPreviousTransaction(t *testing.T) {
	privateKeyBob, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyAlice := &privateKeyAlice.PublicKey

	privateKeyCyril, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyCyril := &privateKeyCyril.PublicKey

	// 1. Create Genesis block
	genesisBlock := NewBlock(nil, pubKeyAlice)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	tx := NewTransaction()
	tx.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx.AddOutput(2, pubKeyBob)
	tx.AddOutput(1.125, pubKeyAlice)
	tx.SignTx(privateKeyAlice, 0)

	TxProcess(tx)

	blockA := BlockCreate(pubKeyBob)
	assert.NotNil(t, blockA, "Block A should be created with valid transaction")
	assert.Equal(t, 2, len(blockA.GetTransactions()), "Block A should contain only coinbase")

	blockB := NewBlock(blockA.GetHash(), pubKeyCyril)
	txB_1 := NewTransaction()
	txB_1.AddInput(blockA.GetCoinbase().GetHash(), 0)
	txB_1.AddOutput(1, pubKeyBob)
	txB_1.AddOutput(2.125, pubKeyAlice)
	txB_1.SignTx(privateKeyBob, 0)

	txB_2 := NewTransaction()
	txB_2.AddInput(tx.GetHash(), 0)
	txB_2.AddOutput(1, pubKeyBob)
	txB_2.AddOutput(1, pubKeyAlice)
	txB_2.SignTx(privateKeyBob, 0)

	blockB.TransactionAdd(txB_2)
	blockB.TransactionAdd(txB_1)

	blockB.Finalizee()

	result := BlockProcess(blockB)
	assert.True(t, result, "Second Block with few valid transactions should be accepted")
}

func TestBlockchain_processTransaction_CreateBlock_ProcessBlockOnTopOfTheGenesisBlockWithATransactionClaimingTheUTXOFromThatPreviousTransaction(t *testing.T) {
	privateKeyBob, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyAlice := &privateKeyAlice.PublicKey

	privateKeyCyril, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyCyril := &privateKeyCyril.PublicKey

	// 1. Create Genesis block
	genesisBlock := NewBlock(nil, pubKeyAlice)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	tx := NewTransaction()
	tx.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	tx.AddOutput(1, pubKeyBob)
	tx.AddOutput(1.125, pubKeyAlice)
	tx.AddOutput(1, pubKeyCyril)
	tx.SignTx(privateKeyAlice, 0)

	TxProcess(tx)

	blockA := BlockCreate(pubKeyBob)
	assert.NotNil(t, blockA, "Block A should be created with valid transaction")
	assert.Equal(t, 2, len(blockA.GetTransactions()), "Block A should contain coinbase transaction and tx1")

	blockB := NewBlock(genesisBlock.GetHash(), pubKeyCyril)
	txB_1 := NewTransaction()
	txB_1.AddInput(tx.GetHash(), 0)
	txB_1.AddOutput(1, pubKeyBob)
	txB_1.SignTx(privateKeyBob, 0)

	txB_2 := NewTransaction()
	txB_1.AddInput(tx.GetHash(), 1)
	txB_1.AddOutput(1, pubKeyBob)
	txB_1.AddOutput(0.125, pubKeyAlice)
	txB_1.SignTx(privateKeyAlice, 0)

	blockB.TransactionAdd(txB_2)
	blockB.TransactionAdd(txB_1)

	blockB.Finalizee()

	result := BlockProcess(blockB)
	assert.False(t, result, "Blockchain must reject Block B, because its trying to spend a UTXO that wa never created in its branch")
}

func TestBlockchain_processMultipleBlocksDirectlyOnTopOfTheGenesiBlock_ThenCreateAnotherBlock_TheOldestBlockAtTheSameHeightAsTheCurrentMaxHeightBlockShouldBecomeTheMaxHeightBlock(t *testing.T) {
	privateKeyBob, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyAlice := &privateKeyAlice.PublicKey

	privateKeyCyril, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyCyril := &privateKeyCyril.PublicKey

	// 1. Create Genesis block
	genesisBlock := NewBlock(nil, pubKeyAlice)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	blockA := NewBlock(genesisBlock.GetHash(), pubKeyCyril)
	blockA.Finalizee()

	result := BlockProcess(blockA)
	assert.True(t, result, "Blockchain must accept Block A with 0 transactions")

	blockB := NewBlock(genesisBlock.GetHash(), pubKeyBob)
	blockB.Finalizee()

	result = BlockProcess(blockB)
	assert.True(t, result, "Blockchain must accept Block B with 0 transactions")

	blockC := NewBlock(genesisBlock.GetHash(), pubKeyCyril)
	blockC.Finalizee()

	result = BlockProcess(blockC)
	assert.True(t, result, "Blockchain must accept Block C with 0 transactions")

	blockD := BlockCreate(pubKeyAlice)
	assert.NotNil(t, blockD, "Block D should be created with zero transaction")
	assert.Equal(t, 1, len(blockA.GetTransactions()), "Block D should contain only coinbase transaction")
	assert.Equal(t, blockA.GetHash(), blockD.GetPrevBlockHash(), "Block D should be created on top of blockA, because it the older one")

}

func TestBlockchain_createMultipleBranchesOfApproximatelyTheSameLength_AndEnsureThatNewBlocksAreAlwaysCreatedInTheCorrectBranch(t *testing.T) {
	privateKeyBob, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyAlice := &privateKeyAlice.PublicKey

	privateKeyCyril, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyCyril := &privateKeyCyril.PublicKey

	// 1. Create Genesis block
	genesisBlock := NewBlock(nil, pubKeyAlice)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	blockA := NewBlock(genesisBlock.GetHash(), pubKeyCyril)
	blockA.Finalizee()

	result := BlockProcess(blockA)
	assert.True(t, result, "Blockchain must accept Block A with 0 transactions")

	blockB := NewBlock(genesisBlock.GetHash(), pubKeyBob)
	blockB.Finalizee()

	result = BlockProcess(blockB)
	assert.True(t, result, "Blockchain must accept Block B with 0 transactions")

	blockC := NewBlock(genesisBlock.GetHash(), pubKeyCyril)
	blockC.Finalizee()

	result = BlockProcess(blockC)
	assert.True(t, result, "Blockchain must accept Block C with 0 transactions")

	blockA1 := NewBlock(blockA.GetHash(), pubKeyCyril)
	blockA1.Finalizee()

	result = BlockProcess(blockA1)
	assert.True(t, result, "Blockchain must accept Block A1 with 0 transactions")

	blockB1 := NewBlock(blockB.GetHash(), pubKeyBob)
	blockB1.Finalizee()

	result = BlockProcess(blockB1)
	assert.True(t, result, "Blockchain must accept Block B1 with 0 transactions")

	blockC1 := NewBlock(blockC.GetHash(), pubKeyCyril)
	blockC1.Finalizee()

	result = BlockProcess(blockC1)
	assert.True(t, result, "Blockchain must accept Block C1 with 0 transactions")

	assert.Equal(t, 3, len(blockchain.MaxHeightNode), "Blockchain must contain 3 block with max height")

	blockB2 := NewBlock(blockB1.GetHash(), pubKeyBob)
	blockB2.Finalizee()

	result = BlockProcess(blockB2)
	assert.True(t, result, "Blockchain must accept Block B2 with 0 transactions")

	blockB3 := BlockCreate(pubKeyAlice)
	assert.NotNil(t, blockB3, "Block D should be created with zero transaction")
	assert.Equal(t, 1, len(blockA.GetTransactions()), "Block D should contain only coinbase transaction")
	assert.Equal(t, blockB2.GetHash(), blockB3.GetPrevBlockHash(), "Block D should be created on top of blockA, because it the older one")
}

func TestBlockchain_similarToThePreviousTestButThenTryToProcessBlocksWhoseParentsAreAtHeight_MaxHeightMinusCUT_OFF_AGE(t *testing.T) {
	privateKeyBob, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyBob := &privateKeyBob.PublicKey

	privateKeyAlice, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyAlice := &privateKeyAlice.PublicKey

	privateKeyCyril, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKeyCyril := &privateKeyCyril.PublicKey

	// 1. Create Genesis block
	genesisBlock := NewBlock(nil, pubKeyAlice)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	blockA := NewBlock(genesisBlock.GetHash(), pubKeyCyril)
	blockA.Finalizee()

	result := BlockProcess(blockA)
	assert.True(t, result, "Blockchain must accept Block A with 0 transactions")

	blockB := NewBlock(genesisBlock.GetHash(), pubKeyBob)
	blockB.Finalizee()

	result = BlockProcess(blockB)
	assert.True(t, result, "Blockchain must accept Block B with 0 transactions")

	blockC := NewBlock(genesisBlock.GetHash(), pubKeyCyril)
	blockC.Finalizee()

	result = BlockProcess(blockC)
	assert.True(t, result, "Blockchain must accept Block C with 0 transactions")

	blockA1 := NewBlock(blockA.GetHash(), pubKeyCyril)
	blockA1.Finalizee()

	result = BlockProcess(blockA1)
	assert.True(t, result, "Blockchain must accept Block A1 with 0 transactions")

	blockB1 := NewBlock(blockB.GetHash(), pubKeyBob)
	blockB1.Finalizee()

	result = BlockProcess(blockB1)
	assert.True(t, result, "Blockchain must accept Block B1 with 0 transactions")

	blockC1 := NewBlock(blockC.GetHash(), pubKeyCyril)
	blockC1.Finalizee()

	result = BlockProcess(blockC1)
	assert.True(t, result, "Blockchain must accept Block C1 with 0 transactions")

	assert.Equal(t, 3, len(blockchain.MaxHeightNode), "Blockchain must contain 3 block with max height")

	blockB2 := NewBlock(blockB1.GetHash(), pubKeyBob)
	blockB2.Finalizee()

	result = BlockProcess(blockB2)
	assert.True(t, result, "Blockchain must accept Block B2 with 0 transactions")

	blockB3 := BlockCreate(pubKeyAlice)
	assert.NotNil(t, blockB3, "Block D should be created with zero transaction")
	assert.Equal(t, 1, len(blockA.GetTransactions()), "Block D should contain only coinbase transaction")
	assert.Equal(t, blockB2.GetHash(), blockB3.GetPrevBlockHash(), "Block D should be created on top of blockA, because it the older one")

	prevHash := make([]byte, len(blockB3.GetHash()))
	copy(prevHash, genesisBlock.GetHash())

	for i := 0; i <= CUT_OFF_AGE+1; i++ {
		block := NewBlock(prevHash, pubKeyAlice)
		block.Finalizee()

		result = BlockProcess(block)
		assert.True(t, result, fmt.Sprintf("Block #%d should be accepted", i+1))

		prevHash = make([]byte, len(block.GetHash()))
		copy(prevHash, block.GetHash())
	}

	oldBlock := NewBlock(genesisBlock.GetHash(), pubKeyAlice)
	oldBlock.Finalizee()

	result = BlockProcess(oldBlock)
	assert.False(t, result, "Block should be rejected because its parent is too deep in history (CUT_OFF_AGE exceeded)")
}

func TestBlockchain_CreateMultisig(t *testing.T) {
	privateKey1, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKey1 := &privateKey1.PublicKey

	privateKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKey2 := &privateKey2.PublicKey

	privateKey3, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKey3 := &privateKey3.PublicKey

	genesisBlock := NewBlock(nil, pubKey2)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	block1 := NewBlock(genesisBlock.GetHash(), pubKey1)

	tx1 := NewTransaction()
	tx1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	addresses := []*rsa.PublicKey{pubKey1, pubKey2, pubKey3}
	multiSigOut := NewMultiSigOutput(3.0, addresses)
	tx1.AddMultisigOutput(multiSigOut)
	tx1.Finalize()

	tx1.SignTx(privateKey2, 0)

	block1.TransactionAdd(tx1)
	block1.Finalizee()

	result := BlockProcess(block1)

	assert.True(t, result, "Block with multisig output transaction should be accepted")
}

func TestBlockchain_CreateTransactionUsingMultisigAndSignByOneUser(t *testing.T) {
	privateKey1, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKey1 := &privateKey1.PublicKey

	privateKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKey2 := &privateKey2.PublicKey

	privateKey3, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKey3 := &privateKey3.PublicKey

	genesisBlock := NewBlock(nil, pubKey2)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	block1 := NewBlock(genesisBlock.GetHash(), pubKey1)

	tx1 := NewTransaction()
	tx1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	addresses := []*rsa.PublicKey{pubKey1, pubKey2, pubKey3}
	multiSigOut := NewMultiSigOutput(3.0, addresses)
	tx1.AddMultisigOutput(multiSigOut)
	tx1.Finalize()

	tx1.SignTx(privateKey2, 0)

	block1.TransactionAdd(tx1)
	block1.Finalizee()

	result := BlockProcess(block1)

	assert.True(t, result, "Block with multisig output transaction should be accepted")

	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 0)
	tx2.AddOutput(1.0, pubKey1)
	tx2.SignMultiSigTx(privateKey1, 0)
	tx2.Finalize()

	block2 := NewBlock(block1.GetHash(), pubKey2)
	block2.TransactionAdd(tx2)
	block2.Finalizee()
	result = BlockProcess(block2)
	assert.False(t, result, "Block with multisig transaction having only one signature should be rejected")
}

func TestBlockchain_CreateTransactionUsingMultisigAndSignByMinimumNumberOfUser(t *testing.T) {
	privateKey1, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKey1 := &privateKey1.PublicKey

	privateKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKey2 := &privateKey2.PublicKey

	privateKey3, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKey3 := &privateKey3.PublicKey

	genesisBlock := NewBlock(nil, pubKey2)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	block1 := NewBlock(genesisBlock.GetHash(), pubKey1)

	tx1 := NewTransaction()
	tx1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	addresses := []*rsa.PublicKey{pubKey1, pubKey2, pubKey3}
	multiSigOut := NewMultiSigOutput(3.0, addresses)
	tx1.AddMultisigOutput(multiSigOut)
	tx1.Finalize()

	tx1.SignTx(privateKey2, 0)

	block1.TransactionAdd(tx1)
	block1.Finalizee()

	result := BlockProcess(block1)

	assert.True(t, result, "Block with multisig output transaction should be accepted")

	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 0)
	tx2.AddOutput(2.0, pubKey1)
	// Sign with Alice and Carol.
	tx2.SignMultiSigTx(privateKey1, 0)
	tx2.SignMultiSigTx(privateKey3, 0)
	tx2.Finalize()

	block2 := NewBlock(block1.GetHash(), pubKey2)
	block2.TransactionAdd(tx2)
	block2.Finalizee()
	result = BlockProcess(block2)
	assert.True(t, result, "Block with multisig transaction having two signatures should be accepted")
}

func TestBlockchain_FullMultiSigTransactions(t *testing.T) {
	privateKey1, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKey1 := &privateKey1.PublicKey

	privateKey2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKey2 := &privateKey2.PublicKey

	privateKey3, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKey3 := &privateKey3.PublicKey

	genesisBlock := NewBlock(nil, pubKey2)
	genesisBlock.Finalizee()

	localBlockchain := NewBlockchain(genesisBlock)
	HandleBlocks(localBlockchain)

	block1 := NewBlock(genesisBlock.GetHash(), pubKey1)

	tx1 := NewTransaction()
	tx1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)
	addresses := []*rsa.PublicKey{pubKey1, pubKey2, pubKey3}
	multiSigOut := NewMultiSigOutput(3.0, addresses)
	tx1.AddMultisigOutput(multiSigOut)
	tx1.Finalize()

	tx1.SignTx(privateKey2, 0)

	block1.TransactionAdd(tx1)
	block1.Finalizee()

	result := BlockProcess(block1)

	assert.True(t, result, "Block with multisig output transaction should be accepted")

	tx2 := NewTransaction()
	tx2.AddInput(tx1.GetHash(), 0)
	tx2.AddOutput(1.0, pubKey1)
	tx2.SignMultiSigTx(privateKey1, 0)
	tx2.Finalize()

	block2 := NewBlock(block1.GetHash(), pubKey2)
	block2.TransactionAdd(tx2)
	block2.Finalizee()
	result = BlockProcess(block2)
	assert.False(t, result, "Block with multisig transaction having only one signature should be rejected")

	tx3 := NewTransaction()
	tx3.AddInput(tx1.GetHash(), 0)
	tx3.AddOutput(2.0, pubKey1)
	// Sign with Alice and Carol.
	tx3.SignMultiSigTx(privateKey1, 0)
	tx3.SignMultiSigTx(privateKey3, 0)
	tx3.Finalize()

	block3 := NewBlock(block1.GetHash(), pubKey2)
	block3.TransactionAdd(tx3)
	block3.Finalizee()
	result = BlockProcess(block3)
	assert.True(t, result, "Block with multisig transaction having two signatures should be accepted")
}
