package client_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
)

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
			data:        []byte("hello"),
		},
		{
			name:        "binary file",
			filename:    "test.bin",
			contentType: "application/octet-stream",
			data:        []byte{0x01, 0x02, 0x03},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := client.PmPacket{
				Filename:    stringPtr(tt.filename),
				ContentType: stringPtr(tt.contentType),
				Data:        &tt.data,
			}
			encoded := *packet.ToPacket()

			raw := encoded
			decoded := client.PmPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}
			if *decoded.Filename != tt.filename {
				t.Errorf("filename mismatch: got %s, want %s", *decoded.Filename, tt.filename)
			}
			if *decoded.ToPacket() != encoded {
				t.Errorf("round-trip mismatch")
			}
		})
	}
}
