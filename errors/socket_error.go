package errors

// SocketError is a custom enum-like type for socket errors.
type SocketError int

const (
	ProtocolError SocketError = iota // Error for protocol violations
	TimeoutError                     // Error for connection timeouts
	DisconnError                     // Error for unexpected disconnections
)

// string returns a string representation of the SocketError.
func (e SocketError) string() string {
	switch e {
	case ProtocolError:
		return "Protocol error occurred"
	case TimeoutError:
		return "Connection timed out"
	case DisconnError:
		return "Client disconnected unexpectedly"
	default:
		return "Unknown socket error"
	}
}

// Error implements the error interface for SocketError.
func (e SocketError) Error() string {
	return e.string()
}
