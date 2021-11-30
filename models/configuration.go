package models

type ServerConfiguration struct {
	ListenPort string `json:"listenPort"`
	BufferSize int    `json:"bufferSize"`
	MaxClients int    `json:"maxClients"`
}
