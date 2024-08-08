package jsonsocket

import "fmt"

// SocketProcess provides common implementations for socket operations.
type SocketProcess struct {
	ID      string         // ID is a unique identifier for the process
	Session *socketSession // The actual session associated with this process

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
			socketMessage.Payload == nil && socketMessage.Flag == DataFlag {
			return nil, ErrInvalidData
		}

		if !initFlag {
			if socketMessage.Flag != StartDataFlag {
				return nil, ErrProtocolError
			}
			initFlag = true
			continue
		}
		fmt.Println(socketMessage)
		switch socketMessage.Flag {
		case StartDataFlag:
			return nil, ErrProtocolError
		case DataFlag:
			dataList = append(dataList, socketMessage.Payload)
		case EndDataFlag:
			nextData = false
		case CloseFlag:
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
	dataResponse = append(dataResponse, SocketMessage{Flag: StartDataFlag})
	for _, v := range data {
		dataResponse = append(dataResponse, SocketMessage{Flag: DataFlag, Payload: v})
	}
	dataResponse = append(dataResponse, SocketMessage{Flag: EndDataFlag})
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
