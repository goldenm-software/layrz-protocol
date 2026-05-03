package layrzprotocol_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/clients"
	"github.com/goldenm-software/layrz-protocol/go/v3/definitions"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/ai"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/trips"
)

var (
	testTime = time.Unix(1700000000, 0)
	testUUID = "12345678-1234-1234-1234-123456789012"
)

func TestDecodeServerOutput(t *testing.T) {
	tests := []struct {
		name        string
		builder     func() *string
		shouldError bool
	}{
		{
			name: "Ab packet",
			builder: func() *string {
				p := server.AbPacket{Devices: &[]definitions.BleData{{MacAddress: strPtr("AA:BB:CC:DD:EE:FF"), Model: strPtr("GENERIC")}}}
				return p.ToPacket()
			},
		},
		{
			name: "Ac packet",
			builder: func() *string {
				p := server.AcPacket{Commands: []definitions.CommandDefinition{{CommandId: 1, CommandName: strPtr("test"), Args: map[string]any{}}}}
				return p.ToPacket()
			},
		},
		{
			name: "Ao packet",
			builder: func() *string {
				p := server.AoPacket{Timestamp: testTime}
				return p.ToPacket()
			},
		},
		{
			name: "Ar packet",
			builder: func() *string {
				p := server.ArPacket{Reason: "error reason"}
				return p.ToPacket()
			},
		},
		{
			name: "As packet",
			builder: func() *string {
				p := server.AsPacket{}
				return p.ToPacket()
			},
		},
		{
			name: "Au packet",
			builder: func() *string {
				p := server.AuPacket{} //nolint:staticcheck
				return p.ToPacket()
			},
		},
		{
			name: "Ts packet",
			builder: func() *string {
				p := trips.TsPacket{Timestamp: testTime, TripId: testUUID}
				return p.ToPacket()
			},
		},
		{
			name: "Te packet",
			builder: func() *string {
				p := trips.TePacket{Timestamp: testTime, TripId: testUUID, DistanceTraveled: 100.0, MaxSpeed: 50.0, Duration: time.Hour}
				return p.ToPacket()
			},
		},
		{
			name: "Im packet",
			builder: func() *string {
				p := ai.ImPacket{Timestamp: testTime, ChatId: testUUID, Message: "hi"}
				return p.ToPacket()
			},
		},
		{
			name:        "invalid packet",
			builder:     func() *string { return strPtr("<Xx>invalid</Xx>") },
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := *tt.builder()
			result, err := clients.DecodeServerOutput(input)
			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("DecodeServerOutput failed: %v", err)
			}
			if result == nil {
				t.Fatal("result is nil")
			}
		})
	}
}

func TestEncodeClientPacket(t *testing.T) {
	tests := []struct {
		name        string
		packet      any
		shouldError bool
	}{
		{
			name:   "PaPacket",
			packet: &client.PaPacket{Ident: strPtr("test"), Password: strPtr("pass")},
		},
		{
			name:   "PbPacket",
			packet: &client.PbPacket{Advertisements: &[]definitions.BleAdvertisement{}},
		},
		{
			name:   "PcPacket",
			packet: &client.PcPacket{Timestamp: testTime, CommandId: 1, Message: strPtr("ok")},
		},
		{
			name:   "PdPacket",
			packet: &client.PdPacket{Timestamp: testTime, Position: nil, ExtraData: map[string]any{}},
		},
		{
			name:   "PiPacket",
			packet: &client.PiPacket{Ident: "test", FirmwareId: "fw1", FirmwareBuild: 1, DeviceId: 1, HardwareId: 1, ModelId: 1, FirmwareBranch: definitions.Stable, FotaEnabled: false},
		},
		{
			name:   "PsPacket",
			packet: &client.PsPacket{Timestamp: testTime, Params: map[string]any{}},
		},
		{
			name:   "PmPacket",
			packet: &client.PmPacket{Filename: strPtr("f.txt"), ContentType: strPtr("text/plain"), Data: &[]byte{1}},
		},
		{
			name:   "PrPacket",
			packet: &client.PrPacket{},
		},
		{
			name:   "TsPacket",
			packet: &trips.TsPacket{Timestamp: testTime, TripId: testUUID},
		},
		{
			name:   "TePacket",
			packet: &trips.TePacket{Timestamp: testTime, TripId: testUUID, DistanceTraveled: 100.0, MaxSpeed: 50.0, Duration: time.Hour},
		},
		{
			name:   "ImPacket",
			packet: &ai.ImPacket{Timestamp: testTime, ChatId: testUUID, Message: "hi"},
		},
		{
			name:        "invalid type",
			packet:      "not a packet",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := clients.EncodeClientPacket(tt.packet)
			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("EncodeClientPacket failed: %v", err)
			}
			if result == nil || len(*result) == 0 {
				t.Fatal("result is nil/empty")
			}
		})
	}
}

func strPtr(s string) *string { return &s }

// suppress unused warning
var _ = fmt.Sprintf
