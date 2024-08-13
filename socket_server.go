package jsonsocket

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/Jdcabreradev/logify/v3"
	logify_enums "github.com/Jdcabreradev/logify/v3/enums"
)

type SocketServer struct {
	Name         string
	Port         uint16
	EnableTLS    bool
	CertPath     string
	KeyPath      string
	listener     net.Listener
	isRunning    bool
	logger       logify.Logger
	socketBroker *socketBroker
}

// REFACTOR THIS
func (ss *SocketServer) Start() error {
	ss.logger = *logify.New(ss.Name, "ServerApp", "./logs/", logify_enums.LogModeVerbose)
	ss.socketBroker = newBroker(&ss.logger)

	if ss.EnableTLS {
		cert, err := tls.LoadX509KeyPair(ss.CertPath, ss.KeyPath)
		if err != nil {
			return err
		}
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
		ss.listener, err = tls.Listen("tcp", fmt.Sprintf(":%d", ss.Port), tlsConfig)
		if err != nil {
			return err
		}
		ss.isRunning = true
		return err
	}
	var err error
	ss.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", ss.Port))
	if err != nil {
		return err
	}
	ss.logger.Log(logify_enums.INFO, "Server started!")
	ss.isRunning = true
	return nil
}

// Bind listen to new connections and if succeed return the client id.
func (ss *SocketServer) Bind() string {
	for {
		conn, err := ss.listener.Accept()
		if err != nil {
			ss.logger.Log(logify_enums.WARNING, "Some connection was refused by the server")
			continue
		}

		clientId := conn.RemoteAddr().String()
		clientSession := newSession(conn)
		ClientSocketProcess := newSocketProcess(clientId, *clientSession, ss.socketBroker)

		ss.socketBroker.AddClient(clientId, ClientSocketProcess)

		ss.logger.Log(logify_enums.DEBUG, "New connection received: "+clientId)
		return clientId
	}
}

// Listen returns data from a specific client.
func (ss *SocketServer) Listen(clientId string) ([]interface{}, error) {
	return ss.socketBroker.clients[clientId].Listen()
}

func (ss *SocketServer) Response(ProcessID string, data []interface{}) error {
	client, exists := ss.socketBroker.clients[ProcessID]
	if !exists {
		return fmt.Errorf("client %s not found", ProcessID)
	}
	return client.Response(data)
}

func (ss *SocketServer) Broadcast(groupId string, data []interface{}) {
	ss.socketBroker.BroadcastToGroup(groupId, data)
}

func (ss *SocketServer) Close(ProcessID string) {
	client, exists := ss.socketBroker.clients[ProcessID]
	if !exists {
		ss.logger.Log(logify_enums.WARNING, "Client "+ProcessID+" not found for disconection process")
	}
	client.Close()
	ss.socketBroker.RemoveClient(ProcessID)
}
