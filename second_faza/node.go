package second_faza

import (
	"math/rand"
)

var (
	k     int
	alpha int
	beta  int
)

type Node interface {
	FolloweesSet(followees []bool)
	PendingTransactionSet(pendingTransaction []*Transaction)
	FollowersSend() []*Transaction
	FollowesReceive(candidates [][]int)
}

type ByzantineNode struct {
}

func CreateByzantineNode(p_graph float64, p_byzantine float64, p_txDistribution float64, numRounds int) Node {
	return &ByzantineNode{}
}

func (node *ByzantineNode) FolloweesSet(followees []bool) {

}

func (node *ByzantineNode) PendingTransactionSet(pendingTransaction []*Transaction) {

}

func (node *ByzantineNode) FollowersSend() []*Transaction {
	return make([]*Transaction, 0)
}

func (node *ByzantineNode) FollowesReceive(candidates [][]int) {

}

type TrustedNode struct {
	followers         []bool
	localTransactions []*Transaction
	txPool            map[int]*Status
}

type Status struct {
	status     TypeOfStatus
	confidence int
}

func NewStatus() *Status {
	return &Status{status: None, confidence: 0}
}

const (
	None TypeOfStatus = iota // initial state
	Valid
	Invalid
)

type TypeOfStatus int

func CreateTrustedNode(p_graph float64, p_byzantine float64, p_txDistribution float64, numRounds int) Node {
	k = int((p_graph / 0.2) * 12 * (float64(numRounds) / 10.0))
	alpha = int(0.58 * float64(k) * (0.30 / p_byzantine))
	beta = int((p_txDistribution / 0.05) * 4 * (float64(numRounds) / 10.0))

	return &TrustedNode{
		followers:         make([]bool, 0),
		localTransactions: make([]*Transaction, 0),
		txPool:            make(map[int]*Status),
	}
}

func (node *TrustedNode) FolloweesSet(followees []bool) {
	node.followers = make([]bool, len(followees))
	copy(node.followers, followees)
}

func (node *TrustedNode) PendingTransactionSet(pendingTransactions []*Transaction) {
	for _, tx := range pendingTransactions {
		node.localTransactions = append(node.localTransactions, tx)
		node.txPool[tx.HashCode()] = NewStatus()
	}
}

func (node *TrustedNode) FollowersSend() []*Transaction {
	return node.localTransactions
}

func (node *TrustedNode) FollowesReceive(candidates [][]int) {
	candidateList := make([]*Candidate, 0, len(candidates))
	for _, data := range candidates {
		candidateList = append(candidateList, NewCandidate(NewTransaction(data[0]), data[1]))
	}

	votesByTx := make(map[int][]int)
	for _, candidate := range candidateList {
		if node.followers[candidate.sender] {
			txID := candidate.tx.HashCode()
			votesByTx[txID] = append(votesByTx[txID], candidate.sender)
		}
	}

	for txId, votes := range votesByTx {
		sampleVotes := votes
		if len(votes) > k {
			rand.Shuffle(len(votes), func(i, j int) {
				votes[i], votes[j] = votes[j], votes[i]
			})
			sampleVotes = votes[:k]
		}

		status, exists := node.txPool[txId]
		if !exists {
			status = NewStatus()
		}

		if len(sampleVotes) >= alpha {
			if status.status == Valid {
				status.confidence++
			} else {
				status.status = Valid
				status.confidence = 1
			}
		}
		node.txPool[txId] = status
	}

	for txId, status := range node.txPool {
		if status.status == Valid && status.confidence >= beta {
			tx := NewTransaction(txId)
			if !containsTransaction(node.localTransactions, tx) {
				node.localTransactions = append(node.localTransactions, tx)
			}
		}
	}
}

func containsTransaction(txs []*Transaction, tx *Transaction) bool {
	for _, t := range txs {
		if t.HashCode() == tx.HashCode() {
			return true
		}
	}
	return false
}
