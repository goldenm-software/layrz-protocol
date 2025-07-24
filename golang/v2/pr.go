package layrzprotocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Based on the `Layrz Protocol v2` specification. PrPacket is the syncronization packet
// sent from the device to the server
type PrPacket struct{}

// FromPacket is a method that converts a raw packet to a PrPacket
// based on the `Layrz Protocol v2` specification
//
// Returns an error if the packet is invalid
func (p *PrPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Pr>") || !strings.HasSuffix(*raw, "</Pr>") {
		return errors.New("invalid package, should be <Pr>...</Pr>")
	}

	*raw = strings.TrimPrefix(*raw, "<Pr>")
	*raw = strings.TrimSuffix(*raw, "</Pr>")

	parts := strings.Split(*raw, ";")
	if len(parts) != 2 {
		return errors.New("invalid package, should contain 2 parts")
	}

	receivedCrc, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return errors.New("cannot convert CRC to integer")
	}

	calculatedCrc := calculateCrc([]byte(";"))

	if calculatedCrc != uint16(receivedCrc) {
		return fmt.Errorf("invalid CRC, received: %04X, calculated: %04X", receivedCrc, calculatedCrc)
	}

	return nil
}

// ToPacket is a method that converts a PaPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *PrPacket) ToPacket() *string {
	crc := calculateCrc([]byte(";"))
	content := fmt.Sprintf("<Pr>;%04X</Pr>", crc)
	return &content
}
