package ai_test

import (
	"testing"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/ai"
)

func TestDecode_Im(t *testing.T) {
	packet := ai.ImPacket{
		Timestamp: time.Unix(1700000000, 0),
		ChatId:    "uuid-abc",
		Message:   "hello",
	}
	encoded := *packet.ToPacket()

	result, err := ai.Decode([]byte(encoded))
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if result == nil {
		t.Fatal("Decode returned nil")
	}
	if *result.ToPacket() != encoded {
		t.Errorf("round-trip mismatch: got %s, want %s", *result.ToPacket(), encoded)
	}
}

func TestDecode_UnknownTag(t *testing.T) {
	_, err := ai.Decode([]byte("<Xx>garbage</Xx>"))
	if err == nil {
		t.Error("expected error for unknown tag")
	}
}

func TestDecode_Empty(t *testing.T) {
	_, err := ai.Decode([]byte(""))
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestDecode_MalformedIm(t *testing.T) {
	_, err := ai.Decode([]byte("<Im>notvalid</Im>"))
	if err == nil {
		t.Error("expected error for malformed Im body")
	}
}
