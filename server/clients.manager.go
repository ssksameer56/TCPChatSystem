package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ssksameer56/TCPChatSystem/models"
)

type ClientsManager struct {
	AllClients    map[string]*Client //Map of All Clients
	ServerChannel chan Message       //Channel to send a receive messages between client and server
}

//Handle data for All Clients
func (manager *ClientsManager) HandleClients(wg *sync.WaitGroup, quit chan bool) {
	log.WithFields(log.Fields{"client": "server"}).Info("Listening for data from all clients")
	for {
		select {
		case data, ok := <-manager.ServerChannel:
			if !ok {
				log.WithFields(log.Fields{"client": "Server"}).Error("Error when reading message from a client")
				break
			}
			log.WithFields(log.Fields{"client": data.ClientName}).Info("Sending to all clients: " + data.MessageText)
			if data.MessageText == models.END_CHAT {
				manager.RemoveClient(data.ClientName)
			} else {
				for _, client := range manager.AllClients {
					select {
					case client.SendChannel <- data:
					default:
						log.WithFields(log.Fields{"client": client.Name}).Info("Send Buffer full. discarding message: " + data.MessageText)
					}
				}
			}
		case <-quit:
			wg.Done()
			return
		default:
			time.Sleep(time.Millisecond * 1000)
		}
	}
}

//Add Client to the List
func (manager *ClientsManager) AddClientToList(client *Client) {
	manager.AllClients[client.Name] = client
	log.WithFields(log.Fields{"client": client.Name}).Info("Added client to connection list")

}

//Remove Client from the List
func (manager *ClientsManager) RemoveClient(name string) {
	delete(manager.AllClients, name)
	log.WithFields(log.Fields{"client": name}).Info("Removed client from connection list")
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
	log.WithFields(log.Fields{"IP": conn.RemoteAddr()}).Info("Creating Client")
	reader := bufio.NewReader(conn)
	for {
		conn.Write([]byte("Please Enter a name for your client\n"))
		name, err := reader.ReadString('\n')
		name = strings.Replace(name, "\n", "", -1) // to remove the delimiter from name
		fmt.Println(name)
		if err != nil {
			log.WithFields(log.Fields{"IP": conn.RemoteAddr()}).Error(err.Error())
			return err
		}
		if !manager.CheckClientName(string(name)) {
			log.WithFields(log.Fields{"client": name}).Info("Creating Client for " + conn.RemoteAddr().String())
			client := NewClient(string(name), &conn, server.Config.BufferSize, manager.ServerChannel)
			manager.AddClientToList(client)
			conn.Write([]byte("Welcome to the Chat!\n"))
			go client.HandleConnection(quitTrigger)
			return nil
		}
	}
}
