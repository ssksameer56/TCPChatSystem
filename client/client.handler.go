package client

import (
	"fmt"
	"net"
	"time"

	"github.com/ssksameer56/TCPChatSystem/models"
)

type Client models.Node

var ioHandler TerminalInput

func (client *Client) HandleChat() {
	for {
		select {
		case data, ok := <-client.ReceiveChannel:
			if !ok {
				fmt.Println("Cant read message from server")
				break
			}
			client.ReceiveMessage(data)
		case data := <-client.SendChannel:
			client.SendMessage(data)

		default:
			time.Sleep(time.Millisecond * 1000)
		}
	}
}

//Function to recieve message from Server
func (client *Client) ReceiveMessage(data string) {

	ioHandler.DisplayMessage(data)
}

//Function to Send Message to Server
func (client *Client) SendMessage(data models.Message) bool {
	_, err := (*client.Connection).Write([]byte(data.MessageText))
	if err != nil {
		fmt.Println("Error in sending message to server: ")
		return false
	}
	return true
}

func NewClient(conn net.Conn, buffSize int, name string) *Client {
	return &Client{
		Name:           name,
		Connection:     &conn,
		SendChannel:    make(chan models.Message, buffSize),
		ReceiveChannel: make(chan string, buffSize),
		SignalChannel:  make(chan int, buffSize),
	}
}
