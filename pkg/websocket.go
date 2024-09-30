package websocket

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net"
	"regexp"
	"strings"
)

const WS_MAGIC_STRING = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

type WebSocket struct {
	Path      string
	TcpServer net.Listener
}

func NewWebSocket(path string, host string, port string) *WebSocket {
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(fmt.Sprintf("error on starting a new server on port %s", port))
	}

	return &WebSocket{
		Path:      path,
		TcpServer: l,
	}
}

func (ws *WebSocket) Accept() *WebSocketConnection {
	socket, err := ws.TcpServer.Accept()
	if err != nil {
		panic("error on accepting the connection")
	}
	if ws.sendHandshake(socket) {
		return NewWebSocketConnection(socket)
	}
	return nil
}

func (ws *WebSocket) sendHandshake(socket net.Conn) bool {
	reader := bufio.NewReader(socket)

	_, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	header, err := getHeader(reader)
	if err != nil {
		return false
	}
	getReqRe := regexp.MustCompile(fmt.Sprintf(`GET %s HTTP/1.1`, ws.Path))
	secretKeyRe := regexp.MustCompile(`Sec-WebSocket-Key: (.*)\r\n`)

	getReqMatches := getReqRe.FindStringSubmatch(header)
	secretKeyMatches := secretKeyRe.FindStringSubmatch(header)

	if len(getReqMatches) > 0 || len(secretKeyMatches) > 0 {
		wsAccept := createWebSocketAccept(secretKeyMatches[1])
		sendHandshakeResponse(socket, wsAccept)
		return true
	}

	send400(socket)
	return false
}

func getHeader(reader *bufio.Reader) (string, error) {
	header := strings.Builder{}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		if line == "\r\n" {
			break
		}
		header.WriteString(line)
	}

	return header.String(), nil
}

func createWebSocketAccept(key string) string {
	// Create SHA-1 hash
	h := sha1.New()
	h.Write([]byte(key + WS_MAGIC_STRING))
	digest := h.Sum(nil)

	// Encode to Base64
	acceptKey := base64.StdEncoding.EncodeToString(digest)

	return acceptKey
}

func send400(socket net.Conn) {
	msg := "HTTP/1.1 400 Bad Request\r\n" +
		"Content-Type: text/plain\r\n" +
		"Connection: close\r\n" +
		"\r\n" +
		"Incorrect request"
	socket.Write([]byte(msg))
	socket.Close()
}

func sendHandshakeResponse(socket net.Conn, wsAccept string) {
	msg := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\n"+
		"Upgrade: websocket\r\n"+
		"Connection: Upgrade\r\n"+
		"Sec-WebSocket-Accept: %s\r\n", wsAccept)
	socket.Write([]byte(msg))
}
