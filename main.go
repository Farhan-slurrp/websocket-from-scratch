package main

import (
	websocket "Farhan-slurrp/websocket-from-scratch/pkg"
	"fmt"
)

func main() {
	ws := websocket.NewWebSocket("/", "8000")
	for {
		webSocketConnection := ws.Accept()
		if webSocketConnection != nil {
			go func(connection *websocket.WebSocketConnection) {
				fmt.Println("Connected")
				message := connection.Recv()
				for message != "" {
					// fmt.Println(message)
					message = connection.Recv()
				}
			}(webSocketConnection)
		}
	}
}
