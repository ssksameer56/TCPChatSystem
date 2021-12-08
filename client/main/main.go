package main

import (
	"os"
	"sync"

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
	wg := sync.WaitGroup{}
	wg.Add(1)
	go client.RunChat(chatClient, &wg)
	wg.Wait()
	(*chatClient.Connection).Close()
}
