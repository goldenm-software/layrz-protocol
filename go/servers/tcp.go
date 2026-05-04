package servers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/helpers"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
	"github.com/pires/go-proxyproto"
)

// TcpServer is a TCP server that listens for incoming connections and processes the incoming packets
type TcpServer struct {
	config             *TcpConfig
	ctx                context.Context
	cancel             context.CancelFunc
	accumulatedPerPort map[int][]byte
}

// TcpConfig is the configuration for the TCP server
type TcpConfig struct {
	// Defines the TCP Port for the TCP server to listen
	Port int
	// Enables Proxy Protocol v2 support, by default is disabled
	ProxyProtocolV2 bool
	// Handler on new packet received, the response is optional, if nil, no response will be sent
	// however, if you need to send a response, you must return a server.ServerPackets
	OnNewPacket func(packet client.ClientPackets, conn net.Conn) (server.ServerPackets, error)

	// Is the defined callback when something went wrong on decoder
	OnDecodeError func(err error, data []byte, conn net.Conn)
}

// Creates a new TCP server with the given configuration
func New(cfg *TcpConfig) (*TcpServer, error) {
	if cfg == nil {
		cfg = &TcpConfig{}
	}

	if cfg.OnNewPacket == nil {
		return nil, fmt.Errorf("on new packet handler is not set")
	}

	if cfg.OnDecodeError == nil {
		cfg.OnDecodeError = func(err error, data []byte, conn net.Conn) {
			log.Printf("Error decoding packet: %s Data: %s", err.Error(), string(data))
		}
	}

	if cfg.Port <= 0 || cfg.Port >= 65535 {
		return nil, fmt.Errorf("port is not valid")
	}

	return &TcpServer{config: cfg}, nil
}

// Starts the TCP server and listens for incoming connections
// This method is blocking and will run until the server is closed
// or an error occurs
func (s *TcpServer) Start(ctx context.Context) error {
	subctx, cancel := context.WithCancel(ctx)
	s.ctx = subctx
	s.cancel = cancel

	s.accumulatedPerPort = make(map[int][]byte)

	defer cancel()

	var ln net.Listener
	var err error

	if s.config.ProxyProtocolV2 {
		rawLn, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
		if err != nil {
			return err
		}

		ln = &proxyproto.Listener{Listener: rawLn}
	} else {
		ln, err = net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	}

	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}

// Handles a new connection and processes the incoming packets
// This method is not thread-safe and should be called in a separate goroutine
func (s *TcpServer) handleConnection(conn net.Conn) {
	port := s.getPort(conn)
	defer func() {
		_ = conn.Close()
		delete(s.accumulatedPerPort, port)
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("Error reading from connection: %s", err.Error())
			return
		}

		s.accumulatedPerPort[port] = append(s.accumulatedPerPort[port], buf[:n]...)
		if !bytes.ContainsRune(s.accumulatedPerPort[port], '\n') {
			continue
		}

		messages := helpers.Split(string(s.accumulatedPerPort[port]))
		s.accumulatedPerPort[port] = []byte{}

		for _, message := range messages {
			packet, err := client.Decode([]byte(message))
			if err != nil {
				s.config.OnDecodeError(err, []byte(message), conn)
				continue
			}

			response, err := s.config.OnNewPacket(packet, conn)
			if err != nil {
				log.Printf("Error in handler callback: %s", err.Error())
				continue
			}

			if response != nil {
				responseStr := response.ToPacket()
				_, err = conn.Write([]byte(*responseStr))
				if err != nil {
					log.Printf("Error writing to connection: %s", err.Error())
					continue
				}
			}
		}
	}
}

// Close the TCP server and release the port
func (s *TcpServer) Close() error {
	s.cancel()
	return nil
}

// Helper function to get the port from a connection
func (s *TcpServer) getPort(conn net.Conn) int {
	if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		return addr.Port
	}
	return 0
}
