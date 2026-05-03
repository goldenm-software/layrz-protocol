package client_test

import (
	"fmt"
	"strings"
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
		{
			name: "advertisement with manufacturer and service data",
			ads: []definitions.BleAdvertisement{
				{
					MacAddress: "11:22:33:44:55:66",
					Timestamp:  fixedTime,
					Latitude:   10.0,
					Longitude:  20.0,
					Altitude:   100.0,
					Rssi:       -60,
					TxPower:    0,
					Model:      "MODEL1",
					DeviceName: "Dev2",
					ManufacturerData: []definitions.BleManufacturerData{
						{CompanyId: 0x004C, Data: []byte{0x01, 0x02}},
					},
					ServiceData: []definitions.BleServiceData{
						{Uuid: 0xFEAA, Data: []byte{0xAB, 0xCD}},
					},
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

func TestPb_FromPacket_Errors(t *testing.T) {
	// Build a valid encoded packet to use as a base for mutations
	ads := []definitions.BleAdvertisement{{
		MacAddress:       "AA:BB:CC:DD:EE:FF",
		Timestamp:        fixedTime,
		Latitude:         19.43,
		Longitude:        -99.18,
		Altitude:         2240.0,
		Rssi:             -50,
		TxPower:          0,
		Model:            "GENERIC",
		DeviceName:       "Device1",
		ManufacturerData: []definitions.BleManufacturerData{},
		ServiceData:      []definitions.BleServiceData{},
	}}
	pbPkt := client.PbPacket{Advertisements: &ads}
	validPacket := *pbPkt.ToPacket()

	tests := []struct {
		name  string
		input string
	}{
		{"wrong prefix/suffix", "<Xx>" + validPacket[4:len(validPacket)-5] + "</Xx>"},
		{"bad outer CRC hex", validPacket[:len(validPacket)-9] + "XXXX" + "</Pb>"},
		{"outer CRC mismatch", validPacket[:len(validPacket)-9] + "0000" + "</Pb>"},
		{"not multiple of 12 parts", "<Pb>;0000</Pb>"},
		{"bad inner CRC hex", buildPbWithBadInnerCRC(validPacket)},
		{"bad timestamp", buildPbWithBadField(validPacket, 1, "notanumber")},
		{"bad latitude", buildPbWithBadField(validPacket, 2, "notafloat")},
		{"bad longitude", buildPbWithBadField(validPacket, 3, "notafloat")},
		{"bad altitude", buildPbWithBadField(validPacket, 4, "notafloat")},
		{"bad rssi", buildPbWithBadField(validPacket, 7, "notanint")},
		{"bad tx power", buildPbWithBadField(validPacket, 8, "notanint")},
		{"bad manufacturer company id", buildPbWithManufacturerData(validPacket, "XXXX:0102")},
		{"bad manufacturer data bytes", buildPbWithManufacturerData(validPacket, "004C:ZZ")},
		{"bad service uuid", buildPbWithServiceData(validPacket, "XXXX:0102")},
		{"bad service data bytes", buildPbWithServiceData(validPacket, "FEAA:ZZ")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := tt.input
			p := client.PbPacket{}
			if err := p.FromPacket(&raw); err == nil {
				t.Errorf("expected error for %q", tt.name)
			}
		})
	}
}

// buildPbWithBadInnerCRC replaces the inner CRC of the first advertisement with "XXXX"
func buildPbWithBadInnerCRC(encoded string) string {
	inner := encoded[4 : len(encoded)-5] // strip <Pb>...</Pb>
	// The inner CRC is the last 4 chars before the outer CRC
	outerCRC := inner[len(inner)-4:]
	body := inner[:len(inner)-4]
	// body ends with ";<inner_crc>" for the first ad — replace inner crc
	lastSemi := strings.LastIndex(body, ";")
	if lastSemi < 0 {
		return encoded
	}
	mutated := body[:lastSemi+1] + "XXXX"
	return fmt.Sprintf("<Pb>%s%s</Pb>", mutated, outerCRC)
}

// buildPbWithBadField replaces field at index fieldIdx (0-based within advertisement parts)
func buildPbWithBadField(encoded string, fieldIdx int, badValue string) string {
	inner := encoded[4 : len(encoded)-5]
	outerCRC := inner[len(inner)-4:]
	body := inner[:len(inner)-4]
	body = strings.TrimSuffix(body, ";")
	parts := strings.Split(body, ";")
	if fieldIdx < len(parts) {
		parts[fieldIdx] = badValue
	}
	mutated := strings.Join(parts, ";") + ";"
	return fmt.Sprintf("<Pb>%s%s</Pb>", mutated, outerCRC)
}

// buildPbWithManufacturerData inserts a given raw manufacturer data string into the packet
func buildPbWithManufacturerData(encoded, mfgData string) string {
	inner := encoded[4 : len(encoded)-5]
	outerCRC := inner[len(inner)-4:]
	body := inner[:len(inner)-4]
	body = strings.TrimSuffix(body, ";")
	parts := strings.Split(body, ";")
	// field 9 is manufacturer data
	if len(parts) > 9 {
		parts[9] = mfgData
	}
	mutated := strings.Join(parts, ";") + ";"
	return fmt.Sprintf("<Pb>%s%s</Pb>", mutated, outerCRC)
}

// buildPbWithServiceData inserts a given raw service data string into the packet
func buildPbWithServiceData(encoded, svcData string) string {
	inner := encoded[4 : len(encoded)-5]
	outerCRC := inner[len(inner)-4:]
	body := inner[:len(inner)-4]
	body = strings.TrimSuffix(body, ";")
	parts := strings.Split(body, ";")
	// field 10 is service data
	if len(parts) > 10 {
		parts[10] = svcData
	}
	mutated := strings.Join(parts, ";") + ";"
	return fmt.Sprintf("<Pb>%s%s</Pb>", mutated, outerCRC)
}
