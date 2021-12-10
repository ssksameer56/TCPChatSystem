package client

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"github.com/ssksameer56/TCPChatSystem/models"
)

type Client models.Client

var connMutex sync.Mutex

//Recieve message from Server
func (client *Client) ListenForMessageFromServer() {
	reader := bufio.NewReader(*client.Connection)
	for {
		var data string
		data, err := reader.ReadString('\n')
		if err != nil {
			if err == net.ErrClosed {
				client.ReceiveChannel <- "exit"
			} else {
				fmt.Println("Error getting message: ", err.Error())
			}
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
		client.SendChannel <- c
	}
}

func (client *Client) DisplayMessage(data string, handler models.InputOutputHandler) {
	err := handler.DisplayMessage(data)
	if err != nil {
		fmt.Println("Error displaying message: ", err.Error())
	}
}

//Send Message to Server
func (client *Client) SendMessageToServer(data []byte) bool {
	connMutex.Lock()
	_, err := (*client.Connection).Write([]byte(data))
	connMutex.Unlock()
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
		SendChannel:    make(chan []byte, buffSize),
		ReceiveChannel: make(chan string, buffSize),
	}
}
