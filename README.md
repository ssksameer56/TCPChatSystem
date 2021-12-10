# TCPChatSystem
A Simple TCP Chat System Written in Golang

# Features
1. Supports multiple clients
2. Uses basic TCP/IP Protocol to enable messaging
3. You can start the server in one terminal and start respective clients in new terminals as required.

# Usage

The Package contains two modules - server and client.



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

# Built Using

1. Go
2. Logrus for logging
