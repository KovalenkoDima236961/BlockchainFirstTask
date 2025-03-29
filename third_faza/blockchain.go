package third_faza

import (
	"encoding/hex"
	"math"
)

const (
	CUT_OFF_AGE           = 12
	MAX_BLOCKS_INT_MEMORY = 100
)

type BlockNode struct {
	B        *Block
	Parent   *BlockNode
	Children []*BlockNode
	Height   uint
	Pool     *UTXOPool
}

func NewBlockNode(b *Block, parent *BlockNode, uPool *UTXOPool) *BlockNode {
	newBlockNode := &BlockNode{
		B:        b,
		Parent:   parent,
		Children: []*BlockNode{},
		Pool:     uPool,
	}

	height := uint(1)
	if parent != nil {
		height = parent.Height + 1
	}
	newBlockNode.Height = height

	return newBlockNode
}

func (blockNode *BlockNode) GetUTXOPoolCopy() *UTXOPool {
	return NewUTXOPoolWithPool(blockNode.Pool)
}

type Blockchain struct {
	BlockChain            map[string]*BlockNode
	MaxHeightNode         []*BlockNode
	GlobalTransactionPool *TransactionPool
	LatestBlocks          []string
}

func NewBlockchain(genesisBlock *Block) *Blockchain {
	blockchainF := new(Blockchain)
	blockchainF.BlockChain = make(map[string]*BlockNode)

	genesisUTXOPool := NewUTXOPool()
	genesisUTXOPool.Put(*NewUTXO(genesisBlock.GetCoinbase().GetHash(),
		0), *genesisBlock.GetCoinbase().GetOutput(0))

	genesisNode := NewBlockNode(genesisBlock, nil, genesisUTXOPool)

	blockchainF.BlockChain[keyFoBlock(genesisBlock.GetHash())] = genesisNode

	blockchainF.MaxHeightNode = make([]*BlockNode, 0)
	blockchainF.MaxHeightNode = append(blockchainF.MaxHeightNode, genesisNode)

	blockchainF.GlobalTransactionPool = NewTransactionPool()

	blockchainF.LatestBlocks = make([]string, 0)
	blockchainF.LatestBlocks = append(blockchainF.LatestBlocks, keyFoBlock(genesisBlock.GetHash()))

	return blockchainF
}

func keyFoBlock(blockHash []byte) string {
	wrapper := NewByteArrayWrapper(blockHash)
	return hex.EncodeToString(wrapper.contents)
}

func (blockChain *Blockchain) GetBlockAtMaxHeight() *Block {
	return blockChain.MaxHeightNode[0].B
}

func (blockChain *Blockchain) GetBlockNodeAtMaxHeight() *BlockNode {
	return blockChain.MaxHeightNode[0]
}

func (blockChain *Blockchain) GetUTXOPoolAtMaxHeight() *UTXOPool {
	return blockChain.MaxHeightNode[0].GetUTXOPoolCopy()
}

func (blockChain *Blockchain) GetTransactionPool() *TransactionPool {
	return blockChain.GlobalTransactionPool
}

func (blockChain *Blockchain) BlockAdd(block *Block) bool {
	parentHash := block.GetPrevBlockHash()
	if parentHash == nil || len(parentHash) == 0 {
		return false
	}
	parentBlock := blockChain.Get(parentHash)
	if parentBlock == nil {
		return false
	}

	newHeight := int(parentBlock.Height + 1)
	maxValidHeight := int(blockChain.MaxHeightNode[0].Height) - CUT_OFF_AGE
	if newHeight <= maxValidHeight {
		return false
	}

	utxo := parentBlock.GetUTXOPoolCopy()
	HandleTxs(utxo)

	blockTxs := block.GetTransactions()
	validTxs := Handler(blockTxs)

	if len(validTxs) != len(blockTxs) {
		return false
	}

	utxo = UTXOPoolGet()

	coinbaseTransaction := block.GetCoinbase()
	if coinbaseTransaction == nil || !CheckCoinbaseTransaction(coinbaseTransaction) {
		return false
	}

	for i, output := range coinbaseTransaction.Outputs {
		utxo.Put(UTXO{txHash: coinbaseTransaction.GetHash(), index: i}, *output)
	}

	newNode := NewBlockNode(block, parentBlock, utxo)
	blochHash := keyFoBlock(block.GetHash())

	blockChain.BlockChain[blochHash] = newNode
	parentBlock.Children = append(parentBlock.Children, newNode)
	blockChain.LatestBlocks = append(blockChain.LatestBlocks, blochHash)

	if newNode.Height > blockChain.MaxHeightNode[0].Height {
		oldMax := blockChain.MaxHeightNode[0]
		blockChain.MaxHeightNode = []*BlockNode{newNode}

		for _, child := range oldMax.Children {
			if child != newNode {
				RemoveFork(child, blockChain)
			}
		}
	} else if newNode.Height == blockChain.MaxHeightNode[0].Height {
		blockChain.MaxHeightNode = append(blockChain.MaxHeightNode, newNode)
	}

	if len(blockChain.LatestBlocks) > MAX_BLOCKS_INT_MEMORY {
		oldestBlockHash := blockChain.LatestBlocks[0]
		delete(blockChain.BlockChain, oldestBlockHash)
		blockChain.LatestBlocks = blockChain.LatestBlocks[1:]
	}

	for _, transaction := range blockTxs {
		blockChain.GlobalTransactionPool.RemoveTransaction(transaction.Hash)
	}
	return true
}

func RemoveFork(forkBlock *BlockNode, blockChain *Blockchain) {
	if forkBlock == nil {
		return
	}

	for _, child := range forkBlock.Children {
		RemoveFork(child, blockChain)
	}

	for _, tx := range forkBlock.B.GetTransactions() {
		for i := range tx.Outputs {
			utxo := UTXO{tx.GetHash(), i}
			forkBlock.Pool.RemoveUTXO(utxo)
		}
	}

	delete(blockChain.BlockChain, keyFoBlock(forkBlock.B.GetHash()))
}

func CheckCoinbaseTransaction(tx *Transaction) bool {
	if tx == nil {
		return false
	}
	coins := 0.0
	for _, output := range tx.GetOutputs() {
		coins += output.Value
	}

	const tolerance = 0.00001
	return math.Abs(coins-COINBASE) <= tolerance
}

func (blockChain *Blockchain) Get(parentHash []byte) *BlockNode {
	return blockChain.BlockChain[keyFoBlock(parentHash)]
}

func (blockChain *Blockchain) TransactionAdd(tx *Transaction) {
	utxo := blockChain.GetUTXOPoolAtMaxHeight()
	if TxIsValid(*tx, utxo) {
		blockChain.GlobalTransactionPool.AddTransaction(tx)
	}
}
