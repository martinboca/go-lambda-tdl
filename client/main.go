package main

import (
	"bufio"
	"net"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	// Goroutine para manejar los mensajes del servidor
	go handleServer(conn)

	for {
		ev := <-eventQueue
		if ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyArrowUp {
				// Enviar acción de flecha hacia arriba al servidor
				conn.Write([]byte(buildPlayerAction(ActionMoveUp)))
				continue
			}
			if ev.Key == termbox.KeyArrowDown {
				conn.Write([]byte(buildPlayerAction(ActionMoveDown)))
				continue
			}
			if ev.Key == termbox.KeyEsc {
				return
			}
		}
	}
}

func handleServer(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if err != nil {
			// Manejar error
			break
		}
		if strings.HasPrefix(message, GameUpdateHeader) {
			handleGameUpdate(message)
			continue
		}
	}
}

func handleGameUpdate(message string) {
	// Ejemplo: GAME_UPDATE:player1Pos,player2Pos,ballX,ballY,scorePlayer1,scorePlayer2\n
	parts := strings.Split(message, ":")
	gameUpdate := strings.Split(parts[1], ",")
	if len(gameUpdate) == 6 {
		player1Pos, _ := strconv.Atoi(gameUpdate[0])
		player2Pos, _ := strconv.Atoi(gameUpdate[1])
		ballX, _ := strconv.Atoi(gameUpdate[2])
		ballY, _ := strconv.Atoi(gameUpdate[3])
		scorePlayer1, _ := strconv.Atoi(gameUpdate[4])
		scorePlayer2, _ := strconv.Atoi(gameUpdate[5])
		DrawPongInterface(player1Pos, player2Pos, ballX, ballY, scorePlayer1, scorePlayer2)
	}
}

func buildPlayerAction(action string) string {
	return PlayerActionHeader + EndOfHeader + action + EndOfMessage
}
