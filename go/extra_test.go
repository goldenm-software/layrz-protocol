package layrzprotocol

import (
	"fmt"
	"testing"
	"time"
)

// TestAb_FromPacket_ToPacket tests Ab packet encoding/decoding round-trip
func TestAb_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name    string
		devices []BleData
	}{
		{
			name: "single device",
			devices: []BleData{
				{
					MacAddress: stringPtr("12:34:56:78:90:AB"),
					Model:      stringPtr("GENERIC"),
				},
			},
		},
		{
			name: "multiple devices",
			devices: []BleData{
				{
					MacAddress: stringPtr("1234567890AB"),
					Model:      stringPtr("GENERIC"),
				},
				{
					MacAddress: stringPtr("BC0987654321"),
					Model:      stringPtr("GENERIC"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build canonical frame
			packet := AbPacket{Devices: &tt.devices}
			encoded := *packet.ToPacket()

			// Decode it back
			raw := copy(encoded)
			decoded := AbPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}

			// Verify devices match
			if decoded.Devices == nil {
				t.Fatal("Devices is nil")
			}
			if len(*decoded.Devices) != len(tt.devices) {
				t.Errorf("device count mismatch: got %d, want %d", len(*decoded.Devices), len(tt.devices))
			}

			// Re-encode and verify round-trip
			reencoded := *decoded.ToPacket()
			if reencoded != encoded {
				t.Errorf("round-trip mismatch:\n  got  %s\n  want %s", reencoded, encoded)
			}
		})
	}
}

// TestAc_FromPacket_ToPacket tests Ac packet encoding/decoding round-trip
func TestAc_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name     string
		commands []CommandDefinition
	}{
		{
			name: "single command no args",
			commands: []CommandDefinition{
				{
					CommandId:   1,
					CommandName: stringPtr("test"),
					Args:        &map[string]any{},
				},
			},
		},
		{
			name: "command with int arg",
			commands: []CommandDefinition{
				{
					CommandId:   42,
					CommandName: stringPtr("setspeed"),
					Args: &map[string]any{
						"value": 100,
					},
				},
			},
		},
		{
			name: "multiple commands",
			commands: []CommandDefinition{
				{
					CommandId:   1,
					CommandName: stringPtr("cmd1"),
					Args: &map[string]any{
						"param1": "hello",
					},
				},
				{
					CommandId:   2,
					CommandName: stringPtr("cmd2"),
					Args: &map[string]any{
						"param2": 42,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := AcPacket{Commands: tt.commands}
			encoded := *packet.ToPacket()

			raw := copy(encoded)
			decoded := AcPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}

			if len(decoded.Commands) != len(tt.commands) {
				t.Errorf("command count mismatch: got %d, want %d", len(decoded.Commands), len(tt.commands))
			}

			reencoded := *decoded.ToPacket()
			if reencoded != encoded {
				t.Errorf("round-trip mismatch:\n  got  %s\n  want %s", reencoded, encoded)
			}
		})
	}
}

