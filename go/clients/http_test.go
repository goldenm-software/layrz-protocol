package clients_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/clients"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
)

func TestHttpComm_NotInitialized(t *testing.T) {
	var c clients.HttpComm

	if _, err := c.Send(&client.PrPacket{}); err == nil {
		t.Error("Send: expected error when not initialized")
	}
	if _, err := c.GetCommands(); err == nil {
		t.Error("GetCommands: expected error when not initialized")
	}
}

func TestHttpComm_Send(t *testing.T) {
	asPacket := server.AsPacket{}
	responseBody := *asPacket.ToPacket()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v2/message" {
			t.Errorf("expected /v2/message, got %s", r.URL.Path)
		}
		authHeader := r.Header.Get("Authorization")
		expected := fmt.Sprintf("LayrzAuth %s;%s", "test-ident", "test-pass")
		if authHeader != expected {
			t.Errorf("expected auth header %q, got %q", expected, authHeader)
		}
		_, _ = fmt.Fprint(w, responseBody)
	}))
	defer srv.Close()

	var c clients.HttpComm
	c.New(clients.HTTP, srv.Listener.Addr().String(), "test-ident", "test-pass")

	result, err := c.Send(&client.PrPacket{})
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if result == nil || *result == nil {
		t.Fatal("result is nil")
	}
}

func TestHttpComm_Send_TransportError(t *testing.T) {
	var c clients.HttpComm
	c.New(clients.HTTP, "127.0.0.1:1", "ident", "pass")

	_, err := c.Send(&client.PrPacket{})
	if err == nil {
		t.Error("expected transport error")
	}
}

func TestHttpComm_Send_DecodeError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "not-a-valid-packet")
	}))
	defer srv.Close()

	var c clients.HttpComm
	c.New(clients.HTTP, srv.Listener.Addr().String(), "ident", "pass")

	_, err := c.Send(&client.PrPacket{})
	if err == nil {
		t.Error("expected decode error for invalid body")
	}
}

func TestHttpComm_GetCommands(t *testing.T) {
	aoPacket := server.AoPacket{Timestamp: parserFixedTime}
	responseBody := *aoPacket.ToPacket()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v2/commands" {
			t.Errorf("expected /v2/commands, got %s", r.URL.Path)
		}
		_, _ = fmt.Fprint(w, responseBody)
	}))
	defer srv.Close()

	var c clients.HttpComm
	c.New(clients.HTTP, srv.Listener.Addr().String(), "ident", "pass")

	result, err := c.GetCommands()
	if err != nil {
		t.Fatalf("GetCommands failed: %v", err)
	}
	if result == nil || *result == nil {
		t.Fatal("result is nil")
	}
}

func TestHttpComm_GetCommands_TransportError(t *testing.T) {
	var c clients.HttpComm
	c.New(clients.HTTP, "127.0.0.1:1", "ident", "pass")

	_, err := c.GetCommands()
	if err == nil {
		t.Error("expected transport error")
	}
}
