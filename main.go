// websockets.go
package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const SOCKET = ":8080"

var clients = make(map[*websocket.Conn]bool)

func main() {
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", homePage)
	fmt.Println("Server started on ", SOCKET)

	errServer := http.ListenAndServe(SOCKET, nil)
	if errServer != nil {
		panic("Error starting server: " + errServer.Error())
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "websockets.html")
}

func echo(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	defer conn.Close()

	clients[conn] = true

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			delete(clients, conn)
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// Write message back to browser
		for client := range clients {
			if err = client.WriteMessage(msgType, msg); err != nil {
				return
			}
		}

	}

}
