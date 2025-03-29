package main

import (
	"DMBLOCK_GO/third_faza"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// showBlockProcessPopup creates a pop-up dialog to select a parent block and transactions,
// then builds and processes a new block.
func showBlockProcessPopup(parentWindow fyne.Window) {
	if parentWindow == nil {
		fmt.Println("Parent window is nil; cannot display dialog.")
		return
	}

	// Build parent block selection options.
	// We'll iterate over the blockchain map and show a short representation.
	blockOptions := []string{}
	blockMap := make(map[string]*third_faza.BlockNode)
	for key, node := range blockchain.BlockChain {
		// Display the first 6 chars of the hash along with its height.
		shortHash := key[:6]
		option := fmt.Sprintf("Block %s (Height: %d)", shortHash, node.Height)
		blockOptions = append(blockOptions, option)
		blockMap[option] = node
	}
	parentBlockSelect := widget.NewSelect(blockOptions, nil)
	parentBlockSelect.PlaceHolder = "Select Parent Block"

	// Build transaction selection using checkboxes.
	txs := blockchain.GetTransactionPool().GetTransactions()
	txCheckboxes := []*widget.Check{}
	txMap := make(map[*widget.Check]*third_faza.Transaction)
	txContainer := container.NewVBox()
	if len(txs) == 0 {
		txContainer.Add(widget.NewLabel("No transactions in the pool."))
	} else {
		fmt.Println("I located before txs loop")
		for _, tx := range txs {
			// Display a short hash for each transaction.
			shortTxHash := fmt.Sprintf("%.6x", tx.GetHash())
			chk := widget.NewCheck(fmt.Sprintf("Tx %s", shortTxHash), nil)
			txCheckboxes = append(txCheckboxes, chk)
			txMap[chk] = tx
			txContainer.Add(chk)
		}
	}
	fmt.Println("I located after txs loop")
	// Optionally, wrap the transaction list in a scroll container.
	scrollTxs := container.NewVScroll(txContainer)
	scrollTxs.SetMinSize(fyne.NewSize(300, 150))

	// Compose the dialog content.
	fmt.Println("I located before creating container")
	content := container.NewVBox(
		widget.NewLabelWithStyle("Select Parent Block", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		parentBlockSelect,
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Select Transactions to Include", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		scrollTxs,
	)

	// Create the custom confirmation dialog.
	dialog.ShowCustomConfirm("Create New Block", "Create", "Cancel", content, func(confirm bool) {
		fmt.Println("I located before first if")
		if !confirm {
			return
		}

		fmt.Println("I located before second if")
		if parentBlockSelect.Selected == "" {
			dialog.ShowInformation("Error", "Please select a parent block", parentWindow)
			return
		}

		fmt.Println("I take parent node")
		parentNode, ok := blockMap[parentBlockSelect.Selected]
		if !ok || parentNode == nil {
			dialog.ShowInformation("Error", "Invalid parent block", parentWindow)
			return
		}

		// For the new block, we need a miner's address.
		// Here, as an example, we use the parent's coinbase output address.
		// You could add another field to let the user choose the miner key.
		minerAddress := parentNode.B.GetCoinbase().Outputs[0].Address

		// Create a new block under the selected parent.
		newBlock := third_faza.NewBlock(parentNode.B.GetHash(), minerAddress)

		// Add all transactions that were selected.
		for _, chk := range txCheckboxes {
			if chk.Checked {
				tx := txMap[chk]
				newBlock.TransactionAdd(tx)
			}
		}

		// Finalize the block (e.g., compute its hash).
		newBlock.Finalizee()

		// Process the block using your BlockProcess function.
		if third_faza.BlockProcess(newBlock) {
			dialog.ShowInformation("Success", fmt.Sprintf("Block created with hash: %.6x", newBlock.GetHash()), parentWindow)
		} else {
			dialog.ShowInformation("Error", "Block processing failed", parentWindow)
		}
	}, parentWindow)
}

func buildMineBlockScreen() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("⛏️ Mine Block", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	minerOptions := make([]string, len(keyPairs))
	for i, kp := range keyPairs {
		minerOptions[i] = kp.Name
	}

	minerSelect := widget.NewSelect(minerOptions, func(s string) {})
	minerSelect.PlaceHolder = "Select Miner Key"

	methodOptions := []string{"BlockCreate", "BlockAdd"}
	methodSelect := widget.NewSelect(methodOptions, func(s string) {})
	methodSelect.PlaceHolder = "Select Mining Method"

	statusLabel := widget.NewLabel("")

	mineBtn := widget.NewButton("Mine Block", func() {
		if minerSelect.Selected == "" {
			statusLabel.SetText("Please select a Miner Key")
			return
		}
		if methodSelect.Selected == "" {
			statusLabel.SetText("Please select a Mining Method")
			return
		}

		var minerKP *KeyPair
		for i := range keyPairs {
			if keyPairs[i].Name == minerSelect.Selected {
				minerKP = &keyPairs[i]
				break
			}
		}
		if minerKP == nil {
			statusLabel.SetText("Miner Key not found")
			return
		}

		if methodSelect.Selected == "BlockCreate" {
			newBlock := third_faza.BlockCreate(minerKP.PublicKey)
			if newBlock != nil {
				statusLabel.SetText(fmt.Sprintf("Block created with hash: %.6x", newBlock.GetHash()))
			} else {
				statusLabel.SetText("BlockCreate failed.")
			}
		} else {
			showBlockProcessPopup(mainWindow)
		}

	})

	content := container.NewVBox(
		title,
		widget.NewLabel("Select Miner Key:"),
		minerSelect,
		widget.NewLabel("Select Mining Method:"),
		methodSelect,
		mineBtn,
		statusLabel,
	)
	return container.NewCenter(content)
}
