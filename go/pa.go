package layrzprotocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Based on the `Layrz Protocol v2` specification. PaPacket is the Authentication packet
// sent from the device to the server.
type PaPacket struct {
	// Defines the ident of the device
	Ident *string

	// Defines the password of the device
	Password *string
}

// FromPacket is a method that converts a raw packet to a PaPacket
// based on the `Layrz Protocol v2` specification
//
// Returns an error if the packet is invalid
func (p *PaPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Pa>") || !strings.HasSuffix(*raw, "</Pa>") {
		return errors.New("invalid package, should be <Pa>...</Pa>")
	}

	*raw = strings.TrimPrefix(*raw, "<Pa>")
	*raw = strings.TrimSuffix(*raw, "</Pa>")

	parts := strings.Split(*raw, ";")

	if len(parts) != 3 {
		return errors.New("invalid package, should contain 3 parts")
	}

	receivedCrc, err := strconv.ParseUint(parts[2], 16, 16)

	if err != nil {
		return errors.New("cannot convert CRC to integer")
	}

	content := fmt.Sprintf("%s;%s;", parts[0], parts[1])
	calculatedCrc := calculateCrc([]byte(content))

	if calculatedCrc != uint16(receivedCrc) {
		return fmt.Errorf("invalid CRC, received: %04X, calculated: %04X", receivedCrc, calculatedCrc)
	}

	p.Ident = &parts[0]
	p.Password = &parts[1]

	return nil
}

// ToPacket is a method that converts a PaPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *PaPacket) ToPacket() *string {
	content := fmt.Sprintf("%s;%s;", (*p.Ident), (*p.Password))
	crc := calculateCrc([]byte(content))
	content += fmt.Sprintf("%04X", crc)

	content = fmt.Sprintf("<Pa>%s</Pa>", content)
	return &content
}
