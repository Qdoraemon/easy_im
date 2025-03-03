package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// "net/http"
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var urgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

// 用来升级HTTP连接为WebSocket连接
func handleConnents(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-------------")
	ws, err := urgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
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

// 监听广播通道，并将消息发送给所有连接的客户端
func handleMessages() {
	for {
		msg := <-broadcast
		fmt.Println("msg: ", msg)
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// 注册链接
	http.HandleFunc("/ws", handleConnents)
	go handleMessages()
	fmt.Println("start ")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
