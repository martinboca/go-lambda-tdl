package main

import (
	"bufio"
	"fmt"
	"github.com/nsf/termbox-go"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	var conn net.Conn
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	for { //Loop until connected
		conn, err = net.Dial("tcp", "localhost:8080")
		if err != nil {
			fmt.Println(err, " | Connection failed, retrying...")
			time.Sleep(5 * time.Second)
			continue
		}
		break // Successful connection
	}
	defer conn.Close()

	identifyClient(conn)

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	// Goroutine para manejar los mensajes del servidor
	go handleServer(conn)
	showMessage(width/2-10, height/2, "Waiting for player...")

	// Esperar hasta que el otro cliente se conecte
	_, err = conn.Read(make([]byte, 1))
	if err != nil {
		fmt.Println("Error waiting for another client connection:", err)
		return
	}

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
			if ev.Key == termbox.KeySpace {
				conn.Write([]byte(buildPlayerAction(ActionExplosion)))
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
			fmt.Println("Connection error: ", err)
			break
		}
		if strings.HasPrefix(message, GameUpdateHeader) || strings.HasPrefix(message, GameEndHeader) {
			handleGameUpdate(message)
			continue
		}
	}
}

func handleGameUpdate(message string) {
	parts := strings.Split(message, ":")

	// Ejemplo: GAME_END:player1Name\n represents the winner
	if parts[0] == GameEndHeader {
		if len(parts[1]) == 1 {
			showMessage(width/2-10, height/2, fmt.Sprintf("Juego Finalizado%v¡El ganador es: %s!", "\n", parts[1]))
			return
		}
	}

	// Ejemplo: GAME_UPDATE:player1Pos,player2Pos,ballX,ballY,scorePlayer1,scorePlayer2\n
	if parts[0] == GameUpdateHeader {
		gameUpdate := strings.Split(parts[1], ",")
		if len(gameUpdate) == 7 {
			player1Pos, _ := strconv.Atoi(gameUpdate[0])
			player2Pos, _ := strconv.Atoi(gameUpdate[1])
			ballX, _ := strconv.Atoi(gameUpdate[2])
			ballY, _ := strconv.Atoi(gameUpdate[3])
			scorePlayer1, _ := strconv.Atoi(gameUpdate[4])
			scorePlayer2, _ := strconv.Atoi(gameUpdate[5])
			explosion, _ := strconv.Atoi(gameUpdate[6])
			DrawPongInterface(player1Pos, player2Pos, ballX, ballY, scorePlayer1, scorePlayer2, explosion)
		}
	}
}

func buildPlayerAction(action string) string {
	return PlayerActionHeader + EndOfHeader + action + EndOfMessage
}

func identifyClient(conn net.Conn) string {
	showMessage(1, 1, "Enter your ID/name:")

	name := readInput()

	if len(name) == 0 {
		name = "GUEST"
	}

	conn.Write([]byte(ClientConnectionHeader + EndOfHeader + name + EndOfMessage))
	return name
}

func readInput() string {
	input := ""

inputLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				break inputLoop
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if len(input) > 0 {
					input = input[:len(input)-1]
				}
			default:
				input += string(ev.Ch)
			}
		}
		showMessage(1, 1, "Enter your ID/name: "+input)
	}

	return strings.TrimSpace(input)
}
