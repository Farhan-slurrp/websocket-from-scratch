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

func NewWebSocket(path string, port string) *WebSocket {
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
	if key, ok := ws.sendHandshake(socket); ok {
		sendHandshakeResponse(socket, key)
		return NewWebSocketConnection(socket)
	}

	send400(socket)
	return nil
}

func (ws *WebSocket) sendHandshake(socket net.Conn) (string, bool) {
	reader := bufio.NewReader(socket)

	header, err := getHeader(reader)
	if err != nil {
		return "", false
	}

	secretKeyRe := regexp.MustCompile(`Sec-WebSocket-Key: (.*)\r\n`)
	secretKeyMatches := secretKeyRe.FindStringSubmatch(header)
	if len(secretKeyMatches) > 0 {
		return secretKeyMatches[1], true
	}

	return "", false
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
	h := sha1.New()
	h.Write([]byte(key + WS_MAGIC_STRING))
	digest := h.Sum(nil)

	acceptKey := base64.StdEncoding.EncodeToString(digest)

	return acceptKey
}

func send400(socket net.Conn) {
	response := "HTTP/1.1 400 Bad Request\r\n" +
		"Content-Type: text/plain\r\n" +
		"Connection: close\r\n" +
		"\r\n" +
		"Incorrect request"

	socket.Write([]byte(response))
	socket.Close()
}

func sendHandshakeResponse(socket net.Conn, key string) {
	wsAccept := createWebSocketAccept(key)

	response := fmt.Sprintf(
		"HTTP/1.1 101 Switching Protocols\r\n"+
			"Upgrade: websocket\r\n"+
			"Connection: Upgrade\r\n"+
			"Sec-WebSocket-Accept: %s\r\n\r\n", wsAccept)

	_, err := socket.Write([]byte(response))
	if err != nil {
		fmt.Println("Failed to write handshake response:", err)
	}
}
