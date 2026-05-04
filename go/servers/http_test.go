package servers_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
	"github.com/goldenm-software/layrz-protocol/go/v3/servers"
)

// realHttpServer starts the HttpServer on a free port and returns its base URL and a stop func.
func realHttpServer(t *testing.T, cfg *servers.HttpConfig) (baseURL string, stop func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("find free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()

	cfg.Port = port
	srv, err := servers.NewHttp(cfg)
	if err != nil {
		t.Fatalf("NewHttp: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() { _ = srv.Start(ctx) }()
	time.Sleep(30 * time.Millisecond)

	return fmt.Sprintf("http://127.0.0.1:%d", port), cancel
}

// --- Constructor validation ---

func TestNewHttp_NilConfigMissingHandler(t *testing.T) {
	_, err := servers.NewHttp(nil)
	if err == nil {
		t.Error("expected error for nil config (missing OnNewPacket)")
	}
}

func TestNewHttp_MissingOnNewPacket(t *testing.T) {
	_, err := servers.NewHttp(&servers.HttpConfig{Port: 8080})
	if err == nil {
		t.Error("expected error for missing OnNewPacket")
	}
}

func TestNewHttp_InvalidPort_Zero(t *testing.T) {
	_, err := servers.NewHttp(&servers.HttpConfig{
		Port:        0,
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
	})
	if err == nil {
		t.Error("expected error for port=0")
	}
}

func TestNewHttp_InvalidPort_MaxBound(t *testing.T) {
	_, err := servers.NewHttp(&servers.HttpConfig{
		Port:        65535,
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
	})
	if err == nil {
		t.Error("expected error for port=65535")
	}
}

func TestNewHttp_DefaultsOnDecodeError(t *testing.T) {
	srv, err := servers.NewHttp(&servers.HttpConfig{
		Port:        18080,
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
	})
	if err != nil {
		t.Fatalf("NewHttp failed: %v", err)
	}
	if srv == nil {
		t.Fatal("expected non-nil server")
	}
}

// --- Close without Start ---

func TestHttpServer_Close_WithoutStart(t *testing.T) {
	srv, err := servers.NewHttp(&servers.HttpConfig{
		Port:        18080,
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
	})
	if err != nil {
		t.Fatalf("NewHttp failed: %v", err)
	}
	if err := srv.Close(); err != nil {
		t.Errorf("Close without Start should not error: %v", err)
	}
}

// --- Start / Close integration ---

func TestHttpServer_StartClose(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()

	srv, err := servers.NewHttp(&servers.HttpConfig{
		Port:        port,
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
	})
	if err != nil {
		t.Fatalf("NewHttp failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- srv.Start(ctx) }()
	time.Sleep(50 * time.Millisecond)

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Start returned error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Error("timeout waiting for server to stop")
	}
}

// --- handleMessage tests ---

func TestHandleMessage_MethodNotAllowed(t *testing.T) {
	url, stop := realHttpServer(t, &servers.HttpConfig{
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
	})
	defer stop()

	resp, err := http.Get(url + "/v2/message")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestHandleMessage_MissingAuth(t *testing.T) {
	url, stop := realHttpServer(t, &servers.HttpConfig{
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
	})
	defer stop()

	req, _ := http.NewRequest(http.MethodPost, url+"/v2/message", bytes.NewBufferString("body"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestHandleMessage_AuthCallbackRejects(t *testing.T) {
	url, stop := realHttpServer(t, &servers.HttpConfig{
		OnNewPacket:    func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
		OnAuthenticate: func(ident, passwd string, r *http.Request) bool { return false },
	})
	defer stop()

	req, _ := http.NewRequest(http.MethodPost, url+"/v2/message", bytes.NewBufferString("body"))
	req.Header.Set("Authorization", "LayrzAuth ident;pass")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestHandleMessage_DecodeError(t *testing.T) {
	decodeErrCalled := make(chan struct{}, 1)
	url, stop := realHttpServer(t, &servers.HttpConfig{
		OnNewPacket:   func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
		OnDecodeError: func(e error, data []byte, r *http.Request) { decodeErrCalled <- struct{}{} },
	})
	defer stop()

	req, _ := http.NewRequest(http.MethodPost, url+"/v2/message", bytes.NewBufferString("garbage"))
	req.Header.Set("Authorization", "LayrzAuth ident;pass")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
	select {
	case <-decodeErrCalled:
	case <-time.After(2 * time.Second):
		t.Error("expected OnDecodeError to be called")
	}
}

func TestHandleMessage_HandlerError(t *testing.T) {
	url, stop := realHttpServer(t, &servers.HttpConfig{
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) {
			return nil, fmt.Errorf("handler error")
		},
	})
	defer stop()

	body := *(&client.PrPacket{}).ToPacket()
	req, _ := http.NewRequest(http.MethodPost, url+"/v2/message", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "LayrzAuth ident;pass")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}

func TestHandleMessage_NilResponse(t *testing.T) {
	url, stop := realHttpServer(t, &servers.HttpConfig{
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
	})
	defer stop()

	body := *(&client.PrPacket{}).ToPacket()
	req, _ := http.NewRequest(http.MethodPost, url+"/v2/message", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "LayrzAuth ident;pass")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}
}

func TestHandleMessage_WithResponse(t *testing.T) {
	asPacket := &server.AsPacket{}
	url, stop := realHttpServer(t, &servers.HttpConfig{
		OnNewPacket: func(p client.ClientPackets, r *http.Request) (server.ServerPackets, error) {
			return asPacket, nil
		},
	})
	defer stop()

	body := *(&client.PrPacket{}).ToPacket()
	req, _ := http.NewRequest(http.MethodPost, url+"/v2/message", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "LayrzAuth ident;pass")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	respBody, _ := io.ReadAll(resp.Body)
	if strings.TrimSpace(string(respBody)) != *asPacket.ToPacket() {
		t.Errorf("body mismatch: got %q, want %q", string(respBody), *asPacket.ToPacket())
	}
}

// --- handleCommands tests ---

func TestHandleCommands_MethodNotAllowed(t *testing.T) {
	url, stop := realHttpServer(t, &servers.HttpConfig{
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
	})
	defer stop()

	req, _ := http.NewRequest(http.MethodPost, url+"/v2/commands", nil)
	req.Header.Set("Authorization", "LayrzAuth ident;pass")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestHandleCommands_MissingAuth(t *testing.T) {
	url, stop := realHttpServer(t, &servers.HttpConfig{
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
	})
	defer stop()

	resp, err := http.Get(url + "/v2/commands")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestHandleCommands_AuthRejects(t *testing.T) {
	url, stop := realHttpServer(t, &servers.HttpConfig{
		OnNewPacket:    func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
		OnAuthenticate: func(ident, passwd string, r *http.Request) bool { return false },
	})
	defer stop()

	req, _ := http.NewRequest(http.MethodGet, url+"/v2/commands", nil)
	req.Header.Set("Authorization", "LayrzAuth ident;pass")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestHandleCommands_NilHandler(t *testing.T) {
	url, stop := realHttpServer(t, &servers.HttpConfig{
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
		// OnPullCommands intentionally nil
	})
	defer stop()

	req, _ := http.NewRequest(http.MethodGet, url+"/v2/commands", nil)
	req.Header.Set("Authorization", "LayrzAuth ident;pass")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}
}

func TestHandleCommands_HandlerError(t *testing.T) {
	url, stop := realHttpServer(t, &servers.HttpConfig{
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
		OnPullCommands: func(ident, passwd string, r *http.Request) (server.ServerPackets, error) {
			return nil, fmt.Errorf("cmd error")
		},
	})
	defer stop()

	req, _ := http.NewRequest(http.MethodGet, url+"/v2/commands", nil)
	req.Header.Set("Authorization", "LayrzAuth ident;pass")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}

func TestHandleCommands_NilResponse(t *testing.T) {
	url, stop := realHttpServer(t, &servers.HttpConfig{
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
		OnPullCommands: func(ident, passwd string, r *http.Request) (server.ServerPackets, error) {
			return nil, nil
		},
	})
	defer stop()

	req, _ := http.NewRequest(http.MethodGet, url+"/v2/commands", nil)
	req.Header.Set("Authorization", "LayrzAuth ident;pass")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}
}

func TestHandleCommands_WithResponse(t *testing.T) {
	asPacket := &server.AsPacket{}
	url, stop := realHttpServer(t, &servers.HttpConfig{
		OnNewPacket: func(client.ClientPackets, *http.Request) (server.ServerPackets, error) { return nil, nil },
		OnPullCommands: func(ident, passwd string, r *http.Request) (server.ServerPackets, error) {
			return asPacket, nil
		},
	})
	defer stop()

	req, _ := http.NewRequest(http.MethodGet, url+"/v2/commands", nil)
	req.Header.Set("Authorization", "LayrzAuth ident;pass")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	respBody, _ := io.ReadAll(resp.Body)
	if strings.TrimSpace(string(respBody)) != *asPacket.ToPacket() {
		t.Errorf("body mismatch: got %q, want %q", string(respBody), *asPacket.ToPacket())
	}
}
