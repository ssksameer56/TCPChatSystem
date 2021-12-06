package client

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
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
	fmt.Println("Connected to Chat! Use C and Enter to start sending message, Enter after to send the message")
	go client.ListenForMessageFromServer() //Start a routine to check for messages from server
	go client.ListenForInput(&ioHandler)   //Start a routine to get messages from input
	for {
		select {
		case data, ok := <-client.ReceiveChannel:
			if !ok {
				fmt.Println("Cant read message from server")
			} else if strings.EqualFold(data, "exit") {
				{
					return
				}
			} else {
				client.DisplayMessage(data, &ioHandler)
			}
		case data := <-client.SendChannel:
			if strings.EqualFold(data, "exit") {
				return
			}
			client.SendMessageToServer(data)
		default:
			time.Sleep(time.Millisecond * 1000)
		}
	}
}
