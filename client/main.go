package main

import (
	"bufio"
	"fmt"
	"github.com/nsf/termbox-go"
	"net"
	"os"
	"os/exec"
	"runtime"
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

	name := identifyClient(conn)
	fmt.Println("Hello: ", name)

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	// Goroutine para manejar los mensajes del servidor
	go handleServer(conn)
	clearScreen()
	DrawWelcomeScreen()

	// Esperar hasta que el otro cliente se conecte
	_, err = conn.Read(make([]byte, 1))
	if err != nil {
		fmt.Println("Error waiting for another client connection:", err)
		return
	}

	// Borrar el mensaje de espera
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()

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

func lineBreakRune() byte {
	if os.PathSeparator == '\\' {
		return '\r'
	} else {
		return '\n'
	}
}

func identifyClient(conn net.Conn) string {
	fmt.Println("Enter your ID/name:")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString(lineBreakRune())
	if err != nil || len(name) == 0 {
		name = "GUEST"
	}
	conn.Write([]byte(ClientConnectionHeader + EndOfHeader + name + EndOfMessage))
	return name
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

func clearScreen() {
	var cmd *exec.Cmd

	// Verificar el sistema operativo y usar el comando correspondiente
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	// Ejecutar el comando para limpiar la pantalla
	cmd.Stdout = os.Stdout
	cmd.Run()
}
