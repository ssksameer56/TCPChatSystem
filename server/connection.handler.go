package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ssksameer56/TCPChatSystem/models"
)

type Client models.Node
type Message = models.Message

//Send Message to Client
func (client *Client) SendMessage(message Message) error {
	messageToSend := fmt.Sprintf("%s @ %s: %s", message.ClientName, message.Time.Format("15:04:05"), message.MessageText)
	if _, err := (*client.Connection).Write([]byte(messageToSend)); err != nil {
		log.WithFields(log.Fields{"client": client.Name}).Error(err.Error())
		return err
	}
	return nil
}

//Read and convert the message from Client
func (client *Client) ReceiveMessage(text string) Message {
	return Message{
		ClientName:  client.Name,
		MessageText: text,
		Time:        time.Now(),
	}
}

//Read Data from Client
func (client *Client) Read() {
	for {
		var data []byte
		if _, err := (*client.Connection).Read(data); err != nil {
			log.WithFields(log.Fields{"client": client.Name}).Error(err.Error())
			continue
		}
		reader := bufio.NewReader(*client.Connection)
		data, _, err := reader.ReadLine()
		if err == io.EOF {
			client.SignalChannel <- models.END
			log.WithFields(log.Fields{"client": client.Name}).Info("Closing Connection")
		} else if err != nil {
			log.WithFields(log.Fields{"client": client.Name}).Error(err.Error())
			continue
		}
		message := string(data)
		if message == models.END_CHAT {
			select {
			case client.SignalChannel <- models.END:
			default:
				log.WithFields(log.Fields{"client": client.Name}).Info("Buffer full. discarding signal: " + message)
			}
		} else {
			select {
			case client.ReceiveChannel <- message:
			default:
				log.WithFields(log.Fields{"client": client.Name}).Info("Buffer full. discarding message: " + message)
			}
		}
	}
}

//Handle all the communication between client and server
func (client *Client) HandleConnection() {
	go client.Read()
	for {
		select {
		case data := <-client.ReceiveChannel:
			client.ReceiveMessage(data)
		case data := <-client.SendChannel:
			client.SendMessage(data)
		case signal := <-client.SignalChannel:
			if signal == models.END {
				close(client.ReceiveChannel)
				close(client.SignalChannel)
				(*client.Connection).Close()
				log.WithFields(log.Fields{"client": client.Name}).Info("Closing Connection")
			}
		default:
			continue
		}
	}
}

//Generate a new client from the TCP Connection
func NewClient(name string, conn *net.Conn, buffSize int) *Client {
	client := Client{
		Name:           name,                         //The name for client, for example username
		Connection:     conn,                         //The connection object to the TCP Connection
		SendChannel:    make(chan Message, buffSize), //Channel to handle data to be sent to client
		ReceiveChannel: make(chan string, buffSize),  //Channel to recieve data from client
		SignalChannel:  make(chan int, 1),            //Channel to signal status changes to server
	}
	return &client
}
