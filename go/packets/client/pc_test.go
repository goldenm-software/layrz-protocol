package client_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
)

func TestPc_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		commandId int
		message   string
	}{
		{
			name:      "valid response",
			timestamp: 1700000000,
			commandId: 1,
			message:   "OK",
		},
		{
			name:      "error message",
			timestamp: 1700000001,
			commandId: 2,
			message:   "ERROR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := client.PcPacket{
				Timestamp: fixedTime,
				CommandId: tt.commandId,
				Message:   stringPtr(tt.message),
			}
			encoded := *packet.ToPacket()

			raw := encoded
			decoded := client.PcPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}
			if decoded.CommandId != tt.commandId {
				t.Errorf("commandId mismatch: got %d, want %d", decoded.CommandId, tt.commandId)
			}
			if *decoded.Message != tt.message {
				t.Errorf("message mismatch: got %s, want %s", *decoded.Message, tt.message)
			}
			if *decoded.ToPacket() != encoded {
				t.Errorf("round-trip mismatch")
			}
		})
	}
}
