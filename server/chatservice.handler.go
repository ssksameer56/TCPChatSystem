package server

import (
	"encoding/json"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/ssksameer56/TCPChatSystem/models"
)

type ServerConfig struct {
	config  models.ServerConfiguration
	manager ClientsManager
}

var server ServerConfig

func InitServer() error {
	file, _ := os.Open("server.settings.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	err := decoder.Decode(&server.config)
	if err != nil {
		log.WithFields(log.Fields{"client": "Server"}).Error("Cant get config:", err.Error())
		return err
	}
	server.manager = ClientsManager{}
	return nil
}

func ListenForClients() {
	log.WithFields(log.Fields{"client": "Server"}).Info("Starting Server")
	listener, err := net.Listen("tcp", server.config.ListenPort)
	if err != nil {
		log.WithFields(log.Fields{"client": "Server"}).Error("Listener: Listen Error", err)
		os.Exit(1)
	}
	defer listener.Close()

	log.WithFields(log.Fields{"client": "Server"}).Info("Listener: Listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.WithFields(log.Fields{"client": "Server"}).Error("Listener: Listen Error", err)
			continue
		}
		if len(server.manager.AllClients) >= server.config.MaxClients {
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
	}
	return nil
}

func RunServer() {
	InitServer()
	go ListenForClients()
}
