package main

import (
	"fmt"

	"github.com/Jdcabreradev/jsonsocket"
)

func main() {
	server := jsonsocket.SocketServer{Name: "testServer", Port: 8000, EnableTLS: false}
	server.Start()
	for {
		newClient := server.Bind()
		fmt.Printf("New client: %v", newClient)
		go HandleConn(&server, newClient)
	}
}

func HandleConn(server *jsonsocket.SocketServer, clientID string) {
	data, _ := server.Receive(clientID)
	fmt.Printf("Data from %v: %v\n", clientID, data)
	server.Reply(clientID, []interface{}{"Hola client"})
	server.Close(clientID)
}
