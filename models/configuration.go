package models

type ServerConfiguration struct {
	ListenPort string `json:"listenPort"`
	BufferSize int    `json:"bufferSize"`
	MaxClients int    `json:"maxClients"`
}

type ClientConfiguration struct {
	DefaultChatHost     string `json:"defaultChatHost"`
	DefaultChatHostPort string `json:"defaultChatHostPort"`
	DefaultProtocol     string `json:"defaultProtocol"`
	BufferSize          int    `json:"bufferSize"`
}
