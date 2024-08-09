package models

import "github.com/Jdcabreradev/jsonsocket/enums"

// SocketMessage represents the structure of the JSON message used for communication.
type SocketMessage struct {
	Flag    enums.SocketFlag `json:"flag"`    // Flag indicates the type of message.
	Payload interface{}      `json:"payload"` // Payload contains the actual data.
	Channel string           `json:"channel"` // Channel sends payload across clients subscribed to it.
}
