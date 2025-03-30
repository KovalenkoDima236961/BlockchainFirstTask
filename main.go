package main

import (
	"DMBLOCK_GO/third_faza"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strconv"
	"strings"
)

type KeyPair struct {
	Name       string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// Global variables
var (
	blockchain *third_faza.Blockchain
	keyPairs   []KeyPair

	txInputs  []TxInputData
	txOutputs []TxOutputData

	mainWindow fyne.Window

	shortTxHashToLongTxHash map[string]string
)

// ===================== STYLED BLOCKS / TREE VIEW =====================

func styledBlockCard(height int, hash string, prevHash string, txCount int) fyne.CanvasObject {
	// Background
	bg := canvas.NewRectangle(color.NRGBA{R: 60, G: 60, B: 100, A: 255})
	bg.SetMinSize(fyne.NewSize(200, 180))

	// White text
	heightText := canvas.NewText(fmt.Sprintf("Height: %d", height), color.White)
	heightText.Alignment = fyne.TextAlignCenter
	heightText.TextStyle.Bold = true

	hashText := canvas.NewText(fmt.Sprintf("Hash: %.6s", hash), color.White)
	hashText.Alignment = fyne.TextAlignCenter

	prevHashText := canvas.NewText(fmt.Sprintf("Prev Hash: %.6s", prevHash), color.White)
	prevHashText.Alignment = fyne.TextAlignCenter

	txsText := canvas.NewText(fmt.Sprintf("%d txs", txCount), color.White)
	txsText.Alignment = fyne.TextAlignCenter

	content := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(heightText),
		container.NewCenter(hashText),
		container.NewCenter(prevHashText),
		container.NewCenter(txsText),
		layout.NewSpacer(),
	)

	return container.NewMax(bg, content)
}

func buildTree(node *third_faza.BlockNode) *fyne.Container {
	block := node.B
	card := styledBlockCard(int(node.Height), fmt.Sprintf("%x", block.GetHash()), fmt.Sprintf("%x", block.GetPrevBlockHash()), len(block.GetTransactions()))

	childContainers := []fyne.CanvasObject{}
	for _, child := range node.Children {
		childContainers = append(childContainers, buildTree(child))
	}

	switch len(childContainers) {
	case 0:
		// No children => Just center the card
		return container.NewVBox(
			container.NewCenter(card),
		)

	case 1:
		// Exactly one child => center the parent and the single child
		return container.NewVBox(
			container.NewCenter(card),
			container.NewCenter(childContainers[0]),
		)

	default:
		// Multiple children => place them in an HBox with spacing
		row := container.NewHBox()
		for i, c := range childContainers {
			if i > 0 {
				// Add a spacer between children to create visible separation
				row.Add(layout.NewSpacer())
			}
			row.Add(c)
		}

		// Center the row of children
		centeredChildren := container.NewCenter(row)

		return container.NewVBox(
			container.NewCenter(card),
			centeredChildren,
		)
	}
}

func buildBlockchainTreeView() fyne.CanvasObject {
	root := getGenesisNode()
	if root == nil {
		return widget.NewLabel("❌ Genesis block not found")
	}

	tree := buildTree(root)
	scroll := container.NewScroll(tree)
	scroll.SetMinSize(fyne.NewSize(800, 600)) // big enough for large blocks
	return scroll
}

func getGenesisNode() *third_faza.BlockNode {
	for _, node := range blockchain.BlockChain {
		if node.Parent == nil {
			return node
		}
	}
	return nil
}

// ===================== ADD TRANSACTION SCREEN (DYNAMIC) =====================

// Data structures for storing user-provided inputs/outputs in the UI
type UTXOInfo struct {
	TxHashHex string
	Index     int
	Value     float64
}

