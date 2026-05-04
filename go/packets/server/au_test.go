package server_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
)

func TestAu_FromPacket_Errors(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{"wrong envelope", "<Xx>garbage</Xx>"},
		{"wrong part count", "<Au>noCrc</Au>"},
		{"bad CRC hex", "<Au>;ZZZZ</Au>"},
		{"CRC mismatch", "<Au>;0000</Au>"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw := tc.raw
			if err := (&server.AuPacket{}).FromPacket(&raw); err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

func TestAu_FromPacket_ToPacket(t *testing.T) {
	packet := server.AuPacket{}
	encoded := *packet.ToPacket()

	raw := encoded
	decoded := server.AuPacket{}
	if err := decoded.FromPacket(&raw); err != nil {
		t.Fatalf("FromPacket failed: %v", err)
	}
	if *decoded.ToPacket() != encoded {
		t.Errorf("round-trip mismatch")
	}
}
