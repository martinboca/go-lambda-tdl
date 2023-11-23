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
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				// Manejar error
				break
			}

			if strings.HasPrefix(message, SRV_GAME_UPDATE) {
				// Parsear y actualizar el estado del juego
				// Ejemplo: GAME_UPDATE:player1Pos,player2Pos,ballX,ballY
				parts := strings.Split(message, ":")
				gameUpdate := strings.Split(parts[1], ",")
				if len(gameUpdate) == 4 {
					player1Pos, _ := strconv.Atoi(gameUpdate[0])
					player2Pos, _ := strconv.Atoi(gameUpdate[1])
					ballX, _ := strconv.Atoi(gameUpdate[2])
					ballY, _ := strconv.Atoi(gameUpdate[3])
					DrawPongInterface(player1Pos, player2Pos, ballX, ballY)
				}

			}
		}
	}()

	for {
		ev := <-eventQueue
		if ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyArrowUp {
				// Enviar acción de flecha hacia arriba al servidor
				conn.Write([]byte("CLIENT_ACTION:ARROW_UP\n"))
			}
			if ev.Key == termbox.KeyEsc {
				return
			}

		}
	}
}
