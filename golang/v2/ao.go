package layrzprotocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Based on the `Layrz Protocol v2` specification. AoPacket is the success packet
// sent from the server to the device
type AoPacket struct {
	// Is the timestamp of the response packet
	Timestamp time.Time
}

// FromPacket is a method that converts a raw packet to a PaPacket
// based on the `Layrz Protocol v2` specification
//
// Returns an error if the packet is invalid
func (p *AoPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Ao>") || !strings.HasSuffix(*raw, "</Ao>") {
		return errors.New("invalid package, should be <Ao>...</Ao>")
	}

	*raw = strings.TrimPrefix(*raw, "<Ao>")
	*raw = strings.TrimSuffix(*raw, "</Ao>")

	parts := strings.Split(*raw, ";")
	if len(parts) != 2 {
		return errors.New("invalid package, should have 2 parts")
	}

	receivedCrc, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return errors.New("cannot convert CRC to integer")
	}

	calculatedCrc := calculateCrc([]byte(fmt.Sprintf("%s;", parts[0])))

	if calculatedCrc != uint16(receivedCrc) {
		return fmt.Errorf("invalid CRC, received: %04X, calculated: %04X", receivedCrc, calculatedCrc)
	}

	rawTimestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return errors.New("cannot convert timestamp to integer")
	}

	p.Timestamp = time.Unix(rawTimestamp, 0)

	return nil
}

// ToPacket is a method that converts a PaPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *AoPacket) ToPacket() *string {
	content := ""
	content += strconv.FormatInt(p.Timestamp.Unix(), 10) + ";"

	crc := calculateCrc([]byte(content))
	content += fmt.Sprintf("%04X", crc)

	content = fmt.Sprintf("<Ao>%s</Ao>", content)
	return &content
}
