package main

import (
	"fmt"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
)

func handleOutput(response *any) {
	var data *string
	var packet = *response

	switch p := packet.(type) {
	case *server.AbPacket:
		data = p.ToPacket()
	case *server.AcPacket:
		data = p.ToPacket()
	case *server.AoPacket:
		data = p.ToPacket()
	case *server.ArPacket:
		data = p.ToPacket()
	case *server.AsPacket:
		data = p.ToPacket()
	case *server.AuPacket:
		data = p.ToPacket()
	default:
		fmt.Printf("Unknown packet type received: %T\n", response)
		return
	}

	fmt.Printf("Packet received: %s\n", *data)
}