// TestPb_FromPacket_ToPacket tests Pb packet encoding/decoding round-trip
func TestPb_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name           string
		advertisements []BleAdvertisement
	}{
		{
			name: "single advertisement minimal",
			advertisements: []BleAdvertisement{
				{
					MacAddress:       "AA:BB:CC:DD:EE:FF",
					Timestamp:        fixedTime,
					Latitude:         10.0,
					Longitude:        20.0,
					Altitude:         30.0,
					Rssi:             -50,
					TxPower:          0,
					Model:            "GENERIC",
					DeviceName:       "test",
					ManufacturerData: []BleManufacturerData{{CompanyId: 0x004C, Data: []byte{0x01}}},
					ServiceData:      []BleServiceData{},
				},
			},
		},
		{
			name: "advertisement with manufacturer data",
			advertisements: []BleAdvertisement{
				{
					MacAddress: "11:22:33:44:55:66",
					Timestamp:  fixedTime,
					Latitude:   5.0,
					Longitude:  15.0,
					Altitude:   25.0,
					Rssi:       -60,
					TxPower:    -10,
					Model:      "BLE_MODEL",
					DeviceName: "device1",
					ManufacturerData: []BleManufacturerData{
						{
							CompanyId: 0x004C,
							Data:      []byte{0x02, 0x01, 0x04},
						},
					},
					ServiceData: []BleServiceData{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := PbPacket{Advertisements: &tt.advertisements}
			encoded := *packet.ToPacket()

			raw := copy(encoded)
			decoded := PbPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}

			if decoded.Advertisements == nil {
				t.Fatal("Advertisements is nil")
			}

			reencoded := *decoded.ToPacket()
			if reencoded != encoded {
				t.Errorf("round-trip mismatch:\n  got  %s\n  want %s", reencoded, encoded)
			}
		})
	}
}

// TestPd_FromPacket_ToPacket tests Pd packet encoding/decoding round-trip
func TestPd_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name      string
		timestamp time.Time
		position  *Position
		extraData *map[string]any
	}{
		{
			name:      "no position, no extras",
			timestamp: fixedTime,
			position:  nil,
			extraData: &map[string]any{},
		},
		{
			name:      "with position",
			timestamp: fixedTime,
			position: &Position{
				Latitude:       floatPtr(10.0),
				Longitude:      floatPtr(20.0),
				Altitude:       floatPtr(30.0),
				Speed:          floatPtr(50.0),
				Direction:      floatPtr(180.0),
				SatelliteCount: intPtr(12),
				Hdop:           floatPtr(1.5),
			},
			extraData: &map[string]any{},
		},
		{
			name:      "with position and extras",
			timestamp: fixedTime,
			position: &Position{
				Latitude:  floatPtr(10.123456),
				Longitude: floatPtr(20.654321),
				Altitude:  floatPtr(100.5),
				Speed:     floatPtr(25.0),
				Direction: floatPtr(90.0),
			},
			extraData: &map[string]any{
				"temperature": 25,
				"humidity":    60.5,
				"status":      "ok",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := PdPacket{
				Timestamp: tt.timestamp,
				Position:  tt.position,
				ExtraData: tt.extraData,
			}
			encoded := *packet.ToPacket()

			raw := copy(encoded)
			decoded := PdPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}

			if decoded.Timestamp.Unix() != tt.timestamp.Unix() {
				t.Errorf("timestamp mismatch: got %d, want %d", decoded.Timestamp.Unix(), tt.timestamp.Unix())
			}

			// Map iteration order is non-deterministic; verify key count instead of string equality
			if tt.extraData != nil && decoded.ExtraData != nil && len(*tt.extraData) > 0 {
				if len(*decoded.ExtraData) != len(*tt.extraData) {
					t.Errorf("extra data count mismatch: got %d, want %d", len(*decoded.ExtraData), len(*tt.extraData))
				}
			} else {
				reencoded := *decoded.ToPacket()
				if reencoded != encoded {
					t.Errorf("round-trip mismatch:\n  got  %s\n  want %s", reencoded, encoded)
				}
			}
		})
	}
}

