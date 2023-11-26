package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type ServerMatch struct {
	player1 net.Conn
	player2 net.Conn
	match   Match
	started bool
}

var (
	activeMatches = make(map[net.Conn]*ServerMatch)
	matchesMtx    sync.Mutex
)

func main() {
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("New client: " + conn.RemoteAddr().String())
		go handleNewClient(conn)
	}
}

func handleNewClient(conn net.Conn) {
	defer matchesMtx.Unlock()
	matchesMtx.Lock()

	// Buscar una partida con un solo jugador
	for _, serverMatch := range activeMatches {
		if !serverMatch.started {
			// Agregar a este cliente a la partida existente y comenzar el juego
			serverMatch.player2 = conn
			serverMatch.started = true
			go handleMatch(serverMatch)
			return
		}
	}

	// Si no hay partidas disponibles, crear una nueva
	newMatch := &ServerMatch{
		player1: conn,
		started: false,
	}
	activeMatches[conn] = newMatch

}

func handleMatch(serverMatch *ServerMatch) {
	fmt.Println("MATCH STARTED")
	defer func() {
		matchesMtx.Lock()
		delete(activeMatches, serverMatch.player1)
		delete(activeMatches, serverMatch.player2)
		matchesMtx.Unlock()
	}()

	player1Conn := serverMatch.player1
	player2Conn := serverMatch.player2
	// Start game and handle updates
	serverMatch.match.StartMatch()
	go handlePlayer(player1Conn, 1, serverMatch)
	go handlePlayer(player2Conn, 2, serverMatch)
	go handleGameUpdates(serverMatch)
}

// Handles player messages
func handlePlayer(playerConn net.Conn, playerNumber int, serverMatch *ServerMatch) {
	defer playerConn.Close()
	reader := bufio.NewReader(playerConn)
	for {
		message, err := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if err != nil {
			// Manejar error
			break
		}
		if strings.HasPrefix(message, PlayerActionHeader) {
			handlePlayerAction(message, playerNumber, serverMatch)
			continue
		}
	}
}

// Update game state based on player actions
func handlePlayerAction(message string, playerNumber int, serverMatch *ServerMatch) {
	// Ejemplo: CLIENT_ACTION:MOVE_UP\n
	parts := strings.Split(message, ":")
	action := parts[1]

	if action == ActionMoveUp {
		serverMatch.match.MovePlayerUp(playerNumber)
	} else if action == ActionMoveDown {
		serverMatch.match.MovePlayerDown(playerNumber)
	}
	sendGameUpdate(serverMatch)
}

// Send game updates to players
func handleGameUpdates(serverMatch *ServerMatch) {
	// Send game updated every 10 ms
	go func() {
		for {
			sendGameUpdate(serverMatch)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Update game state every 1 second
	for {
		waitTime := serverMatch.match.GetNextFrameWaitTime()
		time.Sleep(time.Duration(waitTime) * time.Millisecond)
		serverMatch.match.NextFrame()
	}
}

func sendGameUpdate(serverMatch *ServerMatch) {
	match_string := serverMatch.match.ToString()
	serverMatch.player1.Write([]byte(GameUpdateHeader + ":" + match_string + "\n"))
	serverMatch.player2.Write([]byte(GameUpdateHeader + ":" + match_string + "\n"))
}
