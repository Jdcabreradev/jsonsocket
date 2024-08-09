package jsonsocket

import (
	"encoding/json"
	"net"
	"time"
)

// socketSession represents a single session of a socket connection, which could be either a regular TCP connection or a secure TLS connection.
type socketSession struct {
	Socket  net.Conn      // Socket is the underlying network connections.
	Timeout time.Duration // Timeout defines how long a connection await for reading and writing processes.
	reader  *json.Decoder // reader is a JSON decoder for reading JSON-encoded messages from the connection.
	writer  *json.Encoder // writer is a JSON encoder for writing JSON-encoded messages to the connection.
}

// Init initializes the JSON encoder and decoder for the SocketSession.
func (ss *socketSession) Init() {
	ss.reader = json.NewDecoder(ss.Socket)
	ss.writer = json.NewEncoder(ss.Socket)
}

// Read decodes a JSON-encoded object from the connection into the provided object parameter.
func (ss *socketSession) Read(object any) error {

	return ss.reader.Decode(object)
}

// Write encodes the provided object as JSON and writes it to the connection.
func (ss *socketSession) Write(object any) error {
	return ss.writer.Encode(object)
}

// Close closes the underlying network connection.
func (ss *socketSession) Close() error {
	return ss.Socket.Close()
}
