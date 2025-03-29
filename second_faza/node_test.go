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

func TestNode_firstParameters(t *testing.T) {
	numNodes := 100
	p_graph := .1
	p_byzantine := .15
	p_txDistribution := .01
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_secondParameters(t *testing.T) {
	numNodes := 100
	p_graph := .1
	p_byzantine := .15
	p_txDistribution := .01
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_thirdParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .15
	p_txDistribution := .05
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_fourParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .15
	p_txDistribution := .10
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_fiveParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .15
	p_txDistribution := .05
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_sixParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .15
	p_txDistribution := .10
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_sevenParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .30
	p_txDistribution := .01
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_eightParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .30
	p_txDistribution := .01
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_nineParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .30
	p_txDistribution := .05
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_tenParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .30
	p_txDistribution := .05
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_elevenParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .30
	p_txDistribution := .10
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_twelveParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .30
	p_txDistribution := .10
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_thirteenParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .45
	p_txDistribution := .01
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_fourteenParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .45
	p_txDistribution := .01
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_fifteenParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .45
	p_txDistribution := .05
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_sixteenParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .45
	p_txDistribution := .05
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_seventeenParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .45
	p_txDistribution := .10
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_eighteenParameters(t *testing.T) {
	numNodes := 100

	p_graph := .1
	p_byzantine := .45
	p_txDistribution := .10
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_nineteenParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .15
	p_txDistribution := .01
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_twentyParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .15
	p_txDistribution := .01
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_twenty_oneParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .15
	p_txDistribution := .05
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_twenty_twoParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .15
	p_txDistribution := .05
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_twenty_threeParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .15
	p_txDistribution := .10
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_twenty_fourParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .15
	p_txDistribution := .10
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_twenty_fiveParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .30
	p_txDistribution := .01
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_twenty_sixParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .30
	p_txDistribution := .01
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_twenty_sevenParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .30
	p_txDistribution := .05
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_twenty_eightParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .30
	p_txDistribution := .05
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_twenty_nineParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .30
	p_txDistribution := .10
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_thirtyParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .30
	p_txDistribution := .10
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_thirty_oneParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .45
	p_txDistribution := .01
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_thirty_twoParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .45
	p_txDistribution := .01
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_thirty_threeParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .45
	p_txDistribution := .05
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_thirty_fourParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .45
	p_txDistribution := .05
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_thirty_fiveParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .45
	p_txDistribution := .10
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_thirty_sixParameters(t *testing.T) {
	numNodes := 100

	p_graph := .2
	p_byzantine := .45
	p_txDistribution := .10
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_thirty_sevenParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .15
	p_txDistribution := .01
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_thirty_eightParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .15
	p_txDistribution := .01
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_thirty_nineParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .15
	p_txDistribution := .05
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_fortyParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .15
	p_txDistribution := .05
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_forty_oneParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .15
	p_txDistribution := .10
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_forty_twoParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .15
	p_txDistribution := .10
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_forty_threeParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .30
	p_txDistribution := .01
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_forty_fourParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .30
	p_txDistribution := .01
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_forty_fiveParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .30
	p_txDistribution := .05
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_forty_sixParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .30
	p_txDistribution := .05
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_forty_sevenParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .30
	p_txDistribution := .10
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_forty_eightParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .30
	p_txDistribution := .10
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_forty_nineParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .45
	p_txDistribution := .01
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_fiftyParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .45
	p_txDistribution := .01
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_fifty_oneParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .45
	p_txDistribution := .05
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_fifty_twoParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .45
	p_txDistribution := .05
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_fifty_threeParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .45
	p_txDistribution := .10
	numRounds := 10

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}

func TestNode_fifty_fourParameters(t *testing.T) {
	numNodes := 100

	p_graph := .3
	p_byzantine := .45
	p_txDistribution := .10
	numRounds := 20

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
	}

	for i := 0; i < numNodes; i++ {
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
				if Contain(validTxsIds, tx.id) {
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

	//for i := 0; i < numNodes; i++ {
	//	transactions := nodes[i].FollowersSend()
	//	fmt.Printf("Transaction ids that Node %d believes consensus\n", &i)
	//	for _, tx := range transactions {
	//		fmt.Println(tx.id)
	//	}
	//	fmt.Println()
	//	fmt.Println()
	//}

	txCount := make(map[int]int)
	for _, node := range nodes {
		finalTxs := node.FollowersSend()
		for _, tx := range finalTxs {
			txCount[tx.id]++
		}
	}

	// Calculate threshold
	threshold := int(float64(numNodes) * (1 - p_byzantine))
	successfulTxs := 0

	for _, txId := range validTxsIds {
		count, exists := txCount[txId]
		if exists && count >= threshold {
			successfulTxs++
		}
	}

	totalValidTxs := len(validTxsIds)
	successRate := float64(successfulTxs) / float64(totalValidTxs) * 100.0
	t.Logf("Overall success rate: %.2f%%", successRate)
}
