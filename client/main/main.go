package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ssksameer56/TCPChatSystem/models"
)

var ClientConfig models.ClientConfiguration

func InitClient() error {
	file, _ := os.Open("client.settings.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	err := decoder.Decode(&ClientConfig)
	if err != nil {
		fmt.Println("Cant get config:", err.Error())
		return err
	}
	return nil
}
