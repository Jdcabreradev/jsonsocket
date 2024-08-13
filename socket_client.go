package jsonsocket

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/Jdcabreradev/jsonsocket/enums"
)

// SocketClient represents a client connection to the server.
type SocketClient struct {
	client     *socketProcess
	ServerAddr string
	UseTLS     bool
	timeout    time.Duration
	groups     map[string]bool // Groups the client has joined
}

// Connect establishes a connection to the server.
func (sc *SocketClient) Connect() error {
	var conn net.Conn
	var err error

	if sc.UseTLS {
		conn, err = tls.Dial("tcp", sc.ServerAddr, &tls.Config{})
	} else {
		conn, err = net.DialTimeout("tcp", sc.ServerAddr, sc.timeout)
	}
	if err != nil {
		return err
	}

	clientSession := newSession(conn)
	sc.client = newSocketProcess("ClientManager", *clientSession, nil)
	sc.groups = make(map[string]bool)

	return nil
}

// Send sends data to the server.
func (sc *SocketClient) Send(data []interface{}) error {
	return sc.client.Response(data)
}

// Await waits for data from the server.
func (sc *SocketClient) Await() ([]interface{}, error) {
	return sc.client.Listen()
}

// JoinGroup subscribes the client to a group.
func (sc *SocketClient) JoinGroup(groupID string) error {
	sc.client.Session.Write(SocketMessage{Flag: enums.START_TX})
	msg := SocketMessage{
		Flag:  enums.SUB_CHANNEL,
		Group: groupID,
	}
	err := sc.client.Session.Write(msg)
	if err != nil {
		sc.groups[groupID] = true
	}
	return err
}

// LeaveGroup unsubscribes the client from a group.
func (sc *SocketClient) LeaveGroup(groupID string) error {
	sc.client.Session.Write(SocketMessage{Flag: enums.START_TX})
	msg := SocketMessage{
		Flag:  enums.UNSUB_CHANNEL,
		Group: groupID,
	}
	err := sc.client.Session.Write(msg)
	if err == nil {
		delete(sc.groups, groupID)
	}
	return err
}

// BroadcastToGroup sends a message to all clients in the specified group.
func (sc *SocketClient) BroadcastToGroup(groupID string, data interface{}) error {
	msg := SocketMessage{
		Flag:    enums.BROADCAST_CHANNEL,
		Group:   groupID,
		Payload: data,
	}
	sc.client.Session.Write(SocketMessage{Flag: enums.START_TX})
	return sc.client.Session.Write(msg)
}

// GetGroups returns the list of groups the client has joined.
func (sc *SocketClient) GetGroups() []string {
	var groupList []string
	for groupID := range sc.groups {
		groupList = append(groupList, groupID)
	}
	return groupList
}

// Close closes the client connection.
func (sc *SocketClient) Close() {
	sc.client.Close()
}