// TestPi_FromPacket_ToPacket tests Pi packet encoding/decoding round-trip
func TestPi_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name           string
		ident          string
		firmwareId     string
		firmwareBuild  int
		deviceId       int
		hardwareId     int
		modelId        int
		firmwareBranch FirmwareBranch
		fotaEnabled    bool
	}{
		{
			name:           "stable firmware no fota",
			ident:          "123456789012345",
			firmwareId:     "fw_001",
			firmwareBuild:  100,
			deviceId:       1,
			hardwareId:     2,
			modelId:        3,
			firmwareBranch: Stable,
			fotaEnabled:    false,
		},
		{
			name:           "development firmware with fota",
			ident:          "IMEI12345",
			firmwareId:     "layrz_v2",
			firmwareBuild:  200,
			deviceId:       10,
			hardwareId:     20,
			modelId:        30,
			firmwareBranch: Development,
			fotaEnabled:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := PiPacket{
				Ident:          tt.ident,
				FirmwareId:     tt.firmwareId,
				FirmwareBuild:  tt.firmwareBuild,
				DeviceId:       tt.deviceId,
				HardwareId:     tt.hardwareId,
				ModelId:        tt.modelId,
				FirmwareBranch: tt.firmwareBranch,
				FotaEnabled:    tt.fotaEnabled,
			}
			encoded := *packet.ToPacket()

			raw := copy(encoded)
			decoded := PiPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}

			if decoded.Ident != tt.ident {
				t.Errorf("ident mismatch: got %q, want %q", decoded.Ident, tt.ident)
			}
			if decoded.FirmwareBuild != tt.firmwareBuild {
				t.Errorf("firmware build mismatch: got %d, want %d", decoded.FirmwareBuild, tt.firmwareBuild)
			}
			if decoded.FotaEnabled != tt.fotaEnabled {
				t.Errorf("fota enabled mismatch: got %v, want %v", decoded.FotaEnabled, tt.fotaEnabled)
			}

			reencoded := *decoded.ToPacket()
			if reencoded != encoded {
				t.Errorf("round-trip mismatch:\n  got  %s\n  want %s", reencoded, encoded)
			}
		})
	}
}

// TestPm_FromPacket_ToPacket tests Pm packet encoding/decoding round-trip
func TestPm_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		contentType string
		data        []byte
	}{
		{
			name:        "text file",
			filename:    "test.txt",
			contentType: "text/plain",
			data:        []byte("hello world"),
		},
		{
			name:        "empty file",
			filename:    "empty.bin",
			contentType: "application/octet-stream",
			data:        []byte{},
		},
		{
			name:        "binary data",
			filename:    "data.bin",
			contentType: "application/octet-stream",
			data:        []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := PmPacket{
				Filename:    stringPtr(tt.filename),
				ContentType: stringPtr(tt.contentType),
				Data:        &tt.data,
			}
			encoded := *packet.ToPacket()

			raw := copy(encoded)
			decoded := PmPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}

			if *decoded.Filename != tt.filename {
				t.Errorf("filename mismatch: got %q, want %q", *decoded.Filename, tt.filename)
			}
			if *decoded.ContentType != tt.contentType {
				t.Errorf("content type mismatch: got %q, want %q", *decoded.ContentType, tt.contentType)
			}

			reencoded := *decoded.ToPacket()
			if reencoded != encoded {
				t.Errorf("round-trip mismatch:\n  got  %s\n  want %s", reencoded, encoded)
			}
		})
	}
}

// TestPs_FromPacket_ToPacket tests Ps packet encoding/decoding round-trip
func TestPs_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name      string
		timestamp time.Time
		params    *map[string]any
	}{
		{
			name:      "no params",
			timestamp: fixedTime,
			params:    &map[string]any{},
		},
		{
			name:      "with params",
			timestamp: fixedTime,
			params: &map[string]any{
				"interval":  60,
				"threshold": 25.5,
				"enabled":   true,
				"name":      "device1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := PsPacket{
				Timestamp: tt.timestamp,
				Params:    tt.params,
			}
			encoded := *packet.ToPacket()

			raw := copy(encoded)
			decoded := PsPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}

			if decoded.Timestamp.Unix() != tt.timestamp.Unix() {
				t.Errorf("timestamp mismatch: got %d, want %d", decoded.Timestamp.Unix(), tt.timestamp.Unix())
			}

			// Map iteration order is non-deterministic; verify key count instead of string equality
			if tt.params != nil && decoded.Params != nil {
				if len(*decoded.Params) != len(*tt.params) {
					t.Errorf("params count mismatch: got %d, want %d", len(*decoded.Params), len(*tt.params))
				}
			} else if decoded.Timestamp.Unix() != tt.timestamp.Unix() {
				t.Errorf("round-trip mismatch:\n  got  %s\n  want %s", *decoded.ToPacket(), encoded)
			}
		})
	}
}

