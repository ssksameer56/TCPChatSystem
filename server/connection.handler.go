package server

import (
	"net"

	"github.com/ssksameer56/TCPChatSystem/models"
)

type Client models.Node

func (*Client) SendMessage() {

}

func (*Client) ReceiveMessage() {

}

func (*Client) HandleError() {

}

func (*Client) HandleConnection() {

}

func NewClient(name string, conn *net.Conn, buffSize int) *Client {
	client := Client{
		Name:           name,
		Connection:     conn,
		SendChannel:    make(chan models.Message, buffSize),
		ReceiveChannel: make(chan models.Message, buffSize),
	}
	return &client
}
