package layrzprotocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Based on the `Layrz Protocol v2` specification. ArPacket is the error packet
// sent from the server to the device
type ArPacket struct {
	// Is the reason of the error
	Reason string
}

// FromPacket is a method that converts a raw packet to a PaPacket
// based on the `Layrz Protocol v2` specification
//
// Returns an error if the packet is invalid
func (p *ArPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Ar>") || !strings.HasSuffix(*raw, "</Ar>") {
		return errors.New("invalid package, should be <Ar>...</Ar>")
	}

	*raw = strings.TrimPrefix(*raw, "<Ar>")
	*raw = strings.TrimSuffix(*raw, "</Ar>")

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

	p.Reason = parts[0]
	return nil
}

// ToPacket is a method that converts a PaPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *ArPacket) ToPacket() *string {
	content := ""
	content += p.Reason + ";"
	crc := calculateCrc([]byte(content))
	content += fmt.Sprintf("%04X", crc)

	content = fmt.Sprintf("<Ar>%s</Ar>", content)
	return &content

}
