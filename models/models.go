package models

import "net"

type Message struct {
	ClientName  string
	MessageText string
}

const (
	START int = iota
	END
	PAUSE
	ERROR
)

type Signal struct {
	SignalType int
	Message    string
}

type Node struct {
	Connection     *net.Conn
	Name           string
	SendChannel    chan Message
	ReceiveChannel chan Message
}
