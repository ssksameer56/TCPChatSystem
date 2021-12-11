package server

import (
	"bufio"
	"fmt"
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
	log.WithFields(log.Fields{"client": client.Name}).Info("Got message: " + text)
	return Message{
		ClientName:  client.Name,
		MessageText: text,
		Time:        time.Now(),
	}
}

//Read Data from Client
func (client *Client) Read() {
	log.WithFields(log.Fields{"client": client.Name}).Info("Started reading data from this client")
	reader := bufio.NewReader(*client.Connection)
	for {
		var data string
		data, err := reader.ReadString('\n')
		if err != nil {
			log.WithFields(log.Fields{"client": client.Name}).Error(err.Error())
			client.SignalChannel <- models.END
			return
		}
		message := data
		if message == models.END_CHAT {
			select {
			case client.SignalChannel <- models.END:
			default:
				log.WithFields(log.Fields{"client": client.Name}).Info("Buffer full. discarding signal: " + message)
			}
			return
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
func (client *Client) HandleConnection(quitTrigger chan bool) {
	log.WithFields(log.Fields{"client": client.Name}).Info("Started connection loop")
	go client.Read()
	for {
		select {
		case data, ok := <-client.ReceiveChannel:
			if !ok {
				log.WithFields(log.Fields{"client": client.Name}).Error("Error when reading message from client")
				break
			}
			message := client.ReceiveMessage(data)
			select {
			case client.ServerChannel <- message:
			default:
				log.WithFields(log.Fields{"client": client.Name}).Info("Server Buffer full. discarding message: " + message.MessageText)
			}
		case data, ok := <-client.SendChannel:
			if !ok {
				log.WithFields(log.Fields{"client": client.Name}).Error("Error when sending message to client")
				break
			}
			client.SendMessage(data)
		case signal, ok := <-client.SignalChannel:
			if !ok {
				log.WithFields(log.Fields{"client": client.Name}).Error("Error when reading singal from client")
				break
			}
			if signal == models.END {
				close(client.ReceiveChannel)
				close(client.SignalChannel)
				(*client.Connection).Close()
				msg := models.Message{MessageText: models.END_CHAT, ClientName: client.Name, Time: time.Now()}
				client.ServerChannel <- msg
				log.WithFields(log.Fields{"client": client.Name}).Info("Closed Connection")
				return
			}
		case <-quitTrigger:
			log.WithFields(log.Fields{"client": client.Name}).Info("Close Connection trigger from server received")
			return
		default:
			continue
		}
	}
}

//Generate a new client from the TCP Connection
func NewClient(name string, conn *net.Conn, buffSize int, serverChannel chan<- Message) *Client {
	client := Client{
		Name:           name,                         //The name for client, for example username
		Connection:     conn,                         //The connection object to the TCP Connection
		SendChannel:    make(chan Message, buffSize), //Channel to handle data to be sent to client
		ReceiveChannel: make(chan string, buffSize),  //Channel to recieve data from client
		SignalChannel:  make(chan int, 10),           //Channel to signal status changes to server
		ServerChannel:  serverChannel,                //Channel to send to server
	}
	return &client
}
