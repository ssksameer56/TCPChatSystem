package client

import (
	"fmt"
	"net"

	"github.com/ssksameer56/TCPChatSystem/models"
)

type Client models.Client

//Recieve message from Server
func (client *Client) ListenForMessageFromServer() {
	for {
		var data []byte
		if _, err := (*client.Connection).Read(data); err != nil {
			fmt.Println("Error getting message: ", err.Error())
			continue
		}
		client.ReceiveChannel <- string(data)
	}
}

//Read from any input
func (client *Client) ListenForInput(io models.InputOutputHandler) {
	for {
		c, err := io.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message: ", err.Error())
			continue
		}
		client.SendChannel <- string(c)
	}
}

//Send Message to Server
func (client *Client) SendMessageToServer(data string) bool {
	_, err := (*client.Connection).Write([]byte(data))
	if err != nil {
		fmt.Println("Error in sending message to server: ")
		return false
	}
	return true
}

//Return a new client to handle connection to server
func NewClient(conn net.Conn, buffSize int) *Client {
	return &Client{
		Connection:     &conn,
		SendChannel:    make(chan string, buffSize),
		ReceiveChannel: make(chan string, buffSize),
	}
}
