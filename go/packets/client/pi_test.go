package client_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/definitions"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
)

func TestPi_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name   string
		ident  string
		branch definitions.FirmwareBranch
	}{
		{
			name:   "stable firmware",
			ident:  "test_device",
			branch: definitions.Stable,
		},
		{
			name:   "development firmware",
			ident:  "dev_device",
			branch: definitions.Development,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := client.PiPacket{
				Ident:          tt.ident,
				FirmwareId:     "fw1",
				FirmwareBuild:  1,
				DeviceId:       1,
				HardwareId:     1,
				ModelId:        1,
				FirmwareBranch: tt.branch,
				FotaEnabled:    false,
			}
			encoded := *packet.ToPacket()

			raw := encoded
			decoded := client.PiPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}
			if decoded.FirmwareBranch != tt.branch {
				t.Errorf("branch mismatch: got %s, want %s", decoded.FirmwareBranch, tt.branch)
			}
			if *decoded.ToPacket() != encoded {
				t.Errorf("round-trip mismatch")
			}
		})
	}
}
