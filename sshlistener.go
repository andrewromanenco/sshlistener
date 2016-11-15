package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

const (
	concurrentHandlers = 50
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

func buildSSHConfig(logChannel chan<- string, privateKeyFile string) *ssh.ServerConfig {
	config := &ssh.ServerConfig{
		PasswordCallback: pwdCallbackFactory(logChannel),
	}

	privateBytes, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		panic(err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic(err)
	}

	config.AddHostKey(private)
	return config
}

func readLoginInfo(nConn net.Conn, config *ssh.ServerConfig, semaphore <-chan int) {
	defer nConn.Close()
	defer func() { <-semaphore }()
	ssh.NewServerConn(nConn, config)
}

func runLogServer(privateKeyFile, logFile string) {
	semaphore := make(chan int, concurrentHandlers)
	logChannel := make(chan string)
	config := buildSSHConfig(logChannel, privateKeyFile)

	listener, err := net.Listen("tcp", "0.0.0.0:2022")
	if err != nil {
		panic(err)
	}

	go writeToFile(logChannel, logFile)

	for {
		semaphore <- 1
		nConn, err := listener.Accept()
		go readLoginInfo(nConn, config, semaphore)
		if err != nil {
			log.Println("failed to accept incoming connection: ", err)
		}
	}
}

func main() {
	var privateKeyFile string
	var logFile string
	flag.StringVar(&privateKeyFile, "private", "", "Path to private key (id_rsa)")
	flag.StringVar(&logFile, "output", "", "Path to log file to write to")
	flag.Parse()
	if privateKeyFile == "" || logFile == "" {
		flag.PrintDefaults()
		return
	}

	runLogServer(privateKeyFile, logFile)
}
