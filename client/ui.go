package main

import (
	"github.com/nsf/termbox-go"
)

const (
	width, height = 90, 30
	playerPadSize = 4
)

func DrawPongInterface(player1Pad, player2Pad, ballX, ballY, scorePlayer1, scorePlayer2 int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw horizontal borders
	for x := 0; x < width; x++ {
		termbox.SetCell(x, 0, '-', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, height, '-', termbox.ColorWhite, termbox.ColorDefault)
	}

	// Draw vertical borders and player bars
	for y := 0; y <= height; y++ {
		termbox.SetCell(0, y, '|', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(width-1, y, '|', termbox.ColorWhite, termbox.ColorDefault)

		// Draw player 1 bar
		if y >= player1Pad && y < player1Pad+playerPadSize {
			termbox.SetCell(1, y, '█', termbox.ColorWhite, termbox.ColorDefault)
		}

		// Draw player 2 bar
		if y >= player2Pad && y < player2Pad+playerPadSize {
			termbox.SetCell(width-2, y, '█', termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	// Draw ball
	termbox.SetCell(ballX, ballY, 'O', termbox.ColorRed, termbox.ColorDefault)
	termbox.Flush()
}

func DrawText(x, y int, text string) {
	for i, c := range text {
		termbox.SetCell(x+i, y, c, termbox.ColorWhite, termbox.ColorDefault)
	}
	termbox.Flush()
}

func DrawWelcomeScreen() {
	DrawText(width/2-10, height/2, "Waiting for players...")
}
