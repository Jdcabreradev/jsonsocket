package jsonsocket

import (
	"encoding/json"
	"net"
)

// socketSession represents a single session of a socket connection, which could be either a regular TCP connection or a secure TLS connection.
type socketSession struct {
	Socket net.Conn      // Socket is the underlying network connections.
	Reader *json.Decoder // reader is a JSON decoder for reading JSON-encoded messages from the connection.
	Writer *json.Encoder // writer is a JSON encoder for writing JSON-encoded messages to the connection.
}

// newSession creates a new instance of socketSession with the specified socket and timeout.
func newSession(conn net.Conn) *socketSession {
	return &socketSession{
		Socket: conn,
		Reader: json.NewDecoder(conn),
		Writer: json.NewEncoder(conn),
	}
}

// Read decodes a JSON-encoded object from the connection into the provided object parameter.
func (ss *socketSession) Read(object any) error {
	return ss.Reader.Decode(object)
}

// Write encodes the provided object as JSON and writes it to the connection.
func (ss *socketSession) Write(object any) error {
	return ss.Writer.Encode(object)
}

// Close closes the underlying network connection.
func (ss *socketSession) Close() error {
	return ss.Socket.Close()
}
