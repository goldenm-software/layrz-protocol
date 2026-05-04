package server_test

import (
	"testing"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/definitions"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
)

func TestDecode_AllServerPackets(t *testing.T) {
	fixedTime := time.Unix(1700000000, 0)
	bleDevices := []definitions.BleData{
		{MacAddress: stringPtr("AA:BB:CC:DD:EE:FF"), Model: stringPtr("GENERIC")},
	}
	commands := []definitions.CommandDefinition{
		{CommandId: 1, CommandName: stringPtr("cmd"), Args: map[string]any{}},
	}

	packets := []struct {
		name   string
		packet server.ServerPackets
	}{
		{"Ab", &server.AbPacket{Devices: &bleDevices}},
		{"Ac", &server.AcPacket{Commands: commands}},
		{"Ao", &server.AoPacket{Timestamp: fixedTime}},
		{"Ar", &server.ArPacket{Reason: "test error"}},
		{"As", &server.AsPacket{}},
		{"Au", &server.AuPacket{}},
	}

	for _, tt := range packets {
		t.Run(tt.name, func(t *testing.T) {
			encoded := *tt.packet.ToPacket()
			result, err := server.Decode([]byte(encoded))
			if err != nil {
				t.Fatalf("Decode(%s) failed: %v", tt.name, err)
			}
			if result == nil {
				t.Fatalf("Decode(%s) returned nil", tt.name)
			}
			if *result.ToPacket() != encoded {
				t.Errorf("%s round-trip mismatch: got %s, want %s", tt.name, *result.ToPacket(), encoded)
			}
		})
	}
}

func TestDecode_UnknownServerTag(t *testing.T) {
	_, err := server.Decode([]byte("<Xx>garbage</Xx>"))
	if err == nil {
		t.Error("expected error for unknown tag")
	}
}

func TestDecode_ServerEmpty(t *testing.T) {
	_, err := server.Decode([]byte(""))
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestDecode_MalformedAo(t *testing.T) {
	_, err := server.Decode([]byte("<Ao>notvalid</Ao>"))
	if err == nil {
		t.Error("expected error for malformed Ao body")
	}
}

func TestDecode_MalformedAr(t *testing.T) {
	_, err := server.Decode([]byte("<Ar>notvalid</Ar>"))
	if err == nil {
		t.Error("expected error for malformed Ar body")
	}
}
