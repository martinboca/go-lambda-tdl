package main

/* Documentation:
The messages sent by the client and server are in the format: HEADER:DATA\n

There are different types of messages, each one with its own header and data format.
GAME_UPDATE:
	Full message: "GAME_UPDATE:player1Pos,player2Pos,ballPosX,ballPosY,scorePlayer1,scorePlayer2\n"
	Data format: "player1Pos,player2Pos,ballPosX,ballPosY\n"
		- player1Pos: int (y position in pixels from the top of the screen)
		- player2Pos: int (y position in pixels from the top of the screen)
		- ballPosX: int (x position in pixels from the left of the screen)
		- ballPosY: int (y position in pixels from the top of the screen)
		- scorePlayer1: int (player 1 score)
		- scorePlayer2: int (player 2 score)

CLIENT_ACTION:
	Full message: "CLIENT_ACTION:action\n"
	Data format: "action\n"
		- action: string (action to be performed by the server. ie: "MOVE_UP", "MOVE_DOWN")

*/

// Special characters
const (
	EndOfHeader  = ":"
	EndOfMessage = "\n"
)

// Client headers
const (
	PlayerActionHeader     = "CLIENT_ACTION"
	ClientConnectionHeader = "CLIENT_CONNECTION"
)

// Server headers
const (
	GameUpdateHeader         = "GAME_UPDATE"
	GameStartHeader          = "GAME_START"
	GameEndHeader            = "GAME_END"
	PlayerDisconnectedHeader = "PLAYER_DISCONNECTED"
)

const (
	targetScore = 3
)

// Client actions
const (
	ActionMoveUp    = "MOVE_UP"
	ActionMoveDown  = "MOVE_DOWN"
	ActionExplosion = "EXPLOSION"
)
