package main

import "math/rand"

const (
	width, height = 200, 100
	playerPadSize = 4
)

type Vector2D struct {
	dx float32
	dy float32
}

type Match struct {
	player1Pos   int
	player2Pos   int
	ballX        int
	ballY        int
	ballDir      Vector2D
	scorePlayer1 int
	scorePlayer2 int
}

func (m *Match) ToString() string {
	return string(m.player1Pos) + "," +
		string(m.player2Pos) + "," +
		string(m.ballX) + "," +
		string(m.ballY) + "," +
		string(m.scorePlayer1) + "," +
		string(m.scorePlayer2)
}

func (m *Match) StartMatch() {
	m.player1Pos = (height - playerPadSize) / 2
	m.player2Pos = (height - playerPadSize) / 2
	m.ballX = width / 2
	m.ballY = height / 2
}

func (m *Match) MovePlayerUp(player int) {
	if player == 1 {
		if m.player1Pos == 0 {
			return
		}
		m.player1Pos--
	} else {
		if m.player2Pos == 0 {
			return
		}
		m.player2Pos--
	}
}

func (m *Match) MovePlayerDown(player int) {
	if player == 1 {
		if m.player1Pos == (height - playerPadSize) {
			return
		}
		m.player1Pos++
	} else {
		if m.player2Pos == (height - playerPadSize) {
			return
		}
		m.player2Pos++
	}
}

func (m *Match) NextFrame() {
	// Ball collision with players
	if m.ballX == 1 && m.ballY >= m.player1Pos && m.ballY < m.player1Pos+playerPadSize {
		m.ballDir.dx = -m.ballDir.dx
	}
	if m.ballX == width-2 && m.ballY >= m.player2Pos && m.ballY < m.player2Pos+playerPadSize {
		m.ballDir.dx = -m.ballDir.dx
	}

	// Ball collision with walls
	if m.ballY == 0 || m.ballY == height {
		m.ballDir.dy = -m.ballDir.dy
	}

	// Ball collision with goals
	// Player 1 goal
	if m.ballX == width-1 {
		m.scorePlayer1++
		m.ballX = width / 2
		m.ballY = height / 2
		m.ballDir.dx = 1                    // direction to player 2
		m.ballDir.dy = rand.Float32()*2 - 1 // random number between -1 and 1
		return
	}
	// Player 2 goal
	if m.ballX == 0 {
		m.scorePlayer2++
		m.ballX = width / 2
		m.ballY = height / 2
		m.ballDir.dx = -1                   // direction to player 1
		m.ballDir.dy = rand.Float32()*2 - 1 // random number between -1 and 1
		return
	}

	m.ballX += int(m.ballDir.dx)
	m.ballY += int(m.ballDir.dy)

}
