package websocket

import (
	"bytes"
	"encoding/binary"
	"net"
)

type WebSocketConnection struct {
	socket net.Conn
}

func NewWebSocketConnection(socket net.Conn) *WebSocketConnection {
	return &WebSocketConnection{
		socket: socket,
	}
}

func (wsConn *WebSocketConnection) Close() error {
	return wsConn.Close()
}

func (wsConn *WebSocketConnection) Recv() string {
	// finAndOpCode
	_, err := wsConn.readBytes(1)
	if err != nil {
		return ""
	}

	maskAndLengthIndicator, err := wsConn.readBytes(1)
	if err != nil {
		return ""
	}

	lengthIndicator := int(maskAndLengthIndicator[0] & 0x7F)

	var payloadLength int
	if lengthIndicator <= 125 {
		payloadLength = lengthIndicator
	} else if lengthIndicator == 126 {
		lengthBytes, err := wsConn.readBytes(2)
		if err != nil {
			return ""
		}
		payloadLength = int(binary.BigEndian.Uint16(lengthBytes))
	} else {
		lengthBytes, err := wsConn.readBytes(8)
		if err != nil {
			return ""
		}
		payloadLength = int(binary.BigEndian.Uint64(lengthBytes))
	}

	maskKey, err := wsConn.readBytes(4)
	if err != nil {
		return ""
	}

	encodedPayload, err := wsConn.readBytes(payloadLength)
	if err != nil {
		return ""
	}

	decodedPayload := make([]byte, payloadLength)
	for i := 0; i < payloadLength; i++ {
		decodedPayload[i] = encodedPayload[i] ^ maskKey[i%4]
	}

	return string(decodedPayload)
}

func (wsConn *WebSocketConnection) Send(message string) error {
	var buffer bytes.Buffer
	buffer.WriteByte(129)

	messageSize := len(message)

	if messageSize < 125 {
		buffer.WriteByte(byte(messageSize))
	} else if messageSize < 1<<16 {
		buffer.WriteByte(byte(126))

		err := binary.Write(&buffer, binary.BigEndian, uint16(messageSize))
		if err != nil {
			return err
		}
	} else {
		buffer.WriteByte(byte(127))

		err := binary.Write(&buffer, binary.BigEndian, uint64(messageSize))
		if err != nil {
			return err
		}
	}

	buffer.Write([]byte(message))

	_, err := wsConn.socket.Write(buffer.Bytes())
	return err
}

func (wsConn *WebSocketConnection) readBytes(n int) ([]byte, error) {
	buffer := make([]byte, n)
	_, err := wsConn.socket.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
