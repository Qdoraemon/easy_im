package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// "net/http"
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var urgrader = websocket.Upgrader{}

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

// 用来升级HTTP连接为WebSocket连接
func handleConnents(w http.ResponseWriter, r *http.Request) {
	ws, err := urgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}

}

func main() {

}
