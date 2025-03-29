package main

import (
	"DMBLOCK_GO/third_faza"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

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
		bubble := styledTxBubble(tx)
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
