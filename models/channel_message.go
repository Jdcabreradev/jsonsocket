package models

import "github.com/Jdcabreradev/jsonsocket/enums"

// ServerCallback encapsulates the data from a socketProcess to the server.
type ServerCallback struct {
	ChannelID string           // ID of the channel
	ClientID  string           // ID of the client (SocketProcess)
	Flag      enums.SocketFlag // The flag indicating the type of action
	Payload   interface{}      // Data to be sent to the channel
}
