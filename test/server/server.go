package main

import (
	"fmt"
	"sync"

	"github.com/Jdcabreradev/jsonsocket"
)

type MessageApp struct {
	GroupName string            `json:"group"`
	Members   map[string]string `json:"members"`
	Sync      sync.RWMutex
}

func main() {
	server := &jsonsocket.SocketServer{Name: "LinuxMessageServer", Port: 8000, EnableTLS: false}
	server.Start()
	chatGroup := &MessageApp{GroupName: "Test", Members: make(map[string]string)}
	for {
		newClient := server.Bind()
		fmt.Println(newClient)
		go HandleChatApp(server, newClient, chatGroup)
	}
}

func HandleChatApp(server *jsonsocket.SocketServer, clientId string, chatGroup *MessageApp) {
	for {
		msg, err := server.Listen(clientId)
		if err != nil {
			fmt.Println("User deleted")
			chatGroup.Sync.Lock()
			server.Broadcast(chatGroup.GroupName, []interface{}{map[string]string{"command": "removeMember", "data": chatGroup.Members[clientId]}})
			delete(chatGroup.Members, clientId)
			chatGroup.Sync.Unlock()
			server.Close(clientId)
			return
		}

		switch msg[0].(map[string]interface{})["command"] {
		case "newMember":
			data := msg[0].(map[string]interface{})["data"].(string)
			chatGroup.Sync.Lock()
			chatGroup.Members[clientId] = data
			chatGroup.Sync.Unlock()
			fmt.Println("New user added")
			server.Broadcast(chatGroup.GroupName, msg)
		case "userList":
			server.Response(clientId, []interface{}{map[string]interface{}{"command": "userList", "data": chatGroup.Members}})
		}
	}
}
