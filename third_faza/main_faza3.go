package third_faza

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
)

func mainD() {
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

	privateKeyCyril, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyCyril := &privateKeyCyril.PublicKey

	genesisBlock := NewBlock(nil, pubKeyBob)
	genesisBlock.Finalizee()
	blockchain = NewBlockchain(genesisBlock)
	HandleBlocks(blockchain)

	block1 := NewBlock(genesisBlock.GetHash(), pubKeyAlice)

	tx1 := NewTransaction()
	tx1.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)

	tx1.AddOutput(1, pubKeyAlice)
	tx1.AddOutput(1, pubKeyAlice)
	tx1.AddOutput(1.125, pubKeyAlice)

	// Je len jeden (na pozicii 0) Transaction.Input v tx1
	// a ten obsahuje mince od Boba, a preto je potrebne podpisat transakciu privatnym klucom Boba
	tx1.SignTx(privateKeyBob, 0)

	block1.TransactionAdd(tx1)
	block1.Finalizee()

	println("Block1 Valid check: ", BlockProcess(block1))

	// Bobove vytvorenit alternativneho blocku2 s transakciou sam seve
	block2 := NewBlock(genesisBlock.GetHash(), pubKeyBob)

	tx2 := NewTransaction()

	tx2.AddInput(genesisBlock.GetCoinbase().GetHash(), 0)

	tx2.AddOutput(1, pubKeyBob)
	tx2.AddOutput(1, pubKeyBob)
	tx2.AddOutput(1.125, pubKeyBob)

	tx2.SignTx(privateKeyBob, 0)

	block2.TransactionAdd(tx2)
	block2.Finalizee()

	println("Block2 Valid check: ", BlockProcess(block2))

	// Bobove vytvorenie blocku3, ktory retazi na block1, s transakciou od Alici Cyrilovi
	block3 := NewBlock(genesisBlock.GetHash(), pubKeyBob)

	tx3 := NewTransaction()
	// posielaju sa coiny v hodnote 2
	tx3.AddInput(tx1.GetHash(), 0)
	tx3.AddInput(tx1.GetHash(), 1)
	// tx3.AddInput(tx1.GetHash(), 2)

	tx3.AddOutput(2, pubKeyCyril)

	tx3.SignTx(privateKeyAlice, 0)
	tx3.SignTx(privateKeyAlice, 1)
	//tx3.SignTx(privateKeyAlice, 1)

	block3.TransactionAdd(tx3)
	block3.Finalizee()

	println("Block3 Valid check: ", BlockProcess(block3))

	// Bobove vytvorenit blocku4, ktory retazi na block3 s transakciou od Cyrila Bobovi
	block4 := NewBlock(genesisBlock.GetHash(), pubKeyBob)

	tx4 := NewTransaction()

	tx4.AddInput(tx3.GetHash(), 0)

	tx4.AddOutput(1.5, pubKeyBob)
	tx4.AddOutput(0.5, pubKeyBob)

	tx4.SignTx(privateKeyCyril, 0)

	block4.TransactionAdd(tx4)
	block4.Finalizee()

	println("Block4 Valid check: ", BlockProcess(block4))
}
