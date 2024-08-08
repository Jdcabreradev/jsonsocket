package jsonsocket

import (
	"crypto/tls"
	"net"
	"time"
)

type SocketClient struct {
	client     SocketProcess
	ServerAddr string
	UseTLS     bool
	timeout    time.Duration
}

func (sc *SocketClient) Connect() error {
	var conn net.Conn
	var err error

	if sc.UseTLS {
		conn, err = tls.Dial("tcp", sc.ServerAddr, &tls.Config{})
	} else {
		conn, err = net.DialTimeout("tcp", sc.ServerAddr, sc.timeout)
	}

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return ErrConnectionTimeout
		}
		return ErrConnectionFailed
	}

	clientSession := &socketSession{Socket: conn}
	clientSession.Init()
	sc.client = SocketProcess{ID: "ClientManager", Session: clientSession}

	return nil
}

func (sc *SocketClient) Send(data []interface{}) error {
	return sc.client.Response(data)
}

func (sc *SocketClient) Await() ([]interface{}, error) {
	return sc.client.Listen()
}
