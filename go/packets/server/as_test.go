package server_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
)

func TestAs_FromPacket_Errors(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{"wrong envelope", "<Xx>garbage</Xx>"},
		{"wrong part count", "<As>noCrc</As>"},
		{"bad CRC hex", "<As>;ZZZZ</As>"},
		{"CRC mismatch", "<As>;0000</As>"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw := tc.raw
			if err := (&server.AsPacket{}).FromPacket(&raw); err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

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
