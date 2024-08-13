package jsonsocket

import (
	"io"

	"github.com/Jdcabreradev/jsonsocket/enums"
	"github.com/Jdcabreradev/jsonsocket/errors"
	"github.com/Jdcabreradev/logify/v3"
)

// socketProcess provides common implementations for socket operations.
type socketProcess struct {
	Id               string                  // ID is a unique identifier for the process.
	Session          *socketSession          // The actual session associated with this process.
	SubscribedGroups map[string]*SocketGroup // Groups the client is subscribed to.
	socketBroker     *socketBroker           // Channel to communicate with the server and other clients.
	Logger           logify.Logger           // Logger is used to log events and errors related to the socket process.
}

// newSocketProcess creates a new instance of socketSession with the specified id, session and manager.
func newSocketProcess(id string, session socketSession, socketBroker *socketBroker) *socketProcess {
	return &socketProcess{
		Id:               id,
		Session:          &session,
		SubscribedGroups: make(map[string]*SocketGroup),
		socketBroker:     socketBroker,
	}
}

// Listen reads and processes data from the socket session.
func (sp *socketProcess) Listen() ([]interface{}, error) {
	// Initialize variables to store incoming data and track the processing state.
	var dataList []interface{}
	var initFlag bool
	nextData := true

	for nextData {
		// Create a SocketMessage to store the incoming message from the session.
		var message SocketMessage
		err := sp.Session.Read(&message)
		if err != nil {
			if err == io.EOF {
				return nil, errors.DisconnError
			}
			return nil, err
		}

		// If this is the first message, check if the transmission has started.
		if !initFlag {
			if message.Flag != enums.START_TX {
				return nil, errors.ProtocolError
			}
			initFlag = true
			continue
		}

		// Handle different flags within the received message.
		switch message.Flag {
		case enums.TX:
			// Add payload data to the list.
			dataList = append(dataList, message.Payload)
		case enums.END_TX:
			// End of transmission, stop receiving further data.
			nextData = false
		case enums.CLOSE_CONN:
			// Close the connection and remove the client from the broker.
			sp.socketBroker.RemoveClient(sp.Id)
			return dataList, nil
		case enums.SUB_CHANNEL:
			// Subscribe the client to a specified group.
			sp.socketBroker.SubscribeToGroup(message.Group, sp.Id)
			dataList = []interface{}{}
			initFlag = false
			continue
		case enums.BROADCAST_CHANNEL:
			// Broadcast a message to all clients in the specified group.
			sp.socketBroker.BroadcastToGroup(message.Group, []interface{}{message.Payload})
			dataList = []interface{}{}
			initFlag = false
			continue
		case enums.UNSUB_CHANNEL:
			// Unsubscribe the client from a specified group.
			sp.socketBroker.UnsubscribeFromGroup(message.Group, sp.Id)
			dataList = []interface{}{}
			initFlag = false
			continue
		default:
			// If an unknown flag is encountered, return a protocol error.
			return nil, errors.ProtocolError
		}
	}

	return dataList, nil
}

// Response sends data to the socket session.
func (sp *socketProcess) Response(data []interface{}) error {
	// Start building a response by initializing the message list with a START_TX flag.
	var dataResponse []SocketMessage
	dataResponse = append(dataResponse, SocketMessage{Flag: enums.START_TX})

	// Iterate over the data to create messages for each item.
	for _, v := range data {
		dataResponse = append(dataResponse, SocketMessage{Flag: enums.TX, Payload: v})
	}

	// Add an END_TX message to indicate the end of the transmission.
	dataResponse = append(dataResponse, SocketMessage{Flag: enums.END_TX})

	// Send each message in the dataResponse list through the session.
	for _, v := range dataResponse {
		err := sp.Session.Write(v)
		if err != nil {
			return err
		}
	}

	return nil
}

// Close terminates the socket session.
func (sp *socketProcess) Close() error {
	return sp.Session.Close()
}
