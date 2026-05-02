package main

import (
	"fmt"

	layrzprotocol "github.com/goldenm-software/layrz-protocol/go/v3"
)

func handleOutput(response *any) {
	var data *string
	var packet = *response

	switch p := packet.(type) {
	case *layrzprotocol.AbPacket:
		data = p.ToPacket()
	case *layrzprotocol.AcPacket:
		data = p.ToPacket()
	case *layrzprotocol.AoPacket:
		data = p.ToPacket()
	case *layrzprotocol.ArPacket:
		data = p.ToPacket()
	case *layrzprotocol.AsPacket:
		data = p.ToPacket()
	case *layrzprotocol.AuPacket:
		data = p.ToPacket()
	default:
		fmt.Printf("Unknown packet type received: %T\n", response)
		return
	}

	fmt.Printf("Packet received: %s\n", *data)
}
