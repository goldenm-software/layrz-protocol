package server

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/goldenm-software/layrz-protocol/go/v3/internal/wire"
)

// Based on the `Layrz Protocol v2` specification. AuPacket is the authorization request packet
// sent from the server to the device
//
// Deprecated: This packet is deprecated and will be removed in v4.0
type AuPacket struct{}

// FromPacket is a method that converts a raw packet to a AuPacket
// based on the `Layrz Protocol v2` specification
//
// Returns an error if the packet is invalid
func (p *AuPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Au>") || !strings.HasSuffix(*raw, "</Au>") {
		return errors.New("invalid package, should be <Au>...</Au>")
	}

	*raw = strings.TrimPrefix(*raw, "<Au>")
	*raw = strings.TrimSuffix(*raw, "</Au>")

	parts := strings.Split(*raw, ";")
	if len(parts) != 2 {
		return errors.New("invalid package, should contain 2 parts")
	}

	receivedCrc, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return errors.New("cannot convert CRC to integer")
	}

	calculatedCrc := wire.Calculate([]byte(";"))

	if calculatedCrc != uint16(receivedCrc) {
		return fmt.Errorf("invalid CRC, received: %04X, calculated: %04X", receivedCrc, calculatedCrc)
	}

	return nil
}

// ToPacket is a method that converts a AuPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *AuPacket) ToPacket() *string {
	crc := wire.Calculate([]byte(";"))
	content := fmt.Sprintf("<Au>;%04X</Au>", crc)
	return &content
}
