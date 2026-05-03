package clients_test

import (
	"strings"
	"testing"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/clients"
	"github.com/goldenm-software/layrz-protocol/go/v3/definitions"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/ai"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/trips"
)

var parserFixedTime = time.Unix(1700000000, 0)

func parserPtr[T any](v T) *T { return &v }

func TestDecodeServerOutput(t *testing.T) {
	macAddr := "AA:BB:CC:DD:EE:FF"
	model := "GENERIC"
	abPacket := server.AbPacket{Devices: parserPtr([]definitions.BleData{{MacAddress: &macAddr, Model: &model}})}
	cmdName := "cmd1"
	acPacket := server.AcPacket{Commands: []definitions.CommandDefinition{{CommandId: 1, CommandName: &cmdName, Args: nil}}}
	aoPacket := server.AoPacket{Timestamp: parserFixedTime}
	arPacket := server.ArPacket{Reason: "test-reason"}
	asPacket := server.AsPacket{}
	//nolint:staticcheck
	auPacket := server.AuPacket{}
	tsPacket := trips.TsPacket{Timestamp: parserFixedTime, TripId: "trip-1"}
	tePacket := trips.TePacket{Timestamp: parserFixedTime, TripId: "trip-1"}
	imPacket := ai.ImPacket{Timestamp: parserFixedTime, ChatId: "chat-1", Message: "hello"}

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "Ab packet", input: *abPacket.ToPacket()},
		{name: "Ac packet", input: *acPacket.ToPacket()},
		{name: "Ao packet", input: *aoPacket.ToPacket()},
		{name: "Ar packet", input: *arPacket.ToPacket()},
		{name: "As packet", input: *asPacket.ToPacket()},
		{name: "Au packet", input: *auPacket.ToPacket()},
		{name: "Ts packet", input: *tsPacket.ToPacket()},
		{name: "Te packet", input: *tePacket.ToPacket()},
		{name: "Im packet", input: *imPacket.ToPacket()},
		{name: "unknown prefix", input: "<Xx>garbage</Xx>", wantErr: true},
		{name: "empty string", input: "", wantErr: true},
		{name: "corrupted Ao CRC", input: corruptCRC(*aoPacket.ToPacket()), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := clients.DecodeServerOutput(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result == nil || *result == nil {
				t.Fatal("result is nil")
			}
		})
	}
}

func TestEncodeClientPacket(t *testing.T) {
	ads := []definitions.BleAdvertisement{}
	data := []byte("test")

	tests := []struct {
		name    string
		packet  any
		wantErr bool
	}{
		{name: "PaPacket", packet: &client.PaPacket{Ident: parserPtr("id"), Password: parserPtr("pw")}},
		{name: "PbPacket", packet: &client.PbPacket{Advertisements: &ads}},
		{name: "PcPacket", packet: &client.PcPacket{Timestamp: parserFixedTime, CommandId: 1, Message: parserPtr("ok")}},
		{name: "PdPacket", packet: &client.PdPacket{Timestamp: parserFixedTime, ExtraData: map[string]any{}}},
		{name: "PiPacket", packet: &client.PiPacket{Ident: "IMEI123", FirmwareId: "fw1", FirmwareBuild: 1, DeviceId: 1, HardwareId: 1, ModelId: 1, FirmwareBranch: "main", FotaEnabled: false}},
		{name: "PsPacket", packet: &client.PsPacket{Timestamp: parserFixedTime, Params: map[string]any{}}},
		{name: "PmPacket", packet: &client.PmPacket{Filename: parserPtr("file.jpg"), ContentType: parserPtr("image/jpeg"), Data: &data}},
		{name: "PrPacket", packet: &client.PrPacket{}},
		{name: "TsPacket", packet: &trips.TsPacket{Timestamp: parserFixedTime, TripId: "trip-1"}},
		{name: "TePacket", packet: &trips.TePacket{Timestamp: parserFixedTime, TripId: "trip-1"}},
		{name: "ImPacket", packet: &ai.ImPacket{Timestamp: parserFixedTime, ChatId: "chat-1", Message: "hello"}},
		{name: "unknown type", packet: "not-a-packet", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := clients.EncodeClientPacket(tt.packet)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result == nil {
				t.Fatal("result is nil")
			}
		})
	}
}

// corruptCRC replaces the last 4 characters before the closing tag with "XXXX"
func corruptCRC(packet string) string {
	end := strings.LastIndex(packet, "</")
	if end < 4 {
		return packet
	}
	return packet[:end-4] + "XXXX" + packet[end:]
}
