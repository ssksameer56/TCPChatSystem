package server

import (
	log "github.com/sirupsen/logrus"
)

type ClientsManager struct {
	AllClients map[string]*Client
}

func (manager *ClientsManager) BroadcastMessage(message Message) {
	for _, client := range manager.AllClients {
		select {
		case client.SendChannel <- message:
		default:
			log.WithFields(log.Fields{"client": client.Name}).Info("Buffer full. discarding message: " + message.MessageText)
		}
	}
}

func (manager *ClientsManager) AddClient(client *Client) {
	manager.AllClients[client.Name] = client
}

func (manager *ClientsManager) RemoveClient(name string) {
	log.WithFields(log.Fields{"client": name}).Info("Removing client from connections")
	delete(manager.AllClients, name)
}

func (manager *ClientsManager) CheckClientName(name string) bool {
	exists := false
	for clientName := range manager.AllClients {
		if name == clientName {
			exists = true
		}
	}
	return exists
}
