package client

import (
	"bufio"
	"fmt"
	"io"
	"sync"
)

type TerminalInput struct {
	Reader io.Reader
	Writer io.Writer
}

var consoleMutex sync.Mutex

func (handler *TerminalInput) ReadMessage() (string, error) {
	consoleMutex.Lock()
	defer consoleMutex.Unlock()
	reader := bufio.NewReader(handler.Reader)
	data, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println("Error while reading: ", err.Error())
		return "", err
	}
	return string(data), nil
}

func (handler *TerminalInput) DisplayMessage(message string) error {
	consoleMutex.Lock()
	defer consoleMutex.Unlock()
	writer := bufio.NewWriter(handler.Writer)
	_, err := writer.WriteString(message + "\n")
	if err != nil {
		fmt.Println("Error while writing message to output: ", err.Error())
		return err
	}
	return nil
}
