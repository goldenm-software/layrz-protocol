package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
	lservers "github.com/goldenm-software/layrz-protocol/go/v3/servers"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	s, err := lservers.NewHttp(&lservers.HttpConfig{
		Port: 12345,
		OnAuthenticate: func(ident, passwd string, r *http.Request) bool {
			// Replace with real credential validation
			return ident == "device001" && passwd == "secret"
		},
		OnNewPacket: func(packet client.ClientPackets, r *http.Request) (server.ServerPackets, error) {
			switch packet.(type) {
			case *client.PaPacket:
				fmt.Printf("Pa received\n")
				return &server.AsPacket{}, nil
			case *client.PbPacket:
				fmt.Printf("Pb received\n")
			case *client.PcPacket:
				fmt.Printf("Pc received\n")
			case *client.PdPacket:
				fmt.Printf("Pd received\n")
			case *client.PiPacket:
				fmt.Printf("Pi received\n")
			case *client.PmPacket:
				fmt.Printf("Pm received\n")
			case *client.PrPacket:
				fmt.Printf("Pr received\n")
			case *client.PsPacket:
				fmt.Printf("Ps received\n")
			}

			return &server.AoPacket{Timestamp: time.Now()}, nil
		},
		OnPullCommands: func(ident, passwd string, r *http.Request) (server.ServerPackets, error) {
			fmt.Printf("Commands requested by %s\n", ident)
			// Return nil when there are no pending commands for the device
			return nil, nil
		},
		OnDecodeError: func(err error, data []byte, r *http.Request) {
			fmt.Printf("Error decoding packet: %s\n", err.Error())
		},
	})

	if err != nil {
		panic(err)
	}

	defer s.Close()
	fmt.Printf("Listening on :%d\n", 12345)

	errChan := make(chan error, 1)
	go func() {
		if err := s.Start(ctx); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		fmt.Printf("Error starting server: %s\n", err.Error())
	case <-ctx.Done():
		fmt.Println("Server closed")
	}
}
