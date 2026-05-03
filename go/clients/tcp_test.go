package clients

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
)

func TestTcpComm_New(t *testing.T) {
	var c TcpComm
	c.New("localhost", 5000, "ident", "pass")

	if c.Host != "localhost" {
		t.Errorf("expected Host=localhost, got %s", c.Host)
	}
	if c.Port != 5000 {
		t.Errorf("expected Port=5000, got %d", c.Port)
	}
	if c.Ident != "ident" {
		t.Errorf("expected Ident=ident, got %s", c.Ident)
	}
	if c.Passwd != "pass" {
		t.Errorf("expected Passwd=pass, got %s", c.Passwd)
	}
	if !c.initialized {
		t.Error("expected initialized=true")
	}
	if c.splitRegExp == nil {
		t.Error("expected splitRegExp to be set")
	}
	if c.endRegExp == nil {
		t.Error("expected endRegExp to be set")
	}
}

func TestTcpComm_NotInitialized(t *testing.T) {
	var c TcpComm

	if err := c.SetCallback(func(*any) {}); err == nil {
		t.Error("SetCallback: expected error when not initialized")
	}
	if err := c.Send(&client.PrPacket{}); err == nil {
		t.Error("Send: expected error when not initialized")
	}
	if err := c.Close(); err == nil {
		t.Error("Close: expected error when not initialized")
	}
}

func TestTcpComm_SetCallback(t *testing.T) {
	var c TcpComm
	c.New("localhost", 5000, "ident", "pass")

	called := false
	if err := c.SetCallback(func(*any) { called = true }); err != nil {
		t.Fatalf("SetCallback failed: %v", err)
	}
	if c.callback == nil {
		t.Error("expected callback to be set")
	}
	(*c.callback)(nil)
	if !called {
		t.Error("callback was not called")
	}
}

func TestTcpComm_Send_HappyPath(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer func() { _ = serverConn.Close() }()

	var c TcpComm
	c.New("localhost", 5000, "ident", "pass")
	c.conn = &clientConn

	done := make(chan []byte, 1)
	go func() {
		buf := make([]byte, 4096)
		n, _ := serverConn.Read(buf)
		done <- buf[:n]
	}()

	if err := c.Send(&client.PrPacket{}); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	select {
	case data := <-done:
		if !strings.HasSuffix(strings.TrimRight(string(data), "\r\n"), "</Pr>") {
			t.Errorf("unexpected payload: %q", data)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for sent data")
	}
}

func TestTcpComm_Close_HappyPath(t *testing.T) {
	clientConn, _ := net.Pipe()

	var c TcpComm
	c.New("localhost", 5000, "ident", "pass")
	c.conn = &clientConn

	if err := c.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func TestTcpComm_Listen_AuthenticatesOnAsPacket(t *testing.T) {
	asPacket := server.AsPacket{}
	encoded := *asPacket.ToPacket() + "\r\n"

	clientConn, serverConn := net.Pipe()

	var c TcpComm
	c.New("localhost", 5000, "ident", "pass")
	c.conn = &clientConn

	// listen() panics on read error; wrap in goroutine with recover so test doesn't crash
	go func() {
		defer func() { _ = recover() }()
		c.listen()
	}()

	if _, err := fmt.Fprint(serverConn, encoded); err != nil {
		_ = serverConn.Close()
		t.Fatalf("failed to write to pipe: %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if c.authenticated {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	_ = serverConn.Close()
	_ = clientConn.Close()

	if !c.authenticated {
		t.Error("expected authenticated=true after receiving AsPacket")
	}
}

func TestTcpComm_Listen_CallsCallbackForOtherPackets(t *testing.T) {
	// Use a server packet other than As/Au so the callback branch fires
	aoPacket := server.AoPacket{Timestamp: time.Unix(1700000000, 0)}
	encoded := *aoPacket.ToPacket() + "\r\n"

	clientConn, serverConn := net.Pipe()

	var c TcpComm
	c.New("localhost", 5000, "ident", "pass")
	c.conn = &clientConn

	called := make(chan struct{}, 1)
	_ = c.SetCallback(func(*any) {
		called <- struct{}{}
	})

	go func() {
		defer func() { _ = recover() }()
		c.listen()
	}()

	_, _ = fmt.Fprint(serverConn, encoded)

	select {
	case <-called:
	case <-time.After(2 * time.Second):
		_ = serverConn.Close()
		_ = clientConn.Close()
		t.Error("callback was not called")
		return
	}

	_ = serverConn.Close()
	_ = clientConn.Close()
}

func TestTcpComm_Connect(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}

	asPacket := server.AsPacket{}
	asEncoded := *asPacket.ToPacket() + "\r\n"

	// Hold the server-side connection open for the lifetime of the test process
	// so that listen()'s Read never returns an error and never panics.
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		// intentionally never Close — connection stays alive while the test binary runs
		buf := make([]byte, 4096)
		_, _ = conn.Read(buf)
		_, _ = fmt.Fprint(conn, asEncoded)
		// block forever so the client-side Read in listen() never gets EOF
		select {}
	}()

	addr := ln.Addr().(*net.TCPAddr)
	var c TcpComm
	c.New("127.0.0.1", addr.Port, "ident", "pass")

	if err := c.Connect(); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	if !c.authenticated {
		t.Error("expected authenticated=true after Connect")
	}
}

func TestTcpComm_Connect_NotInitialized(t *testing.T) {
	var c TcpComm
	if err := c.Connect(); err == nil {
		t.Error("expected error when not initialized")
	}
}

func TestTcpComm_Connect_DialError(t *testing.T) {
	var c TcpComm
	c.New("127.0.0.1", 1, "ident", "pass")
	if err := c.Connect(); err == nil {
		t.Error("expected dial error")
	}
}
