package interfaces

// ISocketProcess defines an interface for socket operations.
type ISocketProcess interface {
	Listen() map[string]interface{} // Listen reads data from the socket and returns it as a slice of interfaces.

	Response(data map[string]interface{}) bool // Response sends data to the socket.

	Close() bool // Close closes the socket connection.

}
