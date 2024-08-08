package jsonsocket

// ISocketProcess defines an interface for socket operations.
// It includes methods for handling socket communication:
// - Listen: Reads and returns data from the socket.
// - Response: Sends data to the socket and returns a success status.
// - Close: Closes the socket connection and returns a success status.
type ISocketProcess interface {
	// Listen reads data from the socket and returns it as a slice of interfaces.
	// The actual implementation will handle how the data is read and processed.
	Listen() map[string]interface{}

	// Response sends data to the socket.
	// It takes a slice of interfaces as input and returns a boolean indicating success.
	Response(data map[string]interface{}) bool

	// Close closes the socket connection.
	// It returns a boolean indicating whether the closure was successful.
	Close() bool
}
