package client_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
)

func TestPs_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]any
	}{
		{
			name:   "with parameters",
			params: map[string]any{"key": "value"},
		},
		{
			name:   "empty parameters",
			params: map[string]any{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := client.PsPacket{
				Timestamp: fixedTime,
				Params:    tt.params,
			}
			encoded := *packet.ToPacket()

			raw := encoded
			decoded := client.PsPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}
			if *decoded.ToPacket() != encoded {
				t.Errorf("round-trip mismatch")
			}
		})
	}
}
