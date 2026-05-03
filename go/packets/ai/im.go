package ai

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/internal/wire"
)

// ImPacket is the AI message packet.
type ImPacket struct {
	// Timestamp of the packet
	Timestamp time.Time

	// ChatId is the unique chat identifier (UUID string)
	ChatId string

	// Message is the chat message content; semicolons are escaped as |||
	Message string
}

// FromPacket converts a raw <Im>...</Im> string to an ImPacket.
func (p *ImPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Im>") || !strings.HasSuffix(*raw, "</Im>") {
		return errors.New("invalid packet, should be <Im>...</Im>")
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
	if len(parts) != 3 {
		return errors.New("invalid packet, should have 3 parts")
	}

	rawTimestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return errors.New("cannot convert timestamp to integer")
	}
	p.Timestamp = time.Unix(rawTimestamp, 0)
	p.ChatId = parts[1]
	p.Message = strings.ReplaceAll(parts[2], "|||", ";")

	return nil
}

// ToPacket converts an ImPacket to its wire representation.
func (p *ImPacket) ToPacket() *string {
	escapedMessage := strings.ReplaceAll(p.Message, ";", "|||")
	content := fmt.Sprintf("%d;%s;%s;", p.Timestamp.Unix(), p.ChatId, escapedMessage)
	crc := wire.Calculate([]byte(content))
	result := fmt.Sprintf("<Im>%s%04X</Im>", content, crc)
	return &result
}
