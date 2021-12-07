package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/ssksameer56/TCPChatSystem/models"
)

type ServerConfig struct {
	Config  models.ServerConfiguration
	manager ClientsManager
}

var server ServerConfig

func InitServer() error {
	fmt.Println(os.Getwd())
	file, _ := os.Open("./server.settings.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	err := decoder.Decode(&server.Config)
	if err != nil {
		log.WithFields(log.Fields{"client": "Server"}).Error("Cant get config:", err.Error())
		return err
	}
	server.manager = ClientsManager{}
	return nil
}

func ListenForClients(wg sync.WaitGroup) {
	log.WithFields(log.Fields{"client": "Server"}).Info("Starting Server")
	listener, err := net.Listen("tcp", "localhost:"+server.Config.ListenPort)
	if err != nil {
		log.WithFields(log.Fields{"client": "Server"}).Error("Listener: Listen Error", err)
		os.Exit(1)
	}
	defer listener.Close()

	go server.manager.HandleClients() //Start the loop to handle connected clients

	log.WithFields(log.Fields{"client": "Server"}).Info("Listener: Listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.WithFields(log.Fields{"client": "Server"}).Error("Listener: Listen Error", err)
			continue
		}
		if len(server.manager.AllClients) >= server.Config.MaxClients {
			_, err := conn.Write([]byte("Max Clients Reached. Please try again"))
			if err != nil {
				log.WithFields(log.Fields{"client": "Server"}).Error("Error while closing connection", err)
			}
			err = conn.Close()
			if err != nil {
				log.WithFields(log.Fields{"client": "Server"}).Error("Error while closing connection", err)
			}
			continue
		}
		go HandoverToManager(conn)
	}
}

func HandoverToManager(conn net.Conn) error {
	err := server.manager.CreateClient(conn)
	if err != nil {
		log.WithFields(log.Fields{"client": "Server"}).Error(err.Error())
		conn.Write([]byte("Cant create a client\n"))
		conn.Close()
		return err
	}
	return nil
}

func RunServer() {
	err := InitServer()
	if err != nil {
		log.WithFields(log.Fields{"client": "Server"}).Error("Cant start:", err.Error())
		os.Exit(1)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go ListenForClients(wg)
	wg.Wait()
}