// TestHandleServerOutput tests the handleServerOutput parser function
func TestHandleServerOutput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectType  string
		shouldError bool
	}{
		{
			name:       "Ab packet",
			input:      "<Ab>1234567890AB:GENERIC;362C</Ab>",
			expectType: "*layrzprotocol.AbPacket",
		},
		{
			name:       "Ac packet simple",
			input:      "<Ac>1;test;;5CAD;28E6</Ac>",
			expectType: "*layrzprotocol.AcPacket",
		},
		{
			name:       "Ao packet",
			input:      canonicalFrames["Ao"],
			expectType: "*layrzprotocol.AoPacket",
		},
		{
			name:       "Ar packet",
			input:      canonicalFrames["Ar"],
			expectType: "*layrzprotocol.ArPacket",
		},
		{
			name:       "As packet",
			input:      canonicalFrames["As"],
			expectType: "*layrzprotocol.AsPacket",
		},
		{
			name:       "Au packet",
			input:      canonicalFrames["Au"],
			expectType: "*layrzprotocol.AuPacket",
		},
		{
			name:       "Ts packet",
			input:      canonicalFrames["Ts"],
			expectType: "*layrzprotocol.TsPacket",
		},
		{
			name:       "Te packet",
			input:      canonicalFrames["Te"],
			expectType: "*layrzprotocol.TePacket",
		},
		{
			name:       "Im packet",
			input:      canonicalFrames["Im"],
			expectType: "*layrzprotocol.ImPacket",
		},
		{
			name:        "invalid packet",
			input:       "<Xx>invalid</Xx>",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handleServerOutput(tt.input)

			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("handleServerOutput failed: %v", err)
			}
			if result == nil {
				t.Fatal("result is nil")
			}
		})
	}
}

// TestParsePacketToString tests the parsePacketToString function
func TestParsePacketToString(t *testing.T) {
	tests := []struct {
		name        string
		packet      any
		shouldError bool
	}{
		{
			name:   "PaPacket",
			packet: &PaPacket{Ident: stringPtr("test"), Password: stringPtr("pass")},
		},
		{
			name:   "PbPacket",
			packet: &PbPacket{Advertisements: &[]BleAdvertisement{}},
		},
		{
			name:   "PcPacket",
			packet: &PcPacket{Timestamp: fixedTime, CommandId: 1, Message: stringPtr("ok")},
		},
		{
			name: "PdPacket",
			packet: &PdPacket{
				Timestamp: fixedTime,
				Position:  nil,
				ExtraData: &map[string]any{},
			},
		},
		{
			name: "PiPacket",
			packet: &PiPacket{
				Ident:          "test",
				FirmwareId:     "fw1",
				FirmwareBuild:  1,
				DeviceId:       1,
				HardwareId:     1,
				ModelId:        1,
				FirmwareBranch: Stable,
				FotaEnabled:    false,
			},
		},
		{
			name: "PsPacket",
			packet: &PsPacket{
				Timestamp: fixedTime,
				Params:    &map[string]any{},
			},
		},
		{
			name: "PmPacket",
			packet: &PmPacket{
				Filename:    stringPtr("test.txt"),
				ContentType: stringPtr("text/plain"),
				Data:        &[]byte{1, 2, 3},
			},
		},
		{
			name:   "PrPacket",
			packet: &PrPacket{},
		},
		{
			name:   "TsPacket",
			packet: &TsPacket{Timestamp: fixedTime, TripId: fixedUUID},
		},
		{
			name: "TePacket",
			packet: &TePacket{
				Timestamp:        fixedTime,
				TripId:           fixedUUID,
				DistanceTraveled: 100.0,
				MaxSpeed:         50.0,
				Duration:         time.Hour,
			},
		},
		{
			name:   "ImPacket",
			packet: &ImPacket{Timestamp: fixedTime, ChatId: fixedUUID, Message: "test"},
		},
		{
			name:        "invalid packet type",
			packet:      "invalid",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parsePacketToString(tt.packet)

			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("parsePacketToString failed: %v", err)
			}
			if result == nil {
				t.Fatal("result is nil")
			}
			if len(*result) == 0 {
				t.Fatal("result is empty string")
			}
		})
	}
}

