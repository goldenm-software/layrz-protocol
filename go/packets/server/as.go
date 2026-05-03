package server

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/goldenm-software/layrz-protocol/go/v3/internal/wire"
)

// Based on the `Layrz Protocol v2` specification. AsPacket is the success authentication packet
// sent from the server to the device
type AsPacket struct{}

// FromPacket is a method that converts a raw packet to a AsPacket
// based on the `Layrz Protocol v2` specification
//
// Returns an error if the packet is invalid
func (p *AsPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<As>") || !strings.HasSuffix(*raw, "</As>") {
		return errors.New("invalid package, should be <As>...</As>")
	}

	*raw = strings.TrimPrefix(*raw, "<As>")
	*raw = strings.TrimSuffix(*raw, "</As>")

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

// ToPacket is a method that converts a AsPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *AsPacket) ToPacket() *string {
	crc := wire.Calculate([]byte(";"))
	content := fmt.Sprintf("<As>;%04X</As>", crc)
	return &content
}
