package main

import (
	"fmt"

	"github.com/goldenm-software/layrz-protocol/golang/v2"
)

func handleOutput(response *interface{}) {
	var data *string
	var packet = *response

	switch packet.(type) {
	case *layrzprotocol.AbPacket:
		data = packet.(*layrzprotocol.AbPacket).ToPacket()
	case *layrzprotocol.AcPacket:
		data = packet.(*layrzprotocol.AcPacket).ToPacket()
	case *layrzprotocol.AoPacket:
		data = packet.(*layrzprotocol.AoPacket).ToPacket()
	case *layrzprotocol.ArPacket:
		data = packet.(*layrzprotocol.ArPacket).ToPacket()
	case *layrzprotocol.AsPacket:
		data = packet.(*layrzprotocol.AsPacket).ToPacket()
	case *layrzprotocol.AuPacket:
		data = packet.(*layrzprotocol.AuPacket).ToPacket()
	default:
		fmt.Printf("Unknown packet type received: %T\n", response)
		return
	}

	fmt.Printf("Packet received: %s\n", *data)
}
