package main

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	errReject = errors.New("password rejected")
)

type pwdCallback func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error)

func pwdCallbackFactory(ch chan<- string) pwdCallback {
	return func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
		entry := fmt.Sprintf("%s User: %s Pwd: %s EOL\n", time.Now(), c.User(), string(pass))
		ch <- entry
		return nil, errReject
	}
}

func main() {

}
