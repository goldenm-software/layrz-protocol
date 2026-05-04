package client_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/definitions"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
)

func TestDecode_AllClientPackets(t *testing.T) {
	ads := []definitions.BleAdvertisement{{
		MacAddress:       "AA:BB:CC:DD:EE:FF",
		Timestamp:        fixedTime,
		Latitude:         floatPtr(19.43),
		Longitude:        floatPtr(-99.18),
		Altitude:         floatPtr(2240.0),
		Rssi:             -50,
		TxPower:          0,
		Model:            "GENERIC",
		DeviceName:       "Dev1",
		ManufacturerData: []definitions.BleManufacturerData{},
		ServiceData:      []definitions.BleServiceData{},
	}}
	mediaData := []byte("hello")

	packets := []struct {
		name   string
		packet client.ClientPackets
	}{
		{"Pa", &client.PaPacket{Ident: stringPtr("IMEI123"), Password: stringPtr("secret")}},
		{"Pb", &client.PbPacket{Advertisements: &ads}},
		{"Pc", &client.PcPacket{Timestamp: fixedTime, CommandId: 1, Message: stringPtr("OK")}},
		{"Pd", &client.PdPacket{Timestamp: fixedTime, ExtraData: map[string]any{}}},
		{"Pi", &client.PiPacket{
			Ident:          "device1",
			FirmwareId:     "fw1",
			FirmwareBuild:  1,
			DeviceId:       1,
			HardwareId:     1,
			ModelId:        1,
			FirmwareBranch: definitions.Stable,
			FotaEnabled:    false,
		}},
		{"Pm", &client.PmPacket{
			Filename:    stringPtr("test.txt"),
			ContentType: stringPtr("text/plain"),
			Data:        &mediaData,
		}},
		{"Pr", &client.PrPacket{}},
		{"Ps", &client.PsPacket{Timestamp: fixedTime, Params: map[string]any{}}},
	}

	for _, tt := range packets {
		t.Run(tt.name, func(t *testing.T) {
			encoded := *tt.packet.ToPacket()
			result, err := client.Decode([]byte(encoded))
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

func TestDecode_UnknownClientTag(t *testing.T) {
	_, err := client.Decode([]byte("<Xx>garbage</Xx>"))
	if err == nil {
		t.Error("expected error for unknown tag")
	}
}

func TestDecode_ClientEmpty(t *testing.T) {
	_, err := client.Decode([]byte(""))
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestDecode_MalformedPa(t *testing.T) {
	_, err := client.Decode([]byte("<Pa>notvalid</Pa>"))
	if err == nil {
		t.Error("expected error for malformed Pa body")
	}
}

func TestDecode_MalformedPr(t *testing.T) {
	_, err := client.Decode([]byte("<Pr>notvalid</Pr>"))
	if err == nil {
		t.Error("expected error for malformed Pr body")
	}
}
