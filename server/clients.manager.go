package server

import (
	"bufio"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ssksameer56/TCPChatSystem/models"
)

type ClientsManager struct {
	AllClients    map[string]*Client //Map of All Clients
	ServerChannel chan Message       //Channel to send a receive messages between client and server
}

//Handle data for All Clients
func (manager *ClientsManager) HandleClients(message Message) {
	for {
		select {
		case data, ok := <-manager.ServerChannel:
			if !ok {
				log.WithFields(log.Fields{"client": "Server"}).Error("Error when reading message from a client")
				break
			}
			if data.MessageText == models.END_CHAT {
				client := manager.AllClients[data.ClientName]
				client.SignalChannel <- models.END
				manager.RemoveClient(data.ClientName)
			} else {
				for _, client := range manager.AllClients {
					select {
					case client.SendChannel <- data:
					default:
						log.WithFields(log.Fields{"client": client.Name}).Info("Buffer full. discarding message: " + data.MessageText)
					}
				}
			}
		default:
			time.Sleep(time.Millisecond * 1000)
		}
	}
}

//Add Client to the List
func (manager *ClientsManager) AddClientToList(client *Client) {
	log.WithFields(log.Fields{"client": client.Name}).Info("Adding client to connection list")
	manager.AllClients[client.Name] = client
}

//Remove Client from the List
func (manager *ClientsManager) RemoveClient(name string) {
	log.WithFields(log.Fields{"client": name}).Info("Removing client from connection list")
	delete(manager.AllClients, name)
}

//Check if client name is unique
func (manager *ClientsManager) CheckClientName(name string) bool {
	exists := false
	for clientName := range manager.AllClients {
		if name == clientName {
			exists = true
		}
	}
	return exists
}

func (manager *ClientsManager) CreateClient(conn net.Conn) error {
	for {
		conn.Write([]byte("Please Enter a name for client"))
		reader := bufio.NewReader(conn)
		name, _, err := reader.ReadLine()
		if err != nil {
			log.WithFields(log.Fields{"IP": conn.RemoteAddr()}).Error(err.Error())
			return err
		}
		if !manager.CheckClientName(string(name)) {
			log.WithFields(log.Fields{"client": name}).Info("Creating Client for " + conn.RemoteAddr().String())
			client := NewClient(string(name), &conn, server.BufferSize, manager.ServerChannel)
			manager.AddClientToList(client)
			return nil
		}
	}
}
