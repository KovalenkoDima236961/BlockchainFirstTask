package main

import (
	"DMBLOCK_GO/third_faza"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

func showTxDialog(tx *third_faza.Transaction, parent fyne.Window) {
	inputsContainer := container.NewVBox()
	inputs := tx.GetInputs()
	if len(inputs) == 0 {
		inputsContainer.Add(widget.NewLabel("No inputs available."))
	} else {
		for i, input := range inputs {
			// Show only the first 6 characters of the previous transaction hash.
			hash := input.PrevTxHash
			if len(hash) > 6 {
				hash = hash[:6]
			}
			inputsContainer.Add(widget.NewLabel(fmt.Sprintf("Input %d: Hash: %s, Index: %d", i, hash[:6], input.OutputIndex)))
		}
	}

	outputsContainer := container.NewVBox()
	outputs := tx.GetOutputs()
	if len(outputs) == 0 {
		outputsContainer.Add(widget.NewLabel("No outputs available."))
	} else {
		for i, output := range outputs {
			// Format the address by showing only the first 6 characters.
			addr := fmt.Sprintf("%v", output.Address)
			if len(addr) > 6 {
				addr = addr[:6]
			}
			outputsContainer.Add(widget.NewLabel(fmt.Sprintf("Output %d: Value: %.2f", i, output.Value)))
		}
	}

	content := container.NewVBox(
		widget.NewLabelWithStyle("Transaction Details", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Inputs:"),
		inputsContainer,
		widget.NewLabel("Outputs:"),
		outputsContainer,
	)

	dialog.ShowCustom("Transaction Details", "Close", content, parent)
}

func clickableTxBubble(tx *third_faza.Transaction, parent fyne.Window) fyne.CanvasObject {
	bubble := styledTxBubble(tx)

	btn := widget.NewButton("", func() {
		showTxDialog(tx, parent)
	})
	btn.Importance = widget.LowImportance

	return container.NewMax(bubble, btn)
}

func styledTxBubble(tx *third_faza.Transaction) fyne.CanvasObject {
	shortHash := fmt.Sprintf("%.6x", tx.GetHash())

	// Create a circle
	circle := canvas.NewCircle(color.NRGBA{R: 80, G: 80, B: 150, A: 255})
	circle.Resize(fyne.Size{Width: 200, Height: 200}) // Make it a circle

	// Text content (shortened, because circle is small)
	hashText := canvas.NewText(shortHash, color.White)
	hashText.Alignment = fyne.TextAlignCenter
	hashText.TextSize = 12

	// Center the text on top of the circle
	content := container.NewCenter(hashText)

	return container.NewMax(circle, content)
}

func buildTransactionPoolView() fyne.CanvasObject {
	txPool := blockchain.GetTransactionPool()
	if txPool == nil || len(txPool.H) == 0 {
		return widget.NewLabel("No transactions in the pool.")
	}

	var txBubbles []fyne.CanvasObject
	for _, tx := range txPool.GetTransactions() {
		bubble := clickableTxBubble(tx, mainWindow)
		txBubbles = append(txBubbles, bubble)
	}

	vbox := container.NewVBox()
	for i, b := range txBubbles {
		if i > 0 {
			vbox.Add(layout.NewSpacer())
		}
		vbox.Add(b)
	}

	scroll := container.NewScroll(vbox)
	scroll.SetMinSize(fyne.NewSize(400, 400))

	return container.NewVBox(
		widget.NewLabelWithStyle("📋 Transaction Pool", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		scroll,
	)
}
