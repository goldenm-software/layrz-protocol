package server_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/definitions"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
)

func TestAb_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name    string
		devices []definitions.BleData
	}{
		{
			name: "single device",
			devices: []definitions.BleData{
				{MacAddress: stringPtr("12:34:56:78:90:AB"), Model: stringPtr("GENERIC")},
			},
		},
		{
			name: "multiple devices",
			devices: []definitions.BleData{
				{MacAddress: stringPtr("1234567890AB"), Model: stringPtr("GENERIC")},
				{MacAddress: stringPtr("BC0987654321"), Model: stringPtr("GENERIC")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := server.AbPacket{Devices: &tt.devices}
			encoded := *packet.ToPacket()

			raw := encoded
			decoded := server.AbPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}
			if decoded.Devices == nil {
				t.Fatal("Devices is nil")
			}
			if len(*decoded.Devices) != len(tt.devices) {
				t.Errorf("device count mismatch: got %d, want %d", len(*decoded.Devices), len(tt.devices))
			}
			reencoded := *decoded.ToPacket()
			if reencoded != encoded {
				t.Errorf("round-trip mismatch:\n  got  %s\n  want %s", reencoded, encoded)
			}
		})
	}
}
