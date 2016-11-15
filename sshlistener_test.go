package main

import (
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"
)

func TestPwdCallbackSendsDataToChannel(t *testing.T) {
	stubConn := &stubConnMeta{}
	channel := make(chan string, 1)
	testee := pwdCallbackFactory(channel)
	testee(stubConn, []byte("password"))
	select {
	case entry := <-channel:
		if entry != "Ip: 192.168.1.1 User: login Pwd: password EOL" {
			t.Error("Log entry does not look like it should")
		}
	default:
		t.Error("Entry message was not sent to the channel")
	}
}

func TestWriteToFile(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "")
	os.Remove(file.Name())
	defer os.Remove(file.Name())
	channel := make(chan string, 2)
	channel <- "entry1"
	channel <- "entry2"
	close(channel)
	writeToFile(channel, file.Name())
	data, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Error("Can't read file:", err)
	}
	sData := string(data)
	if !strings.Contains(sData, "entry1") ||
		!strings.Contains(sData, "entry2") {
		t.Error("Data is not logged correctly")
	}
}

type stubConnMeta struct {
}

func (stub *stubConnMeta) User() string {
	return "login"
}

func (stub *stubConnMeta) SessionID() []byte {
	return nil
}

func (stub *stubConnMeta) ClientVersion() []byte {
	return nil
}

func (stub *stubConnMeta) ServerVersion() []byte {
	return nil
}

func (stub *stubConnMeta) RemoteAddr() net.Addr {

	return &Addr{}
}

func (stub *stubConnMeta) LocalAddr() net.Addr {
	return nil
}

type Addr struct{}

func (a *Addr) Network() string {
	return ""
}

func (a *Addr) String() string {
	return "192.168.1.1"
}
