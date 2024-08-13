package errors

// SocketError is a custom enum-like type for socket errors.
type SocketError int

const (
	ProtocolError       SocketError = iota // Error for protocol violations
	DisconnError                           // Error for unexpected disconnections
	ClientNotFoundError                    // Error for missing socketProcess
)

// string returns a string representation of the SocketError.
func (e SocketError) string() string {
	switch e {
	case ProtocolError:
		return "Protocol error occurred"
	case DisconnError:
		return "Client disconnected unexpectedly"
	case ClientNotFoundError:
		return "Client not found in server"
	default:
		return "Unknown socket error"
	}
}

// Error implements the error interface for SocketError.
func (e SocketError) Error() string {
	return e.string()
}
