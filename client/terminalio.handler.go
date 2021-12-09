package client

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

type TerminalInput struct {
	Reader io.Reader
	Writer io.Writer
}

var consoleMutex sync.Mutex

func (handler *TerminalInput) ReadMessage() (string, error) {
	reader := bufio.NewReader(handler.Reader)
	fmt.Print(">")
	consoleMutex.Lock()
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
	writer := bufio.NewWriter(os.Stdout)
	_, err := writer.WriteString(">> " + message + "\n")
	writer.Flush()
	if err != nil {
		fmt.Println("Error while writing message to output: ", err.Error())
		return err
	}
	consoleMutex.Unlock()
	return nil
}
