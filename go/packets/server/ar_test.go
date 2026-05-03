package server_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
)

func TestAr_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name   string
		reason string
	}{
		{
			name:   "error reason",
			reason: "Invalid command",
		},
		{
			name:   "empty reason",
			reason: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := server.ArPacket{Reason: tt.reason}
			encoded := *packet.ToPacket()

			raw := encoded
			decoded := server.ArPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}
			if decoded.Reason != tt.reason {
				t.Errorf("reason mismatch: got %s, want %s", decoded.Reason, tt.reason)
			}
			if *decoded.ToPacket() != encoded {
				t.Errorf("round-trip mismatch")
			}
		})
	}
}
