package ai_test

import (
	"testing"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/ai"
)

func TestIm_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		chatId    string
		message   string
	}{
		{
			name:      "simple message",
			timestamp: 1700000000,
			chatId:    "uuid-1234",
			message:   "hello world",
		},
		{
			name:      "message with semicolons",
			timestamp: 1700000001,
			chatId:    "uuid-5678",
			message:   "a;b;c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := time.Unix(tt.timestamp, 0)
			packet := ai.ImPacket{Timestamp: ts, ChatId: tt.chatId, Message: tt.message}
			encoded := *packet.ToPacket()

			raw := encoded
			decoded := ai.ImPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}
			if decoded.Message != tt.message {
				t.Errorf("message mismatch: got %s, want %s", decoded.Message, tt.message)
			}
			if *decoded.ToPacket() != encoded {
				t.Errorf("round-trip mismatch")
			}
		})
	}
}
