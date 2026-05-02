package layrzprotocol

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Based on the `Layrz Protocol v2` specification. PmPacket is the media packet
// sent from the device to the server
type PmPacket struct {
	// Is the timestamp of the response packet
	Filename *string

	// Is the command id of the response packet
	ContentType *string

	// Is the message of the response packet
	Data *[]byte
}

// FromPacket is a method that converts a raw packet to a PaPacket
// based on the `Layrz Protocol v2` specification
//
// Returns an error if the packet is invalid
func (p *PmPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Pm>") || !strings.HasSuffix(*raw, "</Pm>") {
		return errors.New("invalid package, should be <Pm>...</Pm>")
	}

	*raw = strings.TrimPrefix(*raw, "<Pm>")
	*raw = strings.TrimSuffix(*raw, "</Pm>")

	rawCrc := (*raw)[len(*raw)-4:]
	*raw = (*raw)[:len(*raw)-4]

	receivedCrc, err := strconv.ParseUint(rawCrc, 16, 16)
	if err != nil {
		return errors.New("cannot convert CRC to integer")
	}

	calculatedCrc := calculateCrc([]byte(*raw))

	if calculatedCrc != uint16(receivedCrc) {
		return fmt.Errorf("invalid CRC, received: %04X, calculated: %04X", receivedCrc, calculatedCrc)
	}

	*raw = strings.TrimSuffix(*raw, ";")
	parts := strings.Split(*raw, ";")
	if len(parts) != 3 {
		return errors.New("invalid package, should contain 3 parts")
	}

	p.Filename = &parts[0]
	p.ContentType = &parts[1]

	data, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return errors.New("cannot decode base64 data")
	}
	p.Data = &data

	return nil
}

// ToPacket is a method that converts a PaPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *PmPacket) ToPacket() *string {
	content := ""
	content += *p.Filename + ";"
	content += *p.ContentType + ";"
	content += base64.StdEncoding.EncodeToString(*p.Data) + ";"

	crc := calculateCrc([]byte(content))
	content += fmt.Sprintf("%04X", crc)

	content = fmt.Sprintf("<Pm>%s</Pm>", content)
	return &content
}
