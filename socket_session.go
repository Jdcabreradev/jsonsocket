package jsonsocket

import (
	"crypto/tls"
	"encoding/json"
	"net"

	"github.com/Jdcabreradev/logify/v2"
)

// socketSession represents a single session of a socket connection, which could be either a regular TCP connection or a secure TLS connection.
type socketSession struct {
	ID           string        // ID is a unique identifier for the session.
	Socket       net.Conn      // Socket is the underlying network connection for regular TCP connections.
	SecureSocket *tls.Conn     // SecureSocket is the underlying network connection for TLS connections.
	reader       *json.Decoder // reader is a JSON decoder for reading JSON-encoded messages from the connection.
	writer       *json.Encoder // writer is a JSON encoder for writing JSON-encoded messages to the connection.
	Logger       logify.Logger
}

// Init initializes the JSON encoder and decoder for the SocketSession.
// It checks if the connection is a secure TLS connection and sets up the appropriate reader and writer.
func (ss *socketSession) Init() {
	if ss.SecureSocket != nil {
		ss.reader = json.NewDecoder(ss.SecureSocket)
		ss.writer = json.NewEncoder(ss.SecureSocket)
		return
	}
	ss.reader = json.NewDecoder(ss.Socket)
	ss.writer = json.NewEncoder(ss.Socket)
}

// Read decodes a JSON-encoded object from the connection into the provided object parameter.
// It returns an error if the decoding fails.
func (ss *socketSession) Read(object any) error {
	return ss.reader.Decode(object)
}

// Write encodes the provided object as JSON and writes it to the connection.
// It returns an error if the encoding or writing fails.
func (ss *socketSession) Write(object any) error {
	return ss.writer.Encode(object)
}

// Close closes the underlying network connection.
// It returns an error if the connection cannot be closed.
func (ss *socketSession) Close() error {
	return ss.Socket.Close()
}
