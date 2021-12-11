# TCPChatSystem
A Simple TCP Chat System Written in Go. The purpose is to explore the functionality of goroutines, channels and concurrency in Go. I have also tried to use interfaces to implement certain features(Reader/Writer)

# Features
1. Supports multiple clients
2. Uses basic TCP/IP Protocol to enable messaging
3. You can start the server in one terminal and start respective clients in new terminals as required.
4. Server logs all actions by clients via logrus - configurable to be text or JSON friendly
5. Supports custom client name, client joining and exit as required with custom keywords

# I/O Handler
1. I/O Handler is implementd via interface `models.ioInputOutputHandler` which is a interface that supports any reader or writer interface to send messages to and read message from in client.
2. Currently the interface is implemented in a `terminalIOHandler`. Reading or writing from file can also be done by implementing the interface.

# Usage

The Go Module contains two packages - server and client respectively.

To start the server navigate the server folder and start the server

```
cd server/main
go run main.go
```


To start a client navigate the server folder and start the client

```
cd client/main
go run main.go
```

# Customising

You can modify the settings in Client and Server via the \*.settings.json file

| Setting        | Purpose      |
| ------------- |:-------------:|
| defaultChatHost      | Address for chat server |
| defaultChatHostPort      | Port on which server is listening      |
| bufferSize | Size of buffer to queue messages      |
| maxClients | Maximum number of clients that can connect     |


You can customize the logging by modifying the struct `log.TextFormatter`. 

Refer https://github.com/sirupsen/logrus for more details

```
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
```

# Built Using

- Go
- net/http Package
- Logrus for logging
- Encoding/JSON for loading config files
