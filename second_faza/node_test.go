package second_faza

import (
	"math/rand"
	"testing"
)

func Contain(array []int, value int) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}

	return false
}

func Helper(numNodes int, p_graph float64, p_byzantine float64, p_txDistribution float64, numRounds int, t *testing.T) {
	nodes := make([]Node, numNodes)
	for i := 0; i < numNodes; i++ {
		if rand.Float64() < p_byzantine {
			nodes[i] = CreateByzantineNode(p_graph, p_byzantine, p_txDistribution, numRounds)
		} else {
			nodes[i] = CreateTrustedNode(p_graph, p_byzantine, p_txDistribution, numRounds)
		}
	}

	followees := make([][]bool, numNodes)
	for i := 0; i < numNodes; i++ {

		followees[i] = make([]bool, numNodes)
		for j := 0; j < numNodes; j++ {
			if i == j {
				continue
			}
			if rand.Float64() < p_graph {
				followees[i][j] = true
			}
		}
	}

	// upozorni všetky uzly o ich nasledovníkoch
	for i := 0; i < numNodes; i++ {
		nodes[i].FolloweesSet(followees[i])
	}

	// inicializuj set 500 platných transakcií s náhodnými id
	numTx := 500
	validTxsIds := make([]int, numTx)
	for i := 0; i < numTx; i++ {
		r := rand.Int()
		validTxsIds = append(validTxsIds, r)
	}

	// distribuuje 500 transakcií do všetkých uzlov a inicializuje ich
	// počiatočný stav transakcií, ktoré každý uzol počul. Distribúcia
	// je náhodná s pravdepodobnosťou p_txDistribution pre každý pár
	// Transkacia-Uzol.
	for i := 0; i < numNodes; i++ {
		pendingTransaction := make([]*Transaction, 0)
		for _, txId := range validTxsIds {
			if rand.Float64() < p_txDistribution {
				pendingTransaction = append(pendingTransaction, NewTransaction(txId))
			}
			nodes[i].PendingTransactionSet(pendingTransaction)
		}
	}

	// Simuluj numRounds-krát
	for round := 0; round < numRounds; round++ { // numRounds je buď 10, alebo 20

		// zhromaždiť všetky návrhy do mapy. Kľúčom je index uzla prijímajúceho
		// návrhy. Hodnota je ArrayList obsahujúci polia celých čísel 1x2. Prvým
		// prvkom každého poľa je ID navrhovanej transakcie a druhý
		// element je indexové číslo uzla navrhujúceho transakciu.
		allProposals := make(map[int][][]int, 0)

		for i := 0; i < numNodes; i++ {
			proposals := nodes[i].FollowersSend()
			for _, tx := range proposals {
				if !Contain(validTxsIds, tx.id) {
					continue
				}

				for j := 0; j < numNodes; j++ {
					if !followees[j][i] {
						continue
					}

					if _, ok := allProposals[j]; ok {
						candidate := make([]int, 2)
						candidate[0] = tx.id
						candidate[1] = i
						allProposals[j] = append(allProposals[j], candidate)
					} else {
						candidates := make([][]int, 0)
						candidate := make([]int, 2)

						candidate[0] = tx.id
						candidate[1] = i
						candidates = append(candidates, candidate)
						allProposals[j] = candidates
					}
				}
			}
		}

		for i := 0; i < numNodes; i++ {
			if _, ok := allProposals[i]; ok {
				nodes[i].FollowesReceive(allProposals[i])
			}
		}
	}

	var referenceSet map[int]bool
	firstTruster := true
	for i, node := range nodes {
		trustedNode, ok := node.(*TrustedNode)
		if !ok {
			continue
		}

		finalTxs := trustedNode.FollowersSend()
		txSet := make(map[int]bool)
		for _, tx := range finalTxs {
			txSet[tx.HashCode()] = true
		}

		if firstTruster {
			referenceSet = txSet
			firstTruster = false
		} else {
			if !equalSets(txSet, referenceSet) {
				t.Errorf("Trusted node %d has a different transaction pool than the reference set", i)
			}
		}
	}

	// (Optional) Continue with additional checks, for example, computing a success rate.
	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.HashCode()]++
		}
	}
	threshold := numNodes
	successfulTxs := 0
	for _, txId := range validTxsIds {
		if count, exists := txCount[txId]; exists && count >= threshold {
			successfulTxs++
		}
	}
	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func equalSets(a, b map[int]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for key := range a {
		if !b[key] {
			return false
		}
	}
	return true
}

func TestNode_firstParameters(t *testing.T) {
	Helper(100, .1, .15, .01, 10, t)
}

func TestNode_secondParameters(t *testing.T) {
	Helper(100, .1, .15, .01, 20, t)
}

func TestNode_thirdParameters(t *testing.T) {
	Helper(100, .1, .15, .05, 10, t)
}

func TestNode_fourParameters(t *testing.T) {
	Helper(100, .1, .15, .05, 20, t)
}

func TestNode_fiveParameters(t *testing.T) {
	Helper(100, .1, .15, .10, 10, t)
}

func TestNode_sixParameters(t *testing.T) {
	Helper(100, .1, .15, .10, 20, t)
}

func TestNode_sevenParameters(t *testing.T) {
	Helper(100, .1, .30, .01, 10, t)
}

func TestNode_eightParameters(t *testing.T) {
	Helper(100, .1, .30, .01, 20, t)
}

func TestNode_nineParameters(t *testing.T) {
	Helper(100, .1, .30, .05, 10, t)
}

