package models

import (
	"net"
	"time"
)

type Message struct {
	ClientName  string
	MessageText string
	Time        time.Time
}

const (
	START int = iota
	END
	PAUSE
	ERROR
)

const (
	END_CHAT   string = "END\n"
	PAUSE_CHAT string = "PAUSE\n"
)

type Signal struct {
	SignalType int
	ClientName string
}

type Node struct {
	Connection     *net.Conn
	Name           string
	SendChannel    chan Message
	ReceiveChannel chan string
	SignalChannel  chan int
	ServerChannel  chan<- Message
}

type Client struct {
	Connection     *net.Conn
	SendChannel    chan []byte
	ReceiveChannel chan string
}

type InputOutputHandler interface {
	DisplayMessage(string) error
	ReadMessage() ([]byte, error)
}
