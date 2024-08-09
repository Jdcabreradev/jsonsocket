package jsonsocket

import (
	"encoding/json"
	"net"
	"time"

	"github.com/Jdcabreradev/jsonsocket/errors"
)

// socketSession represents a single session of a socket connection, which could be either a regular TCP connection or a secure TLS connection.
type socketSession struct {
	socket  net.Conn      // Socket is the underlying network connections.
	timeout time.Duration // Timeout defines how long a connection await for reading and writing processes.
	reader  *json.Decoder // reader is a JSON decoder for reading JSON-encoded messages from the connection.
	writer  *json.Encoder // writer is a JSON encoder for writing JSON-encoded messages to the connection.
}

// newSession creates a new instance of socketSession with the specified socket and timeout.
func newSession(conn net.Conn, timeout time.Duration) *socketSession {
	return &socketSession{
		socket:  conn,
		reader:  json.NewDecoder(conn),
		writer:  json.NewEncoder(conn),
		timeout: timeout,
	}
}

// Read decodes a JSON-encoded object from the connection into the provided object parameter.
func (ss *socketSession) Read(object any) error {
	if ss.timeout > 0 {
		ss.socket.SetReadDeadline(time.Now().Add(ss.timeout))
	}
	err := ss.reader.Decode(object)
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return errors.TimeoutError
	}
	return err
}

// Write encodes the provided object as JSON and writes it to the connection.
func (ss *socketSession) Write(object any) error {
	if ss.timeout > 0 {
		ss.socket.SetWriteDeadline(time.Now().Add(ss.timeout))
	}
	err := ss.writer.Encode(object)
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return errors.TimeoutError
	}
	return err
}

// Close closes the underlying network connection.
func (ss *socketSession) Close() error {
	return ss.socket.Close()
}
