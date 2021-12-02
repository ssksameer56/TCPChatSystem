package client

import (
	"bufio"
	"fmt"
	"os"
)

type TerminalInput struct {
}

func (handler *TerminalInput) ReadMessage(clientName string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	data, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println("Error while reading: ", err.Error())
		return "", err
	}
	return string(data), nil
}

func (handler *TerminalInput) DisplayMessage(message string) error {
	writer := bufio.NewWriter(os.Stdout)
	_, err := writer.WriteString(message)
	if err != nil {
		fmt.Println("Error while writing message to output: ", err.Error())
		return err
	}
	return nil
}
