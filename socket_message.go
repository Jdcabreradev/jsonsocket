package jsonsocket

// SocketFlagType represents different types of socket flags.
type SocketFlagType uint8

const (
	// Define different socket flags types
	UndefinedFlag SocketFlagType = iota // Undefined flag
	StartDataFlag                       // Start Data transmission flag
	DataFlag                            // ontinuous data transmission flag
	EndDataFlag                         // End data transmission
	CloseFlag                           // Close connection flag
)

// SocketMessage represents the structure of the JSON message used for communication
type SocketMessage struct {
	Flag    SocketFlagType `json:"Flag"`    // Flag indicates the type of message
	Payload interface{}    `json:"Payload"` // Payload contains the actual data
}
