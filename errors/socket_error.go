package errors

// SocketError is a custom enum-like type for socket errors.
type SocketError int

const (
	ProtocolError SocketError = iota // Error for protocol violations
	TimeoutError                     // Error for connection timeouts
	ReadError                        // Error for reader exception
	WriteError                       // Error for writer exception
	DisconnError                     // Error for unexpected disconnections
)

// string returns a string representation of the SocketError.
func (e SocketError) string() string {
	switch e {
	case ProtocolError:
		return "Protocol error occurred"
	case TimeoutError:
		return "Connection timed out"
	case ReadError:
		return "Reader cannot bind socket data"
	case WriteError:
		return "Writer cannot send data to socket"
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
