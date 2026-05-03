package client_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
)

func TestPr_FromPacket_ToPacket(t *testing.T) {
	packet := client.PrPacket{}
	encoded := *packet.ToPacket()

	raw := encoded
	decoded := client.PrPacket{}
	if err := decoded.FromPacket(&raw); err != nil {
		t.Fatalf("FromPacket failed: %v", err)
	}
	if *decoded.ToPacket() != encoded {
		t.Errorf("round-trip mismatch")
	}
}
