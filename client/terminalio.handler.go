package client

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
)

type TerminalInput struct {
	Reader io.Reader
	Writer io.Writer
}

var consoleMutex sync.Mutex

func (handler *TerminalInput) ReadMessage() (string, error) {
	reader := bufio.NewReader(handler.Reader)
	for {
		check, _ := reader.ReadByte()
		if strings.ToLower(string(check)) == "c" {
			break
		}
	}
	consoleMutex.Lock()
	fmt.Println("Enter Message:")
	fmt.Print(">")
	data, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println("Error while reading: ", err.Error())
		return "", err
	}
	consoleMutex.Unlock()
	return string(data), nil
}

func (handler *TerminalInput) DisplayMessage(message string) error {
	consoleMutex.Lock()
	defer consoleMutex.Unlock()
	writer := bufio.NewWriter(handler.Writer)
	_, err := writer.WriteString(">> " + message + "\n")
	if err != nil {
		fmt.Println("Error while writing message to output: ", err.Error())
		return err
	}
	return nil
}
