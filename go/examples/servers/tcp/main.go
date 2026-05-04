package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
	"github.com/goldenm-software/layrz-protocol/go/v3/servers"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGTERM)
	defer cancel()

	s, err := servers.New(&servers.TcpConfig{
		Port: 12345,
		OnNewPacket: func(packet client.ClientPackets, conn net.Conn) (server.ServerPackets, error) {
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
		OnDecodeError: func(err error, data []byte, conn net.Conn) {
			fmt.Printf("Error decoding packet: %s\n", err.Error())
		},
	})

	if err != nil {
		panic(err)
	}

	defer s.Close()
	fmt.Printf("Listening on %d\n", 12345)
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
		return
	}
}
