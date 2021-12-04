package client

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/ssksameer56/TCPChatSystem/models"
)

var ClientConfig models.ClientConfiguration

var ioHandler TerminalInput

func InitClient() error {
	file, _ := os.Open("client.settings.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	err := decoder.Decode(&ClientConfig)
	if err != nil {
		fmt.Println("Cant get config:", err.Error())
		return err
	}

	//Initialize I/O Handler to Terminal
	ioHandler = TerminalInput{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}
	return nil
}

func CreateChatConnection() (*Client, error) {
	conn, err := net.DialTimeout(ClientConfig.DefaultProtocol, fmt.Sprintf("%s:%s", ClientConfig.DefaultChatHost,
		ClientConfig.DefaultChatHostPort),
		time.Millisecond*200)
	if err != nil {
		fmt.Println("Connecting Error: ", err.Error())
		return &Client{}, err
	}
	client := NewClient(conn, ClientConfig.BufferSize)
	return client, nil
}

func RunChat(client *Client) {
	fmt.Println("Connected to Chat! Use C to start sending message, Enter to send the message")
	go client.ListenForMessageFromServer() //Start a routine to check for messages from server
	go client.ListenForInput(&ioHandler)   //Start a routine to get messages from input
	for {
		select {
		case data, ok := <-client.ReceiveChannel:
			if !ok {
				fmt.Println("Cant read message from server")
				break
			}
			ioHandler.DisplayMessage(data)
		case data := <-client.SendChannel:
			client.SendMessageToServer(data)
		default:
			time.Sleep(time.Millisecond * 1000)
		}
	}
}
