package models

type ServerConfiguration struct {
	ListenPort string `json:"listenPort"`
	BufferSize int    `json:"bufferSize"`
	MaxClients int    `json:"maxClients"`
}

type ClientConfiguration struct {
	ConnectPort string `json:"listenPort"`
	ClientName  string `json:"clientName"`
}
