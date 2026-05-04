package servers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
)

type HttpConfig struct {
	// Port to bind the HTTP listener on
	Port int

	// Called for every POST /v2/message.
	// Return a non-nil ServerPackets to write the encoded packet as the response body.
	// Return nil to respond with 204.
	// Return an error to respond with 500.
	OnNewPacket func(packet client.ClientPackets, r *http.Request) (server.ServerPackets, error)

	// Called for every GET /v2/commands.
	// ident and passwd are extracted from the LayrzAuth header.
	// Return nil to respond with 204; non-nil to write the encoded packet.
	OnPullCommands func(ident, passwd string, r *http.Request) (server.ServerPackets, error)

	// Called to authenticate every request.
	// If nil, all requests are allowed.
	OnAuthenticate func(ident, passwd string, r *http.Request) bool

	// Called when a packet cannot be decoded; parallel to TcpConfig.OnDecodeError.
	OnDecodeError func(err error, data []byte, r *http.Request)
}

type HttpServer struct {
	config *HttpConfig
	srv    *http.Server
}

// NewHttp creates a new HTTP server with the given configuration.
func NewHttp(cfg *HttpConfig) (*HttpServer, error) {
	if cfg == nil {
		cfg = &HttpConfig{}
	}

	if cfg.OnNewPacket == nil {
		return nil, fmt.Errorf("on new packet handler is not set")
	}

	if cfg.OnDecodeError == nil {
		cfg.OnDecodeError = func(err error, data []byte, r *http.Request) {
			log.Printf("Error decoding packet: %s Data: %s", err.Error(), string(data))
		}
	}

	if cfg.Port <= 0 || cfg.Port >= 65535 {
		return nil, fmt.Errorf("port is not valid")
	}

	return &HttpServer{config: cfg}, nil
}

// Start begins listening for HTTP requests.
// This method is blocking and runs until ctx is cancelled or an error occurs.
func (s *HttpServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/v2/message", s.handleMessage)
	mux.HandleFunc("/v2/commands", s.handleCommands)

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: mux,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		return s.Close()
	case err := <-errCh:
		return err
	}
}

// Close shuts down the HTTP server gracefully.
func (s *HttpServer) Close() error {
	if s.srv == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}

func (s *HttpServer) handleMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ident, passwd, ok := parseLayrzAuth(r.Header.Get("Authorization"))
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if s.config.OnAuthenticate != nil && !s.config.OnAuthenticate(ident, passwd, r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	packet, err := client.Decode(data)
	if err != nil {
		s.config.OnDecodeError(err, data, r)
		http.Error(w, "invalid packet", http.StatusBadRequest)
		return
	}

	response, err := s.config.OnNewPacket(packet, r)
	if err != nil {
		log.Printf("Error in handler callback: %s", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if response == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, *response.ToPacket())
}

func (s *HttpServer) handleCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ident, passwd, ok := parseLayrzAuth(r.Header.Get("Authorization"))
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if s.config.OnAuthenticate != nil && !s.config.OnAuthenticate(ident, passwd, r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if s.config.OnPullCommands == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, err := s.config.OnPullCommands(ident, passwd, r)
	if err != nil {
		log.Printf("Error in commands callback: %s", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if response == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, *response.ToPacket())
}

// parseLayrzAuth parses "LayrzAuth <ident>;<passwd>" from the Authorization header.
func parseLayrzAuth(h string) (ident, passwd string, ok bool) {
	const prefix = "LayrzAuth "
	if !strings.HasPrefix(h, prefix) {
		return "", "", false
	}
	rest := strings.TrimPrefix(h, prefix)
	ident, passwd, ok = strings.Cut(rest, ";")
	return
}
