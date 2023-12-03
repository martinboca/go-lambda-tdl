package main

import (
	termbox "github.com/nsf/termbox-go"
	"strconv"
)

const (
	width, height   = 90, 30
	playerPadSize   = 4
	scoreTextLength = 5 // Longitud máxima del texto del puntaje (ejemplo: "99-99")
)

func DrawPongInterface(player1Pad, player2Pad, ballX, ballY, scorePlayer1 int, scorePlayer2 int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw horizontal borders
	for x := 0; x < width; x++ {
		termbox.SetCell(x, 0, '¯', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, height, '_', termbox.ColorWhite, termbox.ColorDefault)
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

	// Calculate the center position for the score text
	centerX := width / 2
	centerY := height / 2

	// Draw player 1 score
	score1Text := strconv.Itoa(scorePlayer1)
	score1X := centerX - len(score1Text) - 2 // Center left align
	DrawText(score1X, centerY, score1Text)

	DrawText(centerX-1, centerY, "-")

	// Draw player 2 score
	score2Text := strconv.Itoa(scorePlayer2)
	score2X := centerX + 1 // Center right align
	DrawText(score2X, centerY, score2Text)

	termbox.Flush()
}

func DrawText(x, y int, text string) {
	for i, c := range text {
		termbox.SetCell(x+i, y, c, termbox.ColorWhite, termbox.ColorDefault)
	}
}

func DrawWelcomeScreen() {
	DrawText(width/2-10, height/2, "Waiting for player...")
	termbox.Flush()
}