func getUTXOsForKey(pubKey *rsa.PublicKey) []UTXOInfo {
	results := []UTXOInfo{}

	utxoPool := blockchain.GetUTXOPoolAtMaxHeight()
	if utxoPool == nil {
		return results
	}

	for utxoKeyStr, txOut := range utxoPool.H {
		// utxoKeyStr is e.g. "abcdef1234:1"
		// We'll parse out hashHex = "abcdef1234", and index = 1
		split := strings.Split(utxoKeyStr, ":")
		if len(split) != 2 {
			continue // or handle error
		}
		hashHex := split[0]
		indexStr := split[1]

		index, err := strconv.Atoi(indexStr)
		if err != nil {
			continue // or handle error
		}

		// If it matches the user's pubKey
		if txOut.Address == pubKey {
			results = append(results, UTXOInfo{
				TxHashHex: hashHex,
				Index:     index,
				Value:     txOut.Value,
			})
		}
	}

	return results
}

type TxInputData struct {
	TxHash string
	Index  int
}

type TxOutputData struct {
	RecipientName string
	Amount        float64
}

func parseUTXOString(sel string) (hashHex string, index int, amount float64, err error) {
	// Expected format: "Tx:abcd12[1] => 10.0 coins"
	// 1. Remove the "Tx:" prefix.
	if !strings.HasPrefix(sel, "Tx:") {
		return "", 0, 0, fmt.Errorf("string does not start with 'Tx:'")
	}
	s := strings.TrimPrefix(sel, "Tx:")

	// 2. Find the position of the '[' character.
	startBracket := strings.Index(s, "[")
	if startBracket == -1 {
		return "", 0, 0, fmt.Errorf("missing '['")
	}
	hashHex = s[:startBracket]

	// 3. Find the closing ']' to extract the index.
	endBracket := strings.Index(s, "]")
	if endBracket == -1 || endBracket <= startBracket {
		return "", 0, 0, fmt.Errorf("missing or misplaced ']'")
	}
	indexStr := s[startBracket+1 : endBracket]
	index, err = strconv.Atoi(indexStr)
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid index: %w", err)
	}

	// 4. Find "=>" to locate the amount.
	arrowIndex := strings.Index(s, "=>")
	if arrowIndex == -1 {
		return "", 0, 0, fmt.Errorf("missing '=>'")
	}
	amountStr := s[arrowIndex+2:]
	// Remove "coins" suffix and trim spaces.
	amountStr = strings.TrimSpace(strings.TrimSuffix(amountStr, "coins"))
	amount, err = strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid amount: %w", err)
	}
	return hashHex, index, amount, nil
}

