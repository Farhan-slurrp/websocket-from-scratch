package main

import (
	websocket "Farhan-slurrp/websocket-from-scratch/pkg"
	"fmt"
)

func main() {
	ws := websocket.NewWebSocket("/", "ws://localhost", "8000")
	for {
		webSocketConnection := ws.Accept()
		fmt.Println("Connected")
		if webSocketConnection != nil {
			for {
				message := webSocketConnection.Recv()
				fmt.Println(message)
			}
		}
	}
}
