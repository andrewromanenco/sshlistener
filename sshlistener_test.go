package main

import (
	"net"
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
		if !strings.HasSuffix(entry, "User: login Pwd: password EOL\n") {
			t.Error("Log entry does not look like it should")
		}
	default:
		t.Error("Entry message was not sent to the channel")
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
	return nil
}

func (stub *stubConnMeta) LocalAddr() net.Addr {
	return nil
}
