package layrzprotocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Based on the `Layrz Protocol v2` specification. PcPacket is the Packet command response
// sent from the device to the server
type PcPacket struct {
	// Is the timestamp of the response packet
	Timestamp time.Time

	// Is the command id of the response packet
	CommandId int

	// Is the message of the response packet
	Message *string
}

// FromPacket is a method that converts a raw packet to a PaPacket
// based on the `Layrz Protocol v2` specification
//
// Returns an error if the packet is invalid
func (p *PcPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Pc>") || !strings.HasSuffix(*raw, "</Pc>") {
		return errors.New("invalid package, should be <Pc>...</Pc>")
	}

	*raw = strings.TrimPrefix(*raw, "<Pc>")
	*raw = strings.TrimSuffix(*raw, "</Pc>")

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

	rawTimestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return errors.New("cannot convert timestamp to integer")
	}

	p.Timestamp = time.Unix(rawTimestamp, 0)
	p.Message = &parts[2]

	p.CommandId, err = strconv.Atoi(parts[1])
	if err != nil {
		return errors.New("cannot convert command id to integer")
	}

	return nil
}

// ToPacket is a method that converts a PaPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *PcPacket) ToPacket() *string {
	content := fmt.Sprintf("%d;%d;%s;", p.Timestamp.Unix(), p.CommandId, *p.Message)
	crc := calculateCrc([]byte(content))
	content = fmt.Sprintf("<Pc>%s%04X</Pc>", content, crc)
	return &content
}
