package jsonsocket

// Flag definitions
const (
	confirmHandshake      uint8 = 0 // 0: Confirm handshake
	startDataTransmission uint8 = 1 // 1: Start data transmission
	endDataTransmission   uint8 = 2 // 2: End data transmission
	closeConnection       uint8 = 3 // 3: Close connection with the server
)

// socketMessage represents the structure of the JSON message used for communication
type socketMessage struct {
	Flag    *uint8      `json:"Flag"`    // Flag indicates the type of message
	Payload interface{} `json:"Payload"` // Payload contains the actual data
}
