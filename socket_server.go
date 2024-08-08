package jsonsocket

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"
)

type SocketServer struct {
	Name        string
	Port        uint16
	EnableTLS   bool
	CertPath    string
	KeyPath     string
	Listener    net.Listener
	Logger      string
	SubProcess  map[string]*SocketProcess
	SyncProcess sync.Mutex
}

func (ss *SocketServer) Start() error {
	ss.SubProcess = make(map[string]*SocketProcess)

	if ss.EnableTLS {
		cert, err := tls.LoadX509KeyPair(ss.CertPath, ss.KeyPath)
		if err != nil {
			return err
		}
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
		ss.Listener, err = tls.Listen("tcp", fmt.Sprintf(":%d", ss.Port), tlsConfig)
		if err != nil {
			return err
		}
		return err
	}
	var err error
	ss.Listener, err = net.Listen("tcp", fmt.Sprintf(":%d", ss.Port))
	if err != nil {
		return err
	}
	fmt.Println("Server Started!")
	return nil
}

func (ss *SocketServer) Bind() string {
	for {
		conn, err := ss.Listener.Accept()
		if err != nil {
			//TODO PRINT LOGGER
			continue
		}

		clientAddr := conn.RemoteAddr().String()
		clientSession := socketSession{Socket: conn}
		clientSession.Init()
		ClientSocketProcess := &SocketProcess{ID: clientAddr, Session: &clientSession}

		ss.SyncProcess.Lock()
		ss.SubProcess[clientAddr] = ClientSocketProcess
		ss.SyncProcess.Unlock()

		return clientAddr
	}
}

func (ss *SocketServer) Receive(ProcessID string) ([]interface{}, error) {
	session, exists := ss.SubProcess[ProcessID]
	if !exists {
		return nil, fmt.Errorf("session %s not found", ProcessID)
	}
	return session.Listen()
}

func (ss *SocketServer) Reply(ProcessID string, data []interface{}) error {
	session, exists := ss.SubProcess[ProcessID]
	if !exists {
		return fmt.Errorf("session %s not found", ProcessID)
	}
	return session.Response(data)
}

func (ss *SocketServer) Close(ProcessID string) error {
	session, exists := ss.SubProcess[ProcessID]
	if !exists {
		return fmt.Errorf("session %s not found", ProcessID)
	}
	return session.Close()

}
