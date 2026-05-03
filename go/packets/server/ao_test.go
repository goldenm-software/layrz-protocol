package server_test

import (
	"testing"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
)

func TestAo_FromPacket_ToPacket(t *testing.T) {
	ts := time.Unix(1700000000, 0)
	packet := server.AoPacket{Timestamp: ts}
	encoded := *packet.ToPacket()

	raw := encoded
	decoded := server.AoPacket{}
	if err := decoded.FromPacket(&raw); err != nil {
		t.Fatalf("FromPacket failed: %v", err)
	}
	if decoded.Timestamp.Unix() != ts.Unix() {
		t.Errorf("timestamp mismatch: got %d, want %d", decoded.Timestamp.Unix(), ts.Unix())
	}
	if *decoded.ToPacket() != encoded {
		t.Errorf("round-trip mismatch")
	}
}
