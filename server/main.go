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
	defer conn.Close()
	matchesMtx.Lock()
	defer matchesMtx.Unlock()

	// Buscar una partida con un solo jugador
	for _, match := range activeMatches {
		if !match.started {
			// Agregar a este cliente a la partida existente y comenzar el juego
			match.player2 = conn
			match.started = true
			go handleMatch(match)
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
	defer func() {
		matchesMtx.Lock()
		delete(activeMatches, serverMatch.player1)
		delete(activeMatches, serverMatch.player2)
		matchesMtx.Unlock()
	}()

	serverMatch.match.StartMatch()
	match := serverMatch.match
	player1Conn := serverMatch.player1
	player2Conn := serverMatch.player2

	// Enviar actualizaciones a los jugadores

	// Handle player player actions
	go handlePlayer(player1Conn, 1, &match)
	go handlePlayer(player2Conn, 2, &match)

	go handleGameUpdates(&match, serverMatch)
}

// Handles player messages
func handlePlayer(playerConn net.Conn, playerNumber int, match *Match) {
	reader := bufio.NewReader(playerConn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			// Manejar error
			break
		}
		if strings.HasPrefix(message, PlayerActionHeader) {
			handlePlayerAction(message, playerNumber, match)
			continue
		}
	}
}

// Update game state based on player actions
func handlePlayerAction(message string, playerNumber int, match *Match) {
	// Ejemplo: CLIENT_ACTION:MOVE_UP\n
	parts := strings.Split(message, ":")
	action := parts[1]
	if action == ActionMoveUp {
		match.MovePlayerUp(playerNumber)
	} else if action == ActionMoveDown {
		match.MovePlayerDown(playerNumber)
	}
}

// Send game updates to players
func handleGameUpdates(match *Match, serverMatch *ServerMatch) {
	for {
		serverMatch.player1.Write([]byte(GameUpdateHeader + ":" + match.ToString() + "\n"))
		serverMatch.player2.Write([]byte(GameUpdateHeader + ":" + match.ToString() + "\n"))
		time.Sleep(100 * time.Millisecond)
	}
}
