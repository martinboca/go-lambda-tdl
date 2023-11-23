package main

import (
	"github.com/nsf/termbox-go"
)

const (
	width, height = 80, 20
	playerPadSize = 4
)

func DrawPongInterface(player1Pad, player2Pad, ballX, ballY int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Dibujar bordes horizontales
	for x := 0; x < width; x++ {
		termbox.SetCell(x, 0, '-', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, height, '-', termbox.ColorWhite, termbox.ColorDefault)
	}

	// Dibujar bordes verticales y barras
	for y := 0; y <= height; y++ {
		termbox.SetCell(0, y, '|', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(width-1, y, '|', termbox.ColorWhite, termbox.ColorDefault)

		// Dibujar barra del jugador 1
		if y >= player1Pad && y < player1Pad+playerPadSize {
			termbox.SetCell(1, y, '█', termbox.ColorWhite, termbox.ColorDefault)
		}

		// Dibujar barra del jugador 2
		if y >= player2Pad && y < player2Pad+playerPadSize {
			termbox.SetCell(width-2, y, '█', termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	// Dibujar la pelota
	termbox.SetCell(ballX, ballY, 'O', termbox.ColorRed, termbox.ColorDefault)

	termbox.Flush()
}
