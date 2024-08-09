package jsonsocket

import (
	"fmt"

	"github.com/Jdcabreradev/logify/v3"
)

// SocketProcess provides common implementations for socket operations.
type SocketProcess struct {
	ID      string         // ID is a unique identifier for the process.
	Session *socketSession // The actual session associated with this process.
	Logger  *logify.Logger // Logger to handle log messages.
}

// Listen reads and processes data from the socket session.
func (sp *SocketProcess) Listen() ([]interface{}, error) {
	var dataList []interface{}
	var initFlag bool
	nextData := true

	for nextData {
		var socketMessage SocketMessage

		err := sp.Session.reader.Decode(&socketMessage)
		if err != nil {
			return nil, ErrEOF
		}
		if socketMessage.Flag == UndefinedFlag ||
			socketMessage.Payload == nil && socketMessage.Flag == TX {
			return nil, ErrInvalidData
		}

		if !initFlag {
			if socketMessage.Flag != START_TX {
				return nil, ErrProtocolError
			}
			initFlag = true
			continue
		}
		fmt.Println(socketMessage)
		switch socketMessage.Flag {
		case START_TX:
			return nil, ErrProtocolError
		case TX:
			dataList = append(dataList, socketMessage.Payload)
		case END_TX:
			nextData = false
		case CLOSE_CONN:
			defer sp.Close()
			nextData = false
		default:
			return nil, ErrProtocolError
		}
	}

	return dataList, nil
}

// Response sends data to the socket session.
func (sp *SocketProcess) Response(data []interface{}) error {
	var dataResponse []SocketMessage
	dataResponse = append(dataResponse, SocketMessage{Flag: START_TX})
	for _, v := range data {
		dataResponse = append(dataResponse, SocketMessage{Flag: TX, Payload: v})
	}
	dataResponse = append(dataResponse, SocketMessage{Flag: END_TX})
	fmt.Println(dataResponse)
	for _, v := range dataResponse {
		err := sp.Session.writer.Encode(v)
		if err != nil {
			return ErrEOF
		}
	}

	return nil
}

// Close closes the socket session.
func (sp *SocketProcess) Close() error {
	return sp.Session.Close()
}
