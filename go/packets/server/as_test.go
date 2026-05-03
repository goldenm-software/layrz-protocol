package server_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
)

func TestAs_FromPacket_ToPacket(t *testing.T) {
	packet := server.AsPacket{}
	encoded := *packet.ToPacket()

	raw := encoded
	decoded := server.AsPacket{}
	if err := decoded.FromPacket(&raw); err != nil {
		t.Fatalf("FromPacket failed: %v", err)
	}
	if *decoded.ToPacket() != encoded {
		t.Errorf("round-trip mismatch")
	}
}
