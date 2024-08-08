package jsonsocket

// SocketErrorType defines various types of socket errors.
type SocketErrorType int

const (
	ErrConnectionFailed SocketErrorType = iota
	ErrInvalidData
	ErrEOF
	ErrProtocolError
	ErrConnectionTimeout
	ErrConnectionRefused
	ErrNetworkUnreachable
	ErrAddressInUse
	ErrProtocolMismatch
	ErrSocketClosed
	ErrMessageTooLarge
	ErrUnknown
)

// Error messages associated with each error type.
var socketErrorMessages = map[SocketErrorType]string{
	ErrConnectionFailed:   "Connection failed",
	ErrInvalidData:        "Invalid data received",
	ErrEOF:                "End of file",
	ErrProtocolError:      "Protocol error",
	ErrConnectionTimeout:  "Connection timed out",
	ErrConnectionRefused:  "Connection refused",
	ErrNetworkUnreachable: "Network unreachable",
	ErrAddressInUse:       "Address in use",
	ErrProtocolMismatch:   "Protocol mismatch",
	ErrSocketClosed:       "Socket closed",
	ErrMessageTooLarge:    "Message too large",
	ErrUnknown:            "Unknown error",
}

// Error returns the error message for the given error type.
func (e SocketErrorType) Error() string {
	msg, exists := socketErrorMessages[e]
	if !exists {
		return "Unknown error"
	}
	return msg
}
