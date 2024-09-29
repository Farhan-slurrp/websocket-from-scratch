package main

import (
	websocket "Farhan-slurrp/websocket-from-scratch/pkg"
	"fmt"
)

func main() {
	ws := websocket.NewWebSocket("/", "8000")
	fmt.Println(ws)
}
