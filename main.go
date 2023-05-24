package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

type Room struct {
	clients  map[*websocket.Conn]bool
	clientsMu sync.Mutex
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow all connections for simplicity
			return true
		},
	}

	rooms     = make(map[string]*Room)
	roomsMu   sync.Mutex
)

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

func main() {
	r := mux.NewRouter()

	corsMiddleware := cors.Default().Handler

	r.HandleFunc("/ws", websocketHandler)
	r.HandleFunc("/rooms", roomsHandler).Methods("GET")

	log.Println("Starting WebSocket server on http://localhost:8080")
	err := http.ListenAndServe(":8080", corsMiddleware(r))
	if err != nil {
		log.Fatal("WebSocket server failed:", err)
	}
}
