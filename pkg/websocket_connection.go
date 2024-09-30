package websocket

import "net"

type WebSocketConnection struct {
	socket net.Conn
}

func NewWebSocketConnection(socket net.Conn) *WebSocketConnection {
	return &WebSocketConnection{
		socket: socket,
	}
}

func (wsC *WebSocketConnection) Recv() string {
	return "here"
}
