package trips

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/internal/wire"
)

// TsPacket is the Trip Start packet sent between Layrz services to identify trips.
type TsPacket struct {
	// Timestamp of the packet
	Timestamp time.Time

	// TripId is the unique trip identifier (UUID string)
	TripId string
}

// FromPacket converts a raw <Ts>...</Ts> string to a TsPacket.
func (p *TsPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Ts>") || !strings.HasSuffix(*raw, "</Ts>") {
		return errors.New("invalid packet, should be <Ts>...</Ts>")
	}

	body := (*raw)[4 : len(*raw)-5]

	rawCrc := body[len(body)-4:]
	body = body[:len(body)-4]

	receivedCrc, err := strconv.ParseUint(rawCrc, 16, 16)
	if err != nil {
		return errors.New("cannot convert CRC to integer")
	}

	calculatedCrc := wire.Calculate([]byte(body))
	if calculatedCrc != uint16(receivedCrc) {
		return fmt.Errorf("invalid CRC, received: %04X, calculated: %04X", receivedCrc, calculatedCrc)
	}

	body = strings.TrimSuffix(body, ";")
	parts := strings.Split(body, ";")
	if len(parts) != 2 {
		return errors.New("invalid packet, should have 2 parts")
	}

	rawTimestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return errors.New("cannot convert timestamp to integer")
	}
	p.Timestamp = time.Unix(rawTimestamp, 0)
	p.TripId = parts[1]

	return nil
}

// ToPacket converts a TsPacket to its wire representation.
func (p *TsPacket) ToPacket() *string {
	content := fmt.Sprintf("%d;%s;", p.Timestamp.Unix(), p.TripId)
	crc := wire.Calculate([]byte(content))
	result := fmt.Sprintf("<Ts>%s%04X</Ts>", content, crc)
	return &result
}
