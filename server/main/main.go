package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/ssksameer56/TCPChatSystem/server"
)

var logFile string

func InitializeLogging(logFile string) {
	var file, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Could Not Open Log File : " + err.Error())
	}
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetOutput(file)
}

func main() {
	logFile = "./log"
	InitializeLogging(logFile)
	server.RunServer()
}
