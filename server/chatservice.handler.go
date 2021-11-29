package server

import (
	"net"
)

type ServerConfig struct {
	BufferSize int
	Manager    ClientsManager
}

var server ServerConfig

func init() {
	server.BufferSize = 10
	server.Manager = ClientsManager{}
}

func ListenForClients() {

}

func HandoverToManager(conn net.Conn) error {
	err := server.Manager.CreateClient(conn)
	if err != nil {
		return err
	}
	return nil
}
