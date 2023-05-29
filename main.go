package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

type PieceType int

const (
	Empty PieceType = iota
	Pawn
	Rook
	Knight
	Bishop
	Queen
	King
)

type Room struct {
	clients   map[*websocket.Conn]bool
	clientsMu sync.Mutex
}

type MoveRequest struct {
	From Coordinates `json:"from"`
	To   Coordinates `json:"to"`
}

type Coordinates struct {
	X int
	Y int
}

type Move struct {
	From Coordinates
	To   Coordinates
}

type Piece struct {
	pieceType  PieceType
	color      string
	legalMoves []Move
}

type Tile struct {
	piece Piece
}

type Board struct {
	pos [8][8]Tile
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow all connections for simplicity
			return true
		},
	}

	rooms    = make(map[string]*Room)
	roomsMu  sync.Mutex
	cpuGames = make(map[string]Board)
)

func initBoard(board *Board) {
	// Initialize the pawns for both white and black
	for file := 0; file < 8; file++ {
		board.pos[1][file] = Tile{Piece{color: "black", pieceType: Pawn, legalMoves: []Move{}}}
		board.pos[6][file] = Tile{Piece{color: "white", pieceType: Pawn, legalMoves: []Move{}}}
	}

	// Initialize the rooks for both white and black
	board.pos[0][0] = Tile{Piece{color: "black", pieceType: Rook, legalMoves: []Move{}}}
	board.pos[0][7] = Tile{Piece{color: "black", pieceType: Rook, legalMoves: []Move{}}}
	board.pos[7][0] = Tile{Piece{color: "white", pieceType: Rook, legalMoves: []Move{}}}
	board.pos[7][7] = Tile{Piece{color: "white", pieceType: Rook, legalMoves: []Move{}}}

	// Initialize the knights for both white and black
	board.pos[0][1] = Tile{Piece{color: "black", pieceType: Knight, legalMoves: []Move{}}}
	board.pos[0][6] = Tile{Piece{color: "black", pieceType: Knight, legalMoves: []Move{}}}
	board.pos[7][1] = Tile{Piece{color: "white", pieceType: Knight, legalMoves: []Move{}}}
	board.pos[7][6] = Tile{Piece{color: "white", pieceType: Knight, legalMoves: []Move{}}}

	// Initialize the bishops for both white and black
	board.pos[0][2] = Tile{Piece{color: "black", pieceType: Bishop, legalMoves: []Move{}}}
	board.pos[0][5] = Tile{Piece{color: "black", pieceType: Bishop, legalMoves: []Move{}}}
	board.pos[7][2] = Tile{Piece{color: "white", pieceType: Bishop, legalMoves: []Move{}}}
	board.pos[7][5] = Tile{Piece{color: "white", pieceType: Bishop, legalMoves: []Move{}}}

	// Initialize the queens for both white and black
	board.pos[0][3] = Tile{Piece{color: "black", pieceType: Queen, legalMoves: []Move{}}}
	board.pos[7][3] = Tile{Piece{color: "white", pieceType: Queen, legalMoves: []Move{}}}

	// Initialize the kings for both white and black
	board.pos[0][4] = Tile{Piece{color: "black", pieceType: King, legalMoves: []Move{}}}
	board.pos[7][4] = Tile{Piece{color: "white", pieceType: King, legalMoves: []Move{}}}
}

func printBoard(board *Board) {
	fmt.Println("  a b c d e f g h")
	for i := 0; i < len(board.pos); i++ {
		fmt.Printf("%d ", i+1)
		for j := 0; j < len(board.pos[i]); j++ {
			piece := board.pos[i][j].piece
			pieceSymbol := getPieceSymbol(piece)
			fmt.Printf("%s ", pieceSymbol)
		}
		fmt.Println()
	}
}

func getPieceSymbol(piece Piece) string {
	switch piece.pieceType {
	case Empty:
		return "-"
	case Pawn:
		return "P"
	case Rook:
		return "R"
	case Knight:
		return "N"
	case Bishop:
		return "B"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		return "?"
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")

	roomsMu.Lock()
	room, ok := rooms[roomID]
	if !ok {
		room = &Room{
			clients: make(map[*websocket.Conn]bool),
		}
		rooms[roomID] = room
	}
	roomsMu.Unlock()

	room.clientsMu.Lock()
	numClients := len(room.clients)
	room.clientsMu.Unlock()

	if numClients >= 2 {
		log.Println("Only two players per room")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	room.clientsMu.Lock()
	room.clients[conn] = true
	room.clientsMu.Unlock()

	defer func() {
		room.clientsMu.Lock()
		delete(room.clients, conn)
		room.clientsMu.Unlock()

		conn.Close()
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading WebSocket message:", err)
			break
		}

		log.Printf("Received message in room %s: %s\n", roomID, message)

		broadcast(room, messageType, message)
	}
}

func broadcast(room *Room, messageType int, message []byte) {
	room.clientsMu.Lock()
	defer room.clientsMu.Unlock()

	for client := range room.clients {
		err := client.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error sending WebSocket message:", err)
		}
	}
}

func roomsHandler(w http.ResponseWriter, r *http.Request) {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	var availableRooms []string

	for roomID, room := range rooms {
		room.clientsMu.Lock()
		numClients := len(room.clients)
		room.clientsMu.Unlock()

		if numClients == 1 {
			availableRooms = append(availableRooms, roomID)
		}
	}

	json.NewEncoder(w).Encode(availableRooms)
}

func initComputerGame(w http.ResponseWriter, r *http.Request) {
	id := uuid.New()
	idString := id.String()
	fmt.Println(idString)
	board := Board{}
	initBoard(&board)
	cpuGames[idString] = board

	json.NewEncoder(w).Encode(idString)
}

func playComputerGame(w http.ResponseWriter, r *http.Request) {
	computerGame := r.URL.Query().Get("game")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Close the request body
	defer r.Body.Close()

	// Parse the JSON data
	var moveReq MoveRequest
	err = json.Unmarshal(body, &moveReq)
	if err != nil {
		log.Println("Error parsing JSON:", err)
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Access the move details
	from := moveReq.From
	to := moveReq.To
	move := Move{From: from, To: to}
	board := cpuGames[computerGame]

	calculateLegalMoves(&board)
	//printBoard(&board)

	movePiece(&move, &board)

	// Close the request body
	defer r.Body.Close()

	json.NewEncoder(w).Encode(getRandomLegalMove(&board))
}

func main() {
	r := mux.NewRouter()

	corsMiddleware := cors.Default().Handler

	r.HandleFunc("/ws", websocketHandler)
	r.HandleFunc("/rooms", roomsHandler).Methods("GET")
	r.HandleFunc("/initComputerGame", initComputerGame).Methods("POST")
	r.HandleFunc("/playComputerGame", playComputerGame).Methods("POST")

	log.Println("Starting WebSocket server on http://localhost:8080")
	err := http.ListenAndServe(":8080", corsMiddleware(r))
	if err != nil {
		log.Fatal("WebSocket server failed:", err)
	}
}
