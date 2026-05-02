package layrzprotocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TePacket is the Trip End packet sent between Layrz services to identify trips.
type TePacket struct {
	// Timestamp of the packet
	Timestamp time.Time

	// TripId is the unique trip identifier (UUID string)
	TripId string

	// DistanceTraveled is the distance in meters
	DistanceTraveled float64

	// MaxSpeed is the maximum speed in km/h
	MaxSpeed float64

	// Duration is the trip duration in seconds
	Duration time.Duration
}

// FromPacket converts a raw <Te>...</Te> string to a TePacket.
func (p *TePacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Te>") || !strings.HasSuffix(*raw, "</Te>") {
		return errors.New("invalid packet, should be <Te>...</Te>")
	}

	body := (*raw)[4 : len(*raw)-5]

	rawCrc := body[len(body)-4:]
	body = body[:len(body)-4]

	receivedCrc, err := strconv.ParseUint(rawCrc, 16, 16)
	if err != nil {
		return errors.New("cannot convert CRC to integer")
	}

	calculatedCrc := calculateCrc([]byte(body))
	if calculatedCrc != uint16(receivedCrc) {
		return fmt.Errorf("invalid CRC, received: %04X, calculated: %04X", receivedCrc, calculatedCrc)
	}

	body = strings.TrimSuffix(body, ";")
	parts := strings.Split(body, ";")
	if len(parts) != 5 {
		return errors.New("invalid packet, should have 5 parts")
	}

	rawTimestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return errors.New("cannot convert timestamp to integer")
	}
	p.Timestamp = time.Unix(rawTimestamp, 0)
	p.TripId = parts[1]

	distanceTraveled, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return errors.New("cannot convert distance_traveled to float")
	}
	p.DistanceTraveled = distanceTraveled

	maxSpeed, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return errors.New("cannot convert max_speed to float")
	}
	p.MaxSpeed = maxSpeed

	durationSecs, err := strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		return errors.New("cannot convert duration to integer")
	}
	p.Duration = time.Duration(durationSecs) * time.Second

	return nil
}

// ToPacket converts a TePacket to its wire representation.
func (p *TePacket) ToPacket() *string {
	content := fmt.Sprintf("%d;%s;%.3f;%.3f;%d;",
		p.Timestamp.Unix(),
		p.TripId,
		p.DistanceTraveled,
		p.MaxSpeed,
		int64(p.Duration.Seconds()),
	)
	crc := calculateCrc([]byte(content))
	result := fmt.Sprintf("<Te>%s%04X</Te>", content, crc)
	return &result
}
