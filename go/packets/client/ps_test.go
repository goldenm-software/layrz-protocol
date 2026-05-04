package client_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
)

func TestPs_FromPacket_Errors(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{"wrong envelope", "<Xx>garbage</Xx>"},
		{"bad CRC hex", "<Ps>1700000000;key:val;ZZZZ</Ps>"},
		{"CRC mismatch", "<Ps>1700000000;key:val;0000</Ps>"},
		{"wrong part count", "<Ps>1700000000;0000</Ps>"},
		{"bad timestamp", "<Ps>notanint;key:val;"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw := tc.raw
			if err := (&client.PsPacket{}).FromPacket(&raw); err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

func TestPs_ToPacket_DefaultParamType(t *testing.T) {
	packet := client.PsPacket{
		Timestamp: fixedTime,
		Params:    map[string]any{"key": []byte("unknown")},
	}
	out := packet.ToPacket()
	if out == nil {
		t.Fatal("ToPacket returned nil")
	}
}

func TestPs_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]any
	}{
		{
			name:   "with parameters",
			params: map[string]any{"key": "value"},
		},
		{
			name:   "empty parameters",
			params: map[string]any{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := client.PsPacket{
				Timestamp: fixedTime,
				Params:    tt.params,
			}
			encoded := *packet.ToPacket()

			raw := encoded
			decoded := client.PsPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}
			if *decoded.ToPacket() != encoded {
				t.Errorf("round-trip mismatch")
			}
		})
	}
}
