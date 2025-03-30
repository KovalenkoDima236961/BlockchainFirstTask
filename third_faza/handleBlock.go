package third_faza

import "crypto/rsa"

var (
	blockchain *Blockchain
)

func HandleBlocks(blockChain *Blockchain) {
	blockchain = blockChain
}

func BlockProcess(block *Block) bool {
	if block == nil {
		return false
	}
	res := blockchain.BlockAdd(block)
	return res
}

func BlockCreate(myAddress *rsa.PublicKey) *Block {
	parent := blockchain.GetBlockAtMaxHeight()
	parentHash := append([]byte{}, parent.GetHash()...)

	current := NewBlock(parentHash, myAddress)
	uPool := blockchain.GetUTXOPoolAtMaxHeight()
	txPool := blockchain.GetTransactionPool()

	HandleTxs(uPool)
	txs := make([]*Transaction, 0)
	txs = txPool.GetTransactions()
	rTxs := Handler(txs)

	for _, tx := range rTxs {
		current.TransactionAdd(tx)
	}

	current.Finalizee()
	if blockchain.BlockAdd(current) {
		return current
	} else {
		return nil
	}
}

func TxProcess(tx *Transaction) {
	blockchain.TransactionAdd(tx)
}
