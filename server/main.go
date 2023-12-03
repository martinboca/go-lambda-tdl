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

type Player struct {
	name string
	conn net.Conn
}

type ServerMatch struct {
	player1 Player
	player2 Player
	match   Match
	started bool
	paused  bool
}

var (
	activeMatches  = make(map[int]*ServerMatch)
	connectedUsers = make(map[string]bool)
	matchIndex     = 0
	matchesMtx     sync.Mutex
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

func identifyClient(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	message = strings.TrimSpace(message)
	if err != nil {
		fmt.Println("Connection error: ", err)
		conn.Close()
		return "", err
	}
	//validar header?
	parts := strings.Split(message, ":")
	name := parts[1]
	return name, nil
}

func handleNewClient(conn net.Conn) {
	defer matchesMtx.Unlock()
	matchesMtx.Lock()

	name, err := identifyClient(conn)
	if err != nil {
		fmt.Println("Connection error identifying client: ", err)
		return
	}
	if _, ok := connectedUsers[name]; ok {
		fmt.Println("User", name, "already connected")
		//TODO avisar al cliente que ya esta conectado y cerrarle la conexion
		conn.Close()
		return
	} else {
		fmt.Println("User", name, "connected")
		connectedUsers[name] = true
	}

	//Buscar si el jugador ya esta en una partida
	for match_id, serverMatch := range activeMatches {
		fmt.Println("checking", match_id)
		if serverMatch.player1.name == name && serverMatch.paused {
			serverMatch.player1 = Player{conn: conn, name: name}
			serverMatch.paused = false
			fmt.Println("Player", name, "reconnected to the match", match_id)
			go handlePlayer(conn, 1, name, serverMatch)
			return
		}
		if serverMatch.player2.name == name && serverMatch.paused {
			serverMatch.player2 = Player{conn: conn, name: name}
			serverMatch.paused = false
			fmt.Println("Player", name, "reconnected to the match", match_id)
			go handlePlayer(conn, 2, name, serverMatch)
			return
		}
		fmt.Println("Player", name, "was not in a match")
		//sale del for si no esta en una partida
	}

	// Buscar una partida con un solo jugador
	for match_id, serverMatch := range activeMatches {
		if !serverMatch.started {
			// Agregar a este cliente a la partida existente y comenzar el juego
			serverMatch.player2 = Player{conn: conn, name: name}
			serverMatch.started = true
			go handleMatch(serverMatch, match_id)
			return
		}
	}

	// Si no hay partidas disponibles, crear una nueva
	newMatch := &ServerMatch{
		player1: Player{conn: conn, name: name},
		started: false,
		paused:  false,
	}
	activeMatches[matchIndex] = newMatch
	matchIndex++

}

func handleMatch(serverMatch *ServerMatch, match_id int) {
	fmt.Println("MATCH STARTED")
	//defer func() {
	//	matchesMtx.Lock()
	//	delete(activeMatches, match_id)
	//	matchesMtx.Unlock()
	//	fmt.Println("MATCH ENDED")
	//}() NO CREO QUE ESTE BIEN TERMINAR EL PARTIDO ACA, PORQUE SE HACEN LOS HANDLERS Y AL TOQUE SE BORRA EL ACTIVE MATCH...

	player1Conn := serverMatch.player1.conn
	player2Conn := serverMatch.player2.conn
	// Start game and handle updates
	serverMatch.match.StartMatch()
	go handlePlayer(player1Conn, 1, serverMatch.player1.name, serverMatch)
	go handlePlayer(player2Conn, 2, serverMatch.player2.name, serverMatch)
	go handleGameUpdates(serverMatch)
}

// Handles player messages
func handlePlayer(playerConn net.Conn, playerNumber int, name string, serverMatch *ServerMatch) {
	fmt.Println("MATCH STARTED")
	defer func() {
		playerConn.Close()
		serverMatch.paused = true
		matchesMtx.Lock()
		delete(connectedUsers, name)
		matchesMtx.Unlock()
	}()
	reader := bufio.NewReader(playerConn)
	for {
		message, err := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if err != nil {
			fmt.Println("Player", playerNumber, "connection error: ", err)
			break
		}
		if serverMatch.started && !serverMatch.paused {
			//actions
			if strings.HasPrefix(message, PlayerActionHeader) {
				handlePlayerAction(message, playerNumber, serverMatch)
				continue
			}
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
	// sendGameUpdate(serverMatch) esto MEPA que esta de mas
}

// Send game updates to players
func handleGameUpdates(serverMatch *ServerMatch) {
	// Send game updated every 10 ms
	go func() {
		for {
			if serverMatch.started && !serverMatch.paused {
				sendGameUpdate(serverMatch)
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	// Update game state every 1 second
	for {
		if serverMatch.started && !serverMatch.paused {
			waitTime := serverMatch.match.GetNextFrameWaitTime()
			time.Sleep(time.Duration(waitTime) * time.Millisecond)
			serverMatch.match.NextFrame()
		}
	}
}

func sendGameUpdate(serverMatch *ServerMatch) {
	match_string := serverMatch.match.ToString()
	serverMatch.player1.conn.Write([]byte(GameUpdateHeader + ":" + match_string + "\n"))
	serverMatch.player2.conn.Write([]byte(GameUpdateHeader + ":" + match_string + "\n"))
}
