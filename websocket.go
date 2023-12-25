package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Opgrader GET-anmodningen til en WebSocket-forbindelse
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Registrer ny klient
	clients[ws] = true
	fmt.Println("Ny klient forbundet")

	for {
		// Read message from the WebSocket connection
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			delete(clients, ws)
			break
		}
		log.Printf("Received message: %s\n", string(msg))

		// Send message to all connected clients
		for client := range clients {
			err := client.WriteMessage(messageType, msg)
			if err != nil {
				log.Println("Error writing message:", err)
				delete(clients, client)
				break
			}
		}
	}
}

func RunWebSocket() {
	http.HandleFunc("/ws", handleConnections)
	log.Println("WebSocket-server kører på http://localhost:8080/ws")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
