package main

import (
	"fmt"
	"math/rand"
)

const (
	width, height  = 90, 30
	playerPadSize  = 4
	initialSpeed   = 3
	speedIncrement = 0.01
)

type BallDirVector struct {
	dx int
	dy int
}

type Match struct {
	player1Pos   int
	player2Pos   int
	ballX        int
	ballY        int
	ballDir      BallDirVector
	scorePlayer1 int
	scorePlayer2 int
	speed        float32
	explosion    int
}

func (m *Match) ToString() string {
	return fmt.Sprint(int(m.player1Pos)) + "," +
		fmt.Sprint(int(m.player2Pos)) + "," +
		fmt.Sprint(int(m.ballX)) + "," +
		fmt.Sprint(int(m.ballY)) + "," +
		fmt.Sprint(int(m.scorePlayer1)) + "," +
		fmt.Sprint(int(m.scorePlayer2)) + "," +
		fmt.Sprint(int(m.explosion))
}

func (m *Match) StartMatch() {
	m.player1Pos = (height - playerPadSize) / 2
	m.player2Pos = (height - playerPadSize) / 2
	m.ballX = width / 2
	m.ballY = height / 2
	m.ballDir.dx = randomHorizontalDirection()
	m.ballDir.dy = randomVerticalDirection()
	m.speed = initialSpeed
	m.explosion = 0
}

func (m *Match) MovePlayerUp(player int) {
	if player == 1 {
		if m.player1Pos == 1 {
			return
		}
		m.player1Pos--
	} else {
		if m.player2Pos == 1 {
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

func (m *Match) GetNextFrameWaitTime() int {
	return int(500 / m.speed)
}

func (m *Match) NextFrame() {
	if m.explosion != 0 {
		m.explosion = 0
	}

	// Ball collision with players
	if m.ballX == 2 && m.ballY >= m.player1Pos && m.ballY < m.player1Pos+playerPadSize {
		m.ballDir.dx = -m.ballDir.dx
		m.ballDir.dy = randomVerticalDirection()
	}
	if m.ballX == width-3 && m.ballY >= m.player2Pos && m.ballY < m.player2Pos+playerPadSize {
		m.ballDir.dx = -m.ballDir.dx
		m.ballDir.dy = randomVerticalDirection()
	}

	// TODO: There is a bug when colliding with the player pad at the border of the screen
	//			 when the randomVerticalDirection points to the wall

	// Ball collision with walls
	if m.ballY == 1 || m.ballY == height-1 {
		m.ballDir.dy = -m.ballDir.dy
	}

	// Ball collision with goals
	// Player 1 scores
	if m.ballX == width-1 {
		goal(1, m)
	}
	// Player 2 scores
	if m.ballX == 0 {
		goal(2, m)
	}

	m.ballX += m.ballDir.dx
	m.ballY += m.ballDir.dy
	m.speed += speedIncrement
}

func goal(player int, m *Match) {
	if player == 1 {
		m.scorePlayer1++
		m.ballDir.dx = 1 // direction to player 2
	} else if player == 2 {
		m.scorePlayer2++
		m.ballDir.dx = -1 // direction to player 1
	}
	m.ballX = width / 2
	m.ballY = height / 2
	m.ballDir.dy = randomVerticalDirection()
	m.speed = 1
}

// Random -1 or 1
func randomHorizontalDirection() int {
	return rand.Intn(2)*2 - 1
}

// Random -1, 0 or 1
func randomVerticalDirection() int {
	return rand.Intn(3) - 1
}

func (m *Match) Explosion() {
	if m.ballX <= 3 && m.ballY >= m.player1Pos && m.ballY < m.player1Pos+playerPadSize {
		m.speed += speedIncrement * 100
		m.explosion = 1
	}
	if m.ballX >= width-4 && m.ballY >= m.player2Pos && m.ballY < m.player2Pos+playerPadSize {
		m.speed += speedIncrement * 100
		m.explosion = 2
	}
}
