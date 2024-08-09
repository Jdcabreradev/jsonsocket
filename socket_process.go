package jsonsocket

import (
	"fmt"

	"github.com/Jdcabreradev/jsonsocket/enums"
	"github.com/Jdcabreradev/jsonsocket/models"
	"github.com/Jdcabreradev/logify/v3"
	logify_enums "github.com/Jdcabreradev/logify/v3/enums"
)

// SocketProcess provides common implementations for socket operations.
type SocketProcess struct {
	Id            string                     // ID is a unique identifier for the process.
	Session       *socketSession             // The actual session associated with this process.
	ParentChannel chan models.ServerCallback // Channel to communicate with its parent (server or client).
	Logger        *logify.Logger             // Logger to handle log messages.
}

// Listen reads and processes data from the socket session.
func (sp *SocketProcess) Listen() []interface{} {
	var dataList []interface{}
	var initFlag bool
	nextData := true

	for nextData {
		var socketMessage models.SocketMessage

		err := sp.Session.reader.Decode(&socketMessage)
		if err != nil {
			sp.Logger.Log(logify_enums.ERROR, sp.Id+": "+err.Error())
			return nil
		}

		if !initFlag {
			if socketMessage.Flag != enums.START_TX {
				sp.Logger.Log(logify_enums.ERROR, sp.Id+": Protocol Error")
				return nil
			}
			initFlag = true
			continue
		}

		switch socketMessage.Flag {
		case enums.START_TX:
			sp.Logger.Log(logify_enums.ERROR, sp.Id+": Protocol Error")
			return nil
		case enums.TX:
			dataList = append(dataList, socketMessage.Payload)
		case enums.END_TX:
			nextData = false
		case enums.CLOSE_CONN:
			defer sp.Close()
			nextData = false
		default:
			return nil, errors.ErrProtocolError
		}
	}

	return dataList, nil
}

// Response sends data to the socket session.
func (sp *SocketProcess) Response(data []interface{}) error {
	var dataResponse []models.SocketMessage
	dataResponse = append(dataResponse, models.SocketMessage{Flag: enums.START_TX})
	for _, v := range data {
		dataResponse = append(dataResponse, models.SocketMessage{Flag: enums.TX, Payload: v})
	}
	dataResponse = append(dataResponse, models.SocketMessage{Flag: enums.END_TX})
	fmt.Println(dataResponse)
	for _, v := range dataResponse {
		err := sp.Session.writer.Encode(v)
		if err != nil {
			return errors.ErrEOF
		}
	}

	return nil
}

// Close closes the socket session.
func (sp *SocketProcess) Close() error {
	return sp.Session.Close()
}