func buildAddTransactionScreen() fyne.CanvasObject {
	inputsListContainer := container.NewVBox()
	outputsListContainer := container.NewVBox()

	// --------------------- FROM KEY SELECT ---------------------
	var fromKeySelect *widget.Select
	var utxoSelect *widget.Select

	keyNames := make([]string, len(keyPairs))
	for i, kp := range keyPairs {
		keyNames[i] = kp.Name
	}

	fromKeySelect = widget.NewSelect(keyNames, func(chosen string) {
		// When a sender is selected, find its key pair.
		var fromKP *KeyPair
		for i := range keyPairs {
			if keyPairs[i].Name == chosen {
				fromKP = &keyPairs[i]
				break
			}
		}
		if fromKP == nil {
			return
		}
		// Get UTXOs for this key.
		utxos := getUTXOsForKey(fromKP.PublicKey)
		utxoNames := []string{}
		for _, u := range utxos {
			// Example. "Tx:abcd12[1] => 10.0 coins"
			s := fmt.Sprintf("Tx:%s[%d] => %.2f coins", u.TxHashHex[:6], u.Index, u.Value)
			utxoNames = append(utxoNames, s)
			shortTxHashToLongTxHash[u.TxHashHex[:6]] = u.TxHashHex
		}
		utxoSelect.Options = utxoNames
		utxoSelect.Refresh()
	})
	fromKeySelect.PlaceHolder = "Select From Key"

	// --------------------- UTXO SELECT & ADD BUTTON ---------------------
	utxoSelect = widget.NewSelect([]string{}, func(chosen string) {})
	utxoSelect.PlaceHolder = "Pick a UTXO"

	utxoAddStatus := widget.NewLabel("")
	addUtxoBtn := widget.NewButton("Add Selected UTXO", func() {
		sel := utxoSelect.Selected
		if sel == "" {
			utxoAddStatus.SetText("Please select a UTXO first.")
			return
		}

		var hashHex string
		var index int
		var amount float64
		fmt.Println(utxoSelect.Selected)
		hashHex, index, amount, err := parseUTXOString(sel)
		if err != nil {
			utxoAddStatus.SetText("Failed to parse UTXO string: " + err.Error())
			return
		}
		txInputs = append(txInputs, TxInputData{
			TxHash: shortTxHashToLongTxHash[hashHex], // У вас буде повнаhash
			Index:  index,
		})
		utxoAddStatus.SetText(fmt.Sprintf("Added input: %s[%d] (%.2f coins)", hashHex, index, amount))

		inputsListContainer.Objects = nil
		for _, inp := range txInputs {
			inputsListContainer.Add(widget.NewLabel(fmt.Sprintf("Input: %s[%d]", inp.TxHash[:6], inp.Index)))
		}
		inputsListContainer.Refresh()
	})

	// --------------------- OUTPUT SECTION ---------------------
	// Recipient Key
	toKeySelect := widget.NewSelect(keyNames, func(chosen string) {})
	toKeySelect.PlaceHolder = "Recipient Key"

	outputAmountEntry := widget.NewEntry()
	outputAmountEntry.SetPlaceHolder("Amount (e.g. 5.0)")

	addOutputStatus := widget.NewLabel("")
	addOutputBtn := widget.NewButton("Add Output", func() {
		if toKeySelect.Selected == "" {
			addOutputStatus.SetText("Error: select recipient key")
			return
		}
		amtStr := outputAmountEntry.Text
		if amtStr == "" {
			addOutputStatus.SetText("Error: enter amount")
			return
		}
		amt, err := strconv.ParseFloat(amtStr, 64)
		if err != nil || amt <= 0 {
			addOutputStatus.SetText("Error: invalid amount")
			return
		}
		txOutputs = append(txOutputs, TxOutputData{
			RecipientName: toKeySelect.Selected,
			Amount:        amt,
		})
		addOutputStatus.SetText(fmt.Sprintf("Added Output: %s => %.2f", toKeySelect.Selected, amt))

		outputsListContainer.Objects = nil
		for _, outp := range txOutputs {
			outputsListContainer.Add(widget.NewLabel(fmt.Sprintf("Output: %s => %.2f", outp.RecipientName, outp.Amount)))
		}
		outputsListContainer.Refresh()

		// clear
		outputAmountEntry.SetText("")
		toKeySelect.ClearSelected()
	})

	// --------------------- CREATE & SIGN BUTTON ---------------------
	statusLabel := widget.NewLabel("")
	createTxBtn := widget.NewButton("Create & Sign TX", func() {
		if fromKeySelect.Selected == "" {
			statusLabel.SetText("No From Key selected!")
			return
		}
		if len(txInputs) == 0 {
			statusLabel.SetText("No inputs added!")
			return
		}
		if len(txOutputs) == 0 {
			statusLabel.SetText("No outputs added!")
			return
		}

		// 1) Знайти fromKeyPair
		var fromKP *KeyPair
		for i := range keyPairs {
			if keyPairs[i].Name == fromKeySelect.Selected {
				fromKP = &keyPairs[i]
			}
		}
		if fromKP == nil {
			statusLabel.SetText("FromKey not found.")
			return
		}

		// 2) Створити транзакцію
		tx := third_faza.NewTransaction()

		// 3) Додати inputs
		for _, inp := range txInputs {
			// У реальному випадку треба повний txHash + parseHexString
			// Зараз у нас лише demo
			hashBytes := parseHexString(inp.TxHash) // parse "abcd12" => 3 bytes
			tx.AddInput(hashBytes, inp.Index)
		}

		// 4) Додати outputs
		for _, outp := range txOutputs {
			var rKP *KeyPair
			for i := range keyPairs {
				if keyPairs[i].Name == outp.RecipientName {
					rKP = &keyPairs[i]
				}
			}
			if rKP == nil {
				statusLabel.SetText("Output recipient not found: " + outp.RecipientName)
				return
			}
			tx.AddOutput(outp.Amount, rKP.PublicKey)
		}

		// 5) Підписуємо всі інпути
		for i := range tx.Inputs {
			tx.SignTx(fromKP.PrivateKey, i)
		}

		// 6) Додаємо в пул
		third_faza.TxProcess(tx)

		statusLabel.SetText("Transaction created & signed. Inputs/Outputs cleared.")
		txInputs = nil
		txOutputs = nil
		// Clear the lists.
		inputsListContainer.Objects = nil
		outputsListContainer.Objects = nil
		inputsListContainer.Refresh()
		outputsListContainer.Refresh()
	})

	// --------------------- Layout ---------------------
	// Layout the input section.
	fromKeyBox := container.NewVBox(
		widget.NewLabelWithStyle("From Key (UTXO owner)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		fromKeySelect,
	)
	utxoBox := container.NewVBox(
		widget.NewLabelWithStyle("Pick UTXO", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		utxoSelect,
		addUtxoBtn,
		utxoAddStatus,
		// Display the list of added inputs.
		widget.NewLabelWithStyle("Added Inputs:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		inputsListContainer,
	)
	inputSection := container.NewVBox(fromKeyBox, utxoBox)

	// Layout the output section.
	outputBox := container.NewVBox(
		widget.NewLabelWithStyle("Add Output", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		toKeySelect,
		outputAmountEntry,
		addOutputBtn,
		addOutputStatus,
		// Display the list of added outputs.
		widget.NewLabelWithStyle("Added Outputs:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		outputsListContainer,
	)

	form := container.NewVBox(
		widget.NewLabelWithStyle("➕ Create a Custom Transaction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewHBox(inputSection, layout.NewSpacer(), outputBox),
		createTxBtn,
		statusLabel,
	)
	return form
}

// parseHexString is a helper to convert hex string → []byte
func parseHexString(str string) []byte {
	data := make([]byte, len(str)/2)
	for i := 0; i < len(data); i++ {
		fmt.Sscanf(str[2*i:2*i+2], "%x", &data[i])
	}
	return data
}

// ===================== MAIN & INIT =====================

func init() {
	// Generate 3 sample keys (User1, User2, User3)
	for i := 1; i <= 3; i++ {
		priv, _ := rsa.GenerateKey(rand.Reader, 1024)
		pub := &priv.PublicKey
		keyPairs = append(keyPairs, KeyPair{
			Name:       fmt.Sprintf("User%d", i),
			PrivateKey: priv,
			PublicKey:  pub,
		})
	}
	shortTxHashToLongTxHash = make(map[string]string)
}

func main() {
	myApp := app.New()
	mainWindow = myApp.NewWindow("Blockchain Visualizer")
	mainWindow.Resize(fyne.NewSize(800, 600))

	// 1) Create genesis block
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	publicKey := &privateKey.PublicKey
	genesis := third_faza.NewBlock(nil, publicKey)
	genesis.Finalizee()
	blockchain = third_faza.NewBlockchain(genesis)
	third_faza.HandleBlocks(blockchain)

	distTx := third_faza.NewTransaction()
	distTx.AddInput(genesis.GetCoinbase().GetHash(), 0)

	for _, kp := range keyPairs {
		distTx.AddOutput(1.0, kp.PublicKey)
	}
	distTx.SignTx(privateKey, 0)
	third_faza.TxProcess(distTx)

	// 3) Mine a block with the distribution TX so it's confirmed
	distBlock := third_faza.NewBlock(genesis.GetHash(), keyPairs[2].PublicKey)
	distBlock.TransactionAdd(distTx)
	distBlock.Finalizee()

	ok := third_faza.BlockProcess(distBlock)
	if !ok {
		fmt.Println("Distribution block was rejected - check if input >= output.")
	} else {
		fmt.Println("Distribution block accepted! All users have 1 coins.")
	}

	// ============== UI: Blockchain screen with Refresh ==============
	var updateBlockchainScreen func()
	blockchainScreen := container.NewVBox()
	updateBlockchainScreen = func() {
		blockchainTree := buildBlockchainTreeView()
		refreshButton := widget.NewButton("🔄 Refresh View", func() {
			updateBlockchainScreen()
		})

		blockchainScreen.Objects = []fyne.CanvasObject{
			widget.NewLabelWithStyle("🔗 Blockchain Chain View", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			blockchainTree,
			layout.NewSpacer(),
			container.NewCenter(refreshButton),
		}
		blockchainScreen.Refresh()
	}
	updateBlockchainScreen()

	var mainContent *fyne.Container

	// Transaction pool screen update.
	updateTransactionPoolScreen := func() {
		poolView := buildTransactionPoolView()
		mainContent.Objects = []fyne.CanvasObject{poolView}
		mainContent.Refresh()
	}
	addTxScreen := buildAddTransactionScreen()
	mineBlockScreen := buildMineBlockScreen()

	// ============== HOME SCREEN ==============
	homeTitle := widget.NewLabelWithStyle("🚀 Blockchain Visualizer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	homeDesc := widget.NewLabel("Welcome! This application allows you to simulate and visualize blockchain mechanics.\n\n" +
		"- Add and validate transactions\n" +
		"- Create blocks manually\n" +
		"- Explore blockchain structure with branching and forks\n\n" +
		"Use the sidebar to navigate through the application.")
	asciiBlock := widget.NewLabel("🧱 → 🧱 → 🧱")
	getStartedBtn := widget.NewButton("👉 Get Started (View Blockchain)", func() {
		mainContent.Objects = []fyne.CanvasObject{blockchainScreen}
		mainContent.Refresh()
	})

	homeScreen := container.NewVBox(
		homeTitle,
		homeDesc,
		asciiBlock,
		getStartedBtn,
	)

	// Main content + sidebar
	mainContent = container.NewMax(homeScreen)

	viewHomeBtn := widget.NewButton("🏠 Home", func() {
		mainContent.Objects = []fyne.CanvasObject{homeScreen}
		mainContent.Refresh()
	})
	viewChainBtn := widget.NewButton("🔗 View Blockchain", func() {
		mainContent.Objects = []fyne.CanvasObject{blockchainScreen}
		mainContent.Refresh()
	})
	viewPoolBtn := widget.NewButton("📋 Transaction Pool", func() {
		updateTransactionPoolScreen()
	})
	addTxBtn := widget.NewButton("➕ Add Transaction", func() {
		mainContent.Objects = []fyne.CanvasObject{addTxScreen}
		mainContent.Refresh()
	})
	mineBlockBtn := widget.NewButton("⛏️ Mine Block", func() {
		mainContent.Objects = []fyne.CanvasObject{mineBlockScreen}
		mainContent.Refresh()
	})

	sidebar := container.NewVBox(
		viewHomeBtn,
		viewChainBtn,
		viewPoolBtn,
		addTxBtn,
		mineBlockBtn,
	)

	mainSplit := container.NewHSplit(sidebar, mainContent)
	mainSplit.Offset = 0.25

	mainWindow.SetContent(mainSplit)
	mainWindow.ShowAndRun()
}
