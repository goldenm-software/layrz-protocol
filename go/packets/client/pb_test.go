package client_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/definitions"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
)

func TestPb_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name string
		ads  []definitions.BleAdvertisement
	}{
		{
			name: "single advertisement",
			ads: []definitions.BleAdvertisement{
				{
					MacAddress:       "AA:BB:CC:DD:EE:FF",
					Timestamp:        fixedTime,
					Latitude:         19.43,
					Longitude:        -99.18,
					Altitude:         2240.0,
					Rssi:             -50,
					TxPower:          -5,
					Model:            "GENERIC",
					DeviceName:       "Device1",
					ManufacturerData: []definitions.BleManufacturerData{},
					ServiceData:      []definitions.BleServiceData{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := client.PbPacket{Advertisements: &tt.ads}
			encoded := *packet.ToPacket()

			raw := encoded
			decoded := client.PbPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}
			if decoded.Advertisements == nil {
				t.Fatal("Advertisements is nil")
			}
			if len(*decoded.Advertisements) != len(tt.ads) {
				t.Errorf("advertisement count mismatch: got %d, want %d", len(*decoded.Advertisements), len(tt.ads))
			}
			if *decoded.ToPacket() != encoded {
				t.Errorf("round-trip mismatch")
			}
		})
	}
}
