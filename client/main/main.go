package main

import (
	"os"

	"github.com/ssksameer56/TCPChatSystem/client"
)

func main() {
	err := client.InitClient()
	if err != nil {
		os.Exit(1)
	}
	var chatClient *client.Client
	chatClient, err = client.CreateChatConnection()
	if err != nil {
		os.Exit(1)
	}
	go client.RunChat(chatClient)
	(*chatClient.Connection).Close()
}
