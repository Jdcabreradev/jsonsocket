package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Jdcabreradev/jsonsocket"
)

type MessageApp struct {
	GroupName string          `json:"group"`
	Members   map[string]bool `json:"members"`
	Sync      sync.RWMutex
}

func main() {
	client := &jsonsocket.SocketClient{ServerAddr: ":8000", UseTLS: false}
	client.Connect()
	chatGroup := &MessageApp{GroupName: "Test", Members: make(map[string]bool)}
	client.JoinGroup("Test")
	var username string
	input := bufio.NewReader(os.Stdin)
	fmt.Println("Type your username")
	fmt.Print("> ")
	username, _ = input.ReadString('\n')
	username = strings.TrimSpace(username)
	client.Send([]interface{}{map[string]string{"command": "newMember", "data": username}})
	go ListenToServer(client, chatGroup)
	client.Send([]interface{}{map[string]string{"command": "userList"}})

	for {
		fmt.Print("> ")
		command, _ := input.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "list":
			fmt.Println(chatGroup.Members)
		case "close":
			client.Close()
			return

		default:
			client.BroadcastToGroup(chatGroup.GroupName, map[string]string{"username": username, "data": command})
		}
	}
}

func ListenToServer(client *jsonsocket.SocketClient, chatGroup *MessageApp) {
	for {
		data, err := client.Await()
		if err != nil {
			fmt.Printf("Error receiving data from server: %v\n", err)
			break
		}

		switch data[0].(map[string]interface{})["command"] {
		case "newMember":
			payload := data[0].(map[string]interface{})["data"].(string)
			chatGroup.Members[payload] = true
		case "removeMember":
			payload := data[0].(map[string]interface{})["data"].(string)
			delete(chatGroup.Members, payload)
		case "userList":
			if data[0].(map[string]interface{})["data"] == nil {
				continue
			}
			payload := data[0].(map[string]interface{})["data"].(map[string]interface{})
			for _, v := range payload {
				chatGroup.Members[v.(string)] = true
			}
		default:
			fmt.Printf("%v:%v\n", data[0].(map[string]interface{})["username"], data[0].(map[string]interface{})["data"])
		}

	}
}