// TestParseArgs tests the parseArgs function with various data types and key remappings
func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		validate func(*map[string]any) error
	}{
		{
			name:  "empty args",
			input: "",
			validate: func(args *map[string]any) error {
				if len(*args) != 0 {
					return fmt.Errorf("expected empty map, got %d items", len(*args))
				}
				return nil
			},
		},
		{
			name:  "integer value",
			input: "count:42",
			validate: func(args *map[string]any) error {
				val, ok := (*args)["count"]
				if !ok {
					return fmt.Errorf("count key not found")
				}
				intVal, ok := val.(int)
				if !ok || intVal != 42 {
					return fmt.Errorf("expected int 42, got %v (type %T)", val, val)
				}
				return nil
			},
		},
		{
			name:  "negative integer",
			input: "temperature:-15",
			validate: func(args *map[string]any) error {
				val, ok := (*args)["temperature"]
				if !ok {
					return fmt.Errorf("temperature key not found")
				}
				intVal, ok := val.(int)
				if !ok || intVal != -15 {
					return fmt.Errorf("expected int -15, got %v", val)
				}
				return nil
			},
		},
		{
			name:  "float value",
			input: "latitude:10.123456",
			validate: func(args *map[string]any) error {
				val, ok := (*args)["latitude"]
				if !ok {
					return fmt.Errorf("latitude key not found")
				}
				floatVal, ok := val.(float64)
				if !ok {
					return fmt.Errorf("expected float64, got %T", val)
				}
				if floatVal != 10.123456 {
					return fmt.Errorf("expected 10.123456, got %v", floatVal)
				}
				return nil
			},
		},
		{
			name:  "boolean true",
			input: "enabled:true",
			validate: func(args *map[string]any) error {
				val, ok := (*args)["enabled"]
				if !ok {
					return fmt.Errorf("enabled key not found")
				}
				boolVal, ok := val.(bool)
				if !ok || !boolVal {
					return fmt.Errorf("expected bool true, got %v", val)
				}
				return nil
			},
		},
		{
			name:  "boolean false",
			input: "disabled:false",
			validate: func(args *map[string]any) error {
				val, ok := (*args)["disabled"]
				if !ok {
					return fmt.Errorf("disabled key not found")
				}
				boolVal, ok := val.(bool)
				if !ok || boolVal {
					return fmt.Errorf("expected bool false, got %v", val)
				}
				return nil
			},
		},
		{
			name:  "string value",
			input: "name:hello",
			validate: func(args *map[string]any) error {
				val, ok := (*args)["name"]
				if !ok {
					return fmt.Errorf("name key not found")
				}
				strVal, ok := val.(string)
				if !ok || strVal != "hello" {
					return fmt.Errorf("expected string hello, got %v", val)
				}
				return nil
			},
		},
		{
			name:  "GPIO digital input remapping",
			input: "io3.di:1",
			validate: func(args *map[string]any) error {
				val, ok := (*args)["gpio.3.digital.input"]
				if !ok {
					return fmt.Errorf("gpio.3.digital.input key not found, have keys: %v", getKeys(*args))
				}
				intVal, ok := val.(int)
				if !ok || intVal != 1 {
					return fmt.Errorf("expected int 1, got %v", val)
				}
				return nil
			},
		},
		{
			name:  "GPIO digital output remapping",
			input: "io5.do:0",
			validate: func(args *map[string]any) error {
				_, ok := (*args)["gpio.5.digital.output"]
				if !ok {
					return fmt.Errorf("gpio.5.digital.output key not found")
				}
				return nil
			},
		},
		{
			name:  "GPIO analog input remapping",
			input: "io2.ai:512",
			validate: func(args *map[string]any) error {
				_, ok := (*args)["gpio.2.analog.input"]
				if !ok {
					return fmt.Errorf("gpio.2.analog.input key not found")
				}
				return nil
			},
		},
		{
			name:  "GPIO counter remapping",
			input: "io1.counter:999",
			validate: func(args *map[string]any) error {
				_, ok := (*args)["gpio.1.event.count"]
				if !ok {
					return fmt.Errorf("gpio.1.event.count key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE temperature celsius remapping",
			input: "ble.0.tempc:25",
			validate: func(args *map[string]any) error {
				val, ok := (*args)["ble.0.temperature.celsius"]
				if !ok {
					return fmt.Errorf("ble.0.temperature.celsius key not found, have: %v", getKeys(*args))
				}
				intVal, ok := val.(int)
				if !ok || intVal != 25 {
					return fmt.Errorf("expected int 25, got %v", val)
				}
				return nil
			},
		},
		{
			name:  "BLE humidity remapping",
			input: "ble.1.hum:60.5",
			validate: func(args *map[string]any) error {
				_, ok := (*args)["ble.1.humidity"]
				if !ok {
					return fmt.Errorf("ble.1.humidity key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE id (mac address) remapping",
			input: "ble.2.id:AA:BB:CC:DD:EE:FF",
			validate: func(args *map[string]any) error {
				val, ok := (*args)["ble.2.mac.address"]
				if !ok {
					return fmt.Errorf("ble.2.mac.address key not found")
				}
				strVal, ok := val.(string)
				if !ok || strVal != "AA:BB:CC:DD:EE:FF" {
					return fmt.Errorf("expected string, got %v", val)
				}
				return nil
			},
		},
		{
			name:  "multiple args",
			input: "count:42,enabled:true,name:test,temperature:23.5",
			validate: func(args *map[string]any) error {
				if len(*args) != 4 {
					return fmt.Errorf("expected 4 args, got %d", len(*args))
				}
				return nil
			},
		},
		{
			name:  "report code remapping",
			input: "report:123",
			validate: func(args *map[string]any) error {
				_, ok := (*args)["report.code"]
				if !ok {
					return fmt.Errorf("report.code key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE temperature fahrenheit remapping",
			input: "ble.0.tempf:77",
			validate: func(args *map[string]any) error {
				_, ok := (*args)["ble.0.temperature.fahrenheit"]
				if !ok {
					return fmt.Errorf("ble.0.temperature.fahrenheit key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE model_id remapping",
			input: "ble.0.model_id:42",
			validate: func(args *map[string]any) error {
				_, ok := (*args)["ble.0.model.id"]
				if !ok {
					return fmt.Errorf("ble.0.model.id key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE battery remapping",
			input: "ble.0.batt:85",
			validate: func(args *map[string]any) error {
				_, ok := (*args)["ble.0.battery.level"]
				if !ok {
					return fmt.Errorf("ble.0.battery.level key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE lux remapping",
			input: "ble.0.lux:500",
			validate: func(args *map[string]any) error {
				_, ok := (*args)["ble.0.light.level.lux"]
				if !ok {
					return fmt.Errorf("ble.0.light.level.lux key not found")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseArgs(tt.input)
			if result == nil {
				t.Fatal("parseArgs returned nil")
			}

			if err := tt.validate(result); err != nil {
				t.Errorf("validation failed: %v", err)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string  { return &s }
func intPtr(i int) *int           { return &i }
func floatPtr(f float64) *float64 { return &f }

func getKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
