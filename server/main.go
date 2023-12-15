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
	name      string
	conn      net.Conn
	connected bool
}

type ServerMatch struct {
	player1   Player
	player2   Player
	match     Match
	preparing bool
	started   bool
	paused    bool
}

func (s *ServerMatch) disconnectPlayer(number int) {
	if number == 1 {
		s.player1.connected = false
	} else {
		s.player2.connected = false
	}
}

func (s *ServerMatch) allPlayersConnected() bool {
	return s.player1.connected && s.player2.connected
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
		if serverMatch.paused {
			fmt.Println("checking paused match", match_id)
			if serverMatch.player1.name == name && serverMatch.paused {
				serverMatch.player1.conn = conn
				serverMatch.player1.connected = true
				if serverMatch.allPlayersConnected() {
					serverMatch.paused = false
					serverMatch.preparing = true
				}
				fmt.Println("Player", name, "reconnected to the match", match_id)
				go handlePlayer(conn, 1, name, serverMatch)
				return
			}
			if serverMatch.player2.name == name && serverMatch.paused {
				serverMatch.player2.conn = conn
				serverMatch.player2.connected = true
				if serverMatch.allPlayersConnected() {
					serverMatch.paused = false
					serverMatch.preparing = true
				}
				fmt.Println("Player", name, "reconnected to the match", match_id)
				go handlePlayer(conn, 2, name, serverMatch)
				return
			}
		}
	}
	fmt.Println("Player", name, "was not in a match")

	// Buscar una partida con un solo jugador
	for match_id, serverMatch := range activeMatches {
		if !serverMatch.started {
			// Agregar a este cliente a la partida existente y comenzar el juego
			serverMatch.player2 = Player{conn: conn, name: name, connected: true}
			serverMatch.started = true
			go handleMatch(serverMatch, match_id)
			return
		}
	}

	// Si no hay partidas disponibles, crear una nueva
	newMatch := &ServerMatch{
		player1:   Player{conn: conn, name: name, connected: true},
		preparing: true,
		started:   false,
		paused:    false,
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
	fmt.Println("PLAYER", playerNumber, "HANDLER STARTED", name)
	defer func() {
		playerConn.Close()
		serverMatch.disconnectPlayer(playerNumber)
		serverMatch.paused = true
		matchesMtx.Lock()
		delete(connectedUsers, name)
		matchesMtx.Unlock()
		fmt.Println("PLAYER", playerNumber, "HANDLER ENDED", name)
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
	parts := strings.Split(message, EndOfHeader)
	action := parts[1]

	if action == ActionMoveUp {
		serverMatch.match.MovePlayerUp(playerNumber)
	} else if action == ActionMoveDown {
		serverMatch.match.MovePlayerDown(playerNumber)
	} else if action == ActionExplosion {
		serverMatch.match.Explosion()
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
			if !serverMatch.allPlayersConnected() {
				sendDisconnectedMessage(serverMatch)
				time.Sleep(1 * time.Second)
			}
		}
	}()

	// Update game state every 1 second
	for {
		if serverMatch.started && !serverMatch.paused && !serverMatch.preparing {
			waitTime := serverMatch.match.GetNextFrameWaitTime()
			time.Sleep(time.Duration(waitTime) * time.Millisecond)
			serverMatch.match.NextFrame()
		}
	}
}

func sendDisconnectedMessage(match *ServerMatch) {
	if !match.player1.connected && match.player2.connected {
		match.player2.conn.Write([]byte(PlayerDisconnectedHeader + EndOfHeader + match.player1.name + EndOfMessage))
	}
	if !match.player2.connected && match.player1.connected {
		match.player1.conn.Write([]byte(PlayerDisconnectedHeader + EndOfHeader + match.player2.name + EndOfMessage))
	}
}

func sendGameUpdate(serverMatch *ServerMatch) {
	if serverMatch.preparing {
		serverMatch.player1.conn.Write([]byte(GameStartHeader + EndOfHeader + serverMatch.player2.name + EndOfMessage))
		serverMatch.player2.conn.Write([]byte(GameStartHeader + EndOfHeader + serverMatch.player1.name + EndOfMessage))

		time.Sleep(3 * time.Second)
		serverMatch.preparing = false
		return
	}

	if serverMatch.match.scorePlayer1 == targetScore || serverMatch.match.scorePlayer2 == targetScore {
		winner := func() string {
			if serverMatch.match.scorePlayer1 == targetScore {
				return serverMatch.player1.name
			}
			return serverMatch.player2.name
		}()

		serverMatch.player1.conn.Write([]byte(GameEndHeader + EndOfHeader + winner + EndOfMessage))
		serverMatch.player2.conn.Write([]byte(GameEndHeader + EndOfHeader + winner + EndOfMessage))

		time.Sleep(2 * time.Second)
		serverMatch.paused = true
		return
	}

	match_string := serverMatch.match.ToString()

	serverMatch.player1.conn.Write([]byte(GameUpdateHeader + EndOfHeader + match_string + EndOfMessage))
	serverMatch.player2.conn.Write([]byte(GameUpdateHeader + EndOfHeader + match_string + EndOfMessage))
}
