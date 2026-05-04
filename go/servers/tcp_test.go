package servers_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
	"github.com/goldenm-software/layrz-protocol/go/v3/servers"
)

// --- Constructor validation ---

func TestNew_NilConfigMissingHandler(t *testing.T) {
	_, err := servers.New(nil)
	if err == nil {
		t.Error("expected error for nil config (missing OnNewPacket)")
	}
}

func TestNew_MissingOnNewPacket(t *testing.T) {
	_, err := servers.New(&servers.TcpConfig{Port: 9000})
	if err == nil {
		t.Error("expected error for missing OnNewPacket")
	}
}

func TestNew_InvalidPort_Zero(t *testing.T) {
	_, err := servers.New(&servers.TcpConfig{
		Port:        0,
		OnNewPacket: func(client.ClientPackets, net.Conn) (server.ServerPackets, error) { return nil, nil },
	})
	if err == nil {
		t.Error("expected error for port=0")
	}
}

func TestNew_InvalidPort_MaxBound(t *testing.T) {
	_, err := servers.New(&servers.TcpConfig{
		Port:        65535,
		OnNewPacket: func(client.ClientPackets, net.Conn) (server.ServerPackets, error) { return nil, nil },
	})
	if err == nil {
		t.Error("expected error for port=65535")
	}
}

func TestNew_DefaultsOnDecodeError(t *testing.T) {
	srv, err := servers.New(&servers.TcpConfig{
		Port:        19000,
		OnNewPacket: func(client.ClientPackets, net.Conn) (server.ServerPackets, error) { return nil, nil },
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if srv == nil {
		t.Fatal("expected non-nil server")
	}
}

// --- Start / Close integration test ---

func TestTcpServer_StartAndClose(t *testing.T) {
	// Find a free port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()

	srv, err := servers.New(&servers.TcpConfig{
		Port:        port,
		OnNewPacket: func(client.ClientPackets, net.Conn) (server.ServerPackets, error) { return nil, nil },
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { _ = srv.Start(ctx) }()
	time.Sleep(50 * time.Millisecond)

	// Close should return without error; Start keeps blocking on Accept until a
	// connection arrives, so we don't wait for it here.
	if err := srv.Close(); err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

// --- handleConnection behaviour tests using a real bound server ---

func startTcpServer(t *testing.T, cfg *servers.TcpConfig) (port int, cancelFn func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("find free port: %v", err)
	}
	port = ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()

	cfg.Port = port
	srv, err := servers.New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() { _ = srv.Start(ctx) }()
	time.Sleep(30 * time.Millisecond) // let the listener bind

	return port, func() {
		cancel()
		_ = srv.Close()
	}
}

func TestTcpServer_ValidPacket_NilResponse(t *testing.T) {
	called := make(chan struct{}, 1)
	port, cancel := startTcpServer(t, &servers.TcpConfig{
		OnNewPacket: func(p client.ClientPackets, conn net.Conn) (server.ServerPackets, error) {
			called <- struct{}{}
			return nil, nil
		},
	})
	defer cancel()

	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	msg := *(&client.PrPacket{}).ToPacket() + "\n"
	if _, err := fmt.Fprint(conn, msg); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case <-called:
	case <-time.After(2 * time.Second):
		t.Error("OnNewPacket was not called")
	}
}

func TestTcpServer_ValidPacket_WithResponse(t *testing.T) {
	asPacket := &server.AsPacket{}
	port, cancel := startTcpServer(t, &servers.TcpConfig{
		OnNewPacket: func(p client.ClientPackets, conn net.Conn) (server.ServerPackets, error) {
			return asPacket, nil
		},
	})
	defer cancel()

	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	msg := *(&client.PrPacket{}).ToPacket() + "\n"
	if _, err := fmt.Fprint(conn, msg); err != nil {
		t.Fatalf("write: %v", err)
	}

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	got := string(buf[:n])
	if got != *asPacket.ToPacket() {
		t.Errorf("response mismatch: got %q, want %q", got, *asPacket.ToPacket())
	}
}

func TestTcpServer_GarbagePacket_DecodeError(t *testing.T) {
	decodeErrCalled := make(chan struct{}, 1)
	port, cancel := startTcpServer(t, &servers.TcpConfig{
		OnNewPacket: func(p client.ClientPackets, conn net.Conn) (server.ServerPackets, error) {
			return nil, nil
		},
		OnDecodeError: func(err error, data []byte, conn net.Conn) {
			decodeErrCalled <- struct{}{}
		},
	})
	defer cancel()

	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	if _, err := fmt.Fprint(conn, "garbage data\n"); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case <-decodeErrCalled:
	case <-time.After(2 * time.Second):
		t.Error("OnDecodeError was not called")
	}
}

func TestTcpServer_MultipleConcatenatedPackets(t *testing.T) {
	callCount := 0
	done := make(chan struct{})
	port, cancel := startTcpServer(t, &servers.TcpConfig{
		OnNewPacket: func(p client.ClientPackets, conn net.Conn) (server.ServerPackets, error) {
			callCount++
			if callCount == 2 {
				close(done)
			}
			return nil, nil
		},
	})
	defer cancel()

	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	// Two packets concatenated, terminated with \n
	pr := *(&client.PrPacket{}).ToPacket()
	msg := pr + pr + "\n"
	if _, err := fmt.Fprint(conn, msg); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Errorf("expected 2 callback invocations, got %d", callCount)
	}
}

func TestTcpServer_HandlerError_NoPanic(t *testing.T) {
	// When the callback returns an error, the server should log and continue — no panic
	called := make(chan struct{}, 1)
	port, cancel := startTcpServer(t, &servers.TcpConfig{
		OnNewPacket: func(p client.ClientPackets, conn net.Conn) (server.ServerPackets, error) {
			called <- struct{}{}
			return nil, fmt.Errorf("handler error")
		},
	})
	defer cancel()

	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	msg := *(&client.PrPacket{}).ToPacket() + "\n"
	if _, err := fmt.Fprint(conn, msg); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case <-called:
	case <-time.After(2 * time.Second):
		t.Error("OnNewPacket was not called")
	}
}
