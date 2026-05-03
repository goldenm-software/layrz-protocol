package trips_test

import (
	"testing"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/trips"
)

func TestTe_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name             string
		timestamp        int64
		tripId           string
		distanceTraveled float64
		maxSpeed         float64
		duration         time.Duration
	}{
		{
			name:             "complete trip",
			timestamp:        1700000000,
			tripId:           "trip-uuid-001",
			distanceTraveled: 100.5,
			maxSpeed:         50.25,
			duration:         time.Hour,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := time.Unix(tt.timestamp, 0)
			packet := trips.TePacket{
				Timestamp:        ts,
				TripId:           tt.tripId,
				DistanceTraveled: tt.distanceTraveled,
				MaxSpeed:         tt.maxSpeed,
				Duration:         tt.duration,
			}
			encoded := *packet.ToPacket()

			raw := encoded
			decoded := trips.TePacket{}
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
