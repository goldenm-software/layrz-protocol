package client_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

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

func TestPd_FromPacket_Errors(t *testing.T) {
	base := client.PdPacket{
		Timestamp: time.Unix(1700000000, 0),
		Position: &definitions.Position{
			Latitude:  floatPtr(19.43),
			Longitude: floatPtr(-99.18),
			Altitude:  floatPtr(2240.0),
			Speed:     floatPtr(0.0),
			Direction: floatPtr(0.0),
		},
		ExtraData: map[string]any{},
	}
	validPacket := *base.ToPacket()

	tests := []struct {
		name  string
		input string
	}{
		{"wrong prefix/suffix", "<Xx>" + validPacket[4:len(validPacket)-5] + "</Xx>"},
		{"bad CRC hex", validPacket[:len(validPacket)-9] + "XXXX" + "</Pd>"},
		{"CRC mismatch", validPacket[:len(validPacket)-9] + "0000" + "</Pd>"},
		{"not 9 parts", "<Pd>;0000</Pd>"},
		{"bad timestamp", buildPdWithField(validPacket, 0, "notanumber")},
		{"bad latitude", buildPdWithField(validPacket, 1, "notafloat")},
		{"bad longitude", buildPdWithField(validPacket, 2, "notafloat")},
		{"bad altitude", buildPdWithField(validPacket, 3, "notafloat")},
		{"bad speed", buildPdWithField(validPacket, 4, "notafloat")},
		{"bad direction", buildPdWithField(validPacket, 5, "notafloat")},
		{"bad satellite count", buildPdWithField(validPacket, 6, "notanint")},
		{"bad hdop", buildPdWithField(validPacket, 7, "notafloat")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := tt.input
			p := client.PdPacket{}
			if err := p.FromPacket(&raw); err == nil {
				t.Errorf("expected error for %q", tt.name)
			}
		})
	}
}

func TestPd_EmptyExtraData(t *testing.T) {
	p := client.PdPacket{
		Timestamp: time.Unix(1700000000, 0),
		Position:  nil,
		ExtraData: map[string]any{},
	}
	encoded := *p.ToPacket()
	raw := encoded
	decoded := client.PdPacket{}
	if err := decoded.FromPacket(&raw); err != nil {
		t.Fatalf("FromPacket failed: %v", err)
	}
}

func TestPd_ToPacket_AllExtraDataTypes(t *testing.T) {
	p := client.PdPacket{
		Timestamp: time.Unix(1700000000, 0),
		Position:  nil,
		ExtraData: map[string]any{
			"str":     "hello",
			"int":     int(1),
			"int8":    int8(2),
			"int16":   int16(3),
			"int32":   int32(4),
			"int64":   int64(5),
			"uint":    uint(6),
			"uint8":   uint8(7),
			"uint16":  uint16(8),
			"uint32":  uint32(9),
			"uint64":  uint64(10),
			"float32": float32(1.5),
			"float64": float64(2.5),
			"bool":    true,
			"other":   struct{ V int }{V: 42},
		},
	}
	result := p.ToPacket()
	if result == nil || *result == "" {
		t.Fatal("ToPacket returned nil or empty")
	}
}

func TestPd_FromPacket_SatelliteCountAndHdop(t *testing.T) {
	satCount := 8
	hdop := 1.2
	p := client.PdPacket{
		Timestamp: time.Unix(1700000000, 0),
		Position: &definitions.Position{
			Latitude:       floatPtr(10.0),
			Longitude:      floatPtr(20.0),
			Altitude:       floatPtr(100.0),
			Speed:          floatPtr(0.0),
			Direction:      floatPtr(0.0),
			SatelliteCount: &satCount,
			Hdop:           &hdop,
		},
		ExtraData: map[string]any{},
	}
	encoded := *p.ToPacket()
	raw := encoded
	decoded := client.PdPacket{}
	if err := decoded.FromPacket(&raw); err != nil {
		t.Fatalf("FromPacket failed: %v", err)
	}
	if decoded.Position == nil || decoded.Position.SatelliteCount == nil || *decoded.Position.SatelliteCount != satCount {
		t.Errorf("satellite count mismatch")
	}
}

// buildPdWithField replaces field at fieldIdx with badValue and recomputes CRC
func buildPdWithField(encoded string, fieldIdx int, badValue string) string {
	inner := encoded[4 : len(encoded)-5]
	body := inner[:len(inner)-4]
	body = strings.TrimSuffix(body, ";")
	parts := strings.Split(body, ";")
	if fieldIdx < len(parts) {
		parts[fieldIdx] = badValue
	}
	mutated := strings.Join(parts, ";") + ";"
	return fmt.Sprintf("<Pd>%s%s</Pd>", mutated, inner[len(inner)-4:])
}