func TestNode_tenParameters(t *testing.T) {
	Helper(100, .1, .30, .05, 20, t)
}

func TestNode_elevenParameters(t *testing.T) {
	Helper(100, .1, .30, .10, 10, t)
}

func TestNode_twelveParameters(t *testing.T) {
	Helper(100, .1, .30, .10, 20, t)
}

func TestNode_thirteenParameters(t *testing.T) {
	Helper(100, .1, .45, .01, 10, t)
}

func TestNode_fourteenParameters(t *testing.T) {
	Helper(100, .1, .45, .01, 20, t)
}

func TestNode_fifteenParameters(t *testing.T) {
	Helper(100, .1, .45, .05, 10, t)
}

func TestNode_sixteenParameters(t *testing.T) {
	Helper(100, .1, .45, .05, 20, t)
}

func TestNode_seventeenParameters(t *testing.T) {
	Helper(100, .1, .45, .10, 10, t)
}

func TestNode_eighteenParameters(t *testing.T) {
	Helper(100, .1, .45, .10, 20, t)
}

func TestNode_nineteenParameters(t *testing.T) {
	Helper(100, .2, .15, .01, 10, t)
}

func TestNode_twentyParameters(t *testing.T) {
	Helper(100, .2, .15, .01, 20, t)
}

func TestNode_twenty_oneParameters(t *testing.T) {
	Helper(100, .2, .15, .05, 10, t)
}

func TestNode_twenty_twoParameters(t *testing.T) {
	Helper(100, .2, .15, .05, 20, t)
}

func TestNode_twenty_threeParameters(t *testing.T) {
	Helper(100, .2, .15, .10, 10, t)
}

func TestNode_twenty_fourParameters(t *testing.T) {
	Helper(100, .2, .15, .10, 20, t)
}

func TestNode_twenty_fiveParameters(t *testing.T) {
	Helper(100, .2, .30, .01, 10, t)
}

func TestNode_twenty_sixParameters(t *testing.T) {
	Helper(100, .2, .30, .01, 20, t)
}

func TestNode_twenty_sevenParameters(t *testing.T) {
	Helper(100, .2, .30, .05, 10, t)
}

func TestNode_twenty_eightParameters(t *testing.T) {
	Helper(100, .2, .30, .05, 20, t)
}

func TestNode_twenty_nineParameters(t *testing.T) {
	Helper(100, .2, .30, .10, 10, t)
}

func TestNode_thirtyParameters(t *testing.T) {
	Helper(100, .2, .30, .10, 20, t)
}

func TestNode_thirty_oneParameters(t *testing.T) {
	Helper(100, .2, .45, .01, 10, t)
}

func TestNode_thirty_twoParameters(t *testing.T) {
	Helper(100, .2, .45, .01, 20, t)
}

func TestNode_thirty_threeParameters(t *testing.T) {
	Helper(100, .2, .45, .05, 10, t)
}

func TestNode_thirty_fourParameters(t *testing.T) {
	Helper(100, .2, .45, .05, 20, t)
}

func TestNode_thirty_fiveParameters(t *testing.T) {
	Helper(100, .2, .45, .10, 10, t)
}

func TestNode_thirty_sixParameters(t *testing.T) {
	Helper(100, .2, .45, .10, 20, t)
}

func TestNode_thirty_sevenParameters(t *testing.T) {
	Helper(100, .3, .15, .01, 10, t)
}

func TestNode_thirty_eightParameters(t *testing.T) {
	Helper(100, .3, .15, .01, 20, t)
}

func TestNode_thirty_nineParameters(t *testing.T) {
	Helper(100, .3, .15, .05, 10, t)
}

func TestNode_fortyParameters(t *testing.T) {
	Helper(100, .3, .15, .05, 20, t)
}

func TestNode_forty_oneParameters(t *testing.T) {
	Helper(100, .3, .15, .10, 10, t)
}

func TestNode_forty_twoParameters(t *testing.T) {
	Helper(100, .3, .15, .10, 20, t)
}

func TestNode_forty_threeParameters(t *testing.T) {
	Helper(100, .3, .30, .01, 10, t)
}

func TestNode_forty_fourParameters(t *testing.T) {
	Helper(100, .3, .30, .01, 20, t)
}

func TestNode_forty_fiveParameters(t *testing.T) {
	Helper(100, .3, .30, .05, 10, t)
}

func TestNode_forty_sixParameters(t *testing.T) {
	Helper(100, .3, .30, .05, 20, t)
}

func TestNode_forty_sevenParameters(t *testing.T) {
	Helper(100, .3, .30, .10, 10, t)
}

func TestNode_forty_eightParameters(t *testing.T) {
	Helper(100, .3, .30, .10, 20, t)
}

func TestNode_forty_nineParameters(t *testing.T) {
	Helper(100, .3, .45, .01, 10, t)
}

func TestNode_fiftyParameters(t *testing.T) {
	Helper(100, .3, .45, .01, 20, t)
}

func TestNode_fifty_oneParameters(t *testing.T) {
	Helper(100, .3, .45, .05, 10, t)
}

func TestNode_fifty_twoParameters(t *testing.T) {
	Helper(100, .3, .45, .05, 20, t)
}

func TestNode_fifty_threeParameters(t *testing.T) {
	Helper(100, .3, .45, .10, 10, t)
}

func TestNode_fifty_fourParameters(t *testing.T) {
	Helper(100, .3, .45, .10, 20, t)
}
