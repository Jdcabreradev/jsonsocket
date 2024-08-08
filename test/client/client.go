package main

import (
	"fmt"
	"time"

	"github.com/Jdcabreradev/jsonsocket"
)

func main() {
	client := jsonsocket.SocketClient{ServerAddr: "10.152.164.143:8000", UseTLS: false}
	client.Connect()
	time.Sleep(time.Second * 2)
	var dataSend []interface{}
	dataSend = append(dataSend, "hola Server")
	client.Send(dataSend)
	data, _ := client.Await()
	fmt.Printf("Datos del servidor: %v", data)
}
