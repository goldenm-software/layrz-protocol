package trips_test

import (
	"testing"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/trips"
)

func TestTs_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		tripId    string
	}{
		{
			name:      "trip start",
			timestamp: 1700000000,
			tripId:    "trip-uuid-001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := time.Unix(tt.timestamp, 0)
			packet := trips.TsPacket{Timestamp: ts, TripId: tt.tripId}
			encoded := *packet.ToPacket()

			raw := encoded
			decoded := trips.TsPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}
			if decoded.TripId != tt.tripId {
				t.Errorf("TripId mismatch: got %s, want %s", decoded.TripId, tt.tripId)
			}
			if *decoded.ToPacket() != encoded {
				t.Errorf("round-trip mismatch")
			}
		})
	}
}
