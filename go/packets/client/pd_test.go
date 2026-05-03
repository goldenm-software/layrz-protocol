package client_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/definitions"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
)

func TestPd_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		position  *definitions.Position
	}{
		{
			name:      "with position",
			timestamp: 1700000000,
			position: &definitions.Position{
				Latitude:  floatPtr(19.43),
				Longitude: floatPtr(-99.18),
				Altitude:  floatPtr(2240.0),
			},
		},
		{
			name:      "without position",
			timestamp: 1700000001,
			position:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := client.PdPacket{
				Timestamp: fixedTime,
				Position:  tt.position,
				ExtraData: map[string]any{},
			}
			encoded := *packet.ToPacket()

			raw := encoded
			decoded := client.PdPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}
			if *decoded.ToPacket() != encoded {
				t.Errorf("round-trip mismatch")
			}
		})
	}
}
