package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

const (
	flushEveryNItems = 100
)

var (
	errReject = errors.New("password rejected")
)

type pwdCallback func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error)

func pwdCallbackFactory(ch chan<- string) pwdCallback {
	return func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
		entry := fmt.Sprintf("User: %s Pwd: %s EOL", c.User(), string(pass))
		ch <- entry
		return nil, errReject
	}
}

func writeToFile(ch <-chan string, filePath string) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
	for {
		entry, isOpen := <-ch
		if !isOpen {
			break
		}
		log.Print(entry)
	}
}

func main() {

}
