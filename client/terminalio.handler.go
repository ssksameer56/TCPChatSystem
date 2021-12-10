package client

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type TerminalInput struct {
	Reader io.Reader
	Writer io.Writer
}

//var consoleMutex sync.Mutex

func (handler *TerminalInput) ReadMessage() ([]byte, error) {
	reader := bufio.NewReader(handler.Reader)
	data, _, err := reader.ReadLine()
	data = append(data, '\n')
	if err != nil {
		fmt.Println("Error while reading: ", err.Error())
		return []byte{}, err
	}
	return data, nil
}

func (handler *TerminalInput) DisplayMessage(message string) error {
	writer := bufio.NewWriter(os.Stdout)
	_, err := writer.WriteString(">> " + message)
	writer.Flush()
	fmt.Print(">>")
	if err != nil {
		fmt.Println("Error while writing message to output: ", err.Error())
		return err
	}
	return nil
}
