package websocket

import (
	"fmt"
	"net"
)

type WebSocket struct {
	Path      string
	TcpServer net.Listener
}

func NewWebSocket(path string, port string) *WebSocket {
	l, err := net.Listen("tcp", port)
	if err != nil {
		panic(fmt.Sprintf("error on starting a new server on port %s", port))
	}
	return &WebSocket{
		Path:      path,
		TcpServer: l,
	}
}

func (ws *WebSocket) Accept() {
	socket, err := ws.TcpServer.Accept()
	if err != nil {
		panic("error on accepting the connection")
	}
	if ws.sendHandshake(socket) {
		// create a new websocket connection here
	}
}

func (ws *WebSocket) sendHandshake(socket net.Conn) bool {
	// check get header
	return false
}
