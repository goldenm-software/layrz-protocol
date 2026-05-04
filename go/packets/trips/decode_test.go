package trips_test

import (
	"testing"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/trips"
)

func TestDecode_Ts(t *testing.T) {
	packet := trips.TsPacket{
		Timestamp: time.Unix(1700000000, 0),
		TripId:    "trip-001",
	}
	encoded := *packet.ToPacket()

	result, err := trips.Decode([]byte(encoded))
	if err != nil {
		t.Fatalf("Decode(Ts) failed: %v", err)
	}
	if result == nil {
		t.Fatal("Decode returned nil")
	}
	if *result.ToPacket() != encoded {
		t.Errorf("round-trip mismatch: got %s, want %s", *result.ToPacket(), encoded)
	}
}

func TestDecode_Te(t *testing.T) {
	packet := trips.TePacket{
		Timestamp:        time.Unix(1700000000, 0),
		TripId:           "trip-001",
		DistanceTraveled: 100.5,
		MaxSpeed:         80.0,
		Duration:         time.Hour,
	}
	encoded := *packet.ToPacket()

	result, err := trips.Decode([]byte(encoded))
	if err != nil {
		t.Fatalf("Decode(Te) failed: %v", err)
	}
	if result == nil {
		t.Fatal("Decode returned nil")
	}
	if *result.ToPacket() != encoded {
		t.Errorf("round-trip mismatch: got %s, want %s", *result.ToPacket(), encoded)
	}
}

func TestDecode_UnknownTag(t *testing.T) {
	_, err := trips.Decode([]byte("<Xx>garbage</Xx>"))
	if err == nil {
		t.Error("expected error for unknown tag")
	}
}

func TestDecode_Empty(t *testing.T) {
	_, err := trips.Decode([]byte(""))
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestDecode_MalformedTs(t *testing.T) {
	_, err := trips.Decode([]byte("<Ts>notvalid</Ts>"))
	if err == nil {
		t.Error("expected error for malformed Ts body")
	}
}

func TestDecode_MalformedTe(t *testing.T) {
	_, err := trips.Decode([]byte("<Te>notvalid</Te>"))
	if err == nil {
		t.Error("expected error for malformed Te body")
	}
}
