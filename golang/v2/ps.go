package layrzprotocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Based on the `Layrz Protocol v2` specification. PsPacket is the Settings packet
// sent from the device to the server
type PsPacket struct {
	// Is the timestamp of the response packet
	Timestamp time.Time

	// Is the current configuration of the device
	Params *map[string]interface{}
}

// FromPacket is a method that converts a raw packet to a PaPacket
// based on the `Layrz Protocol v2` specification
//
// Returns an error if the packet is invalid
func (p *PsPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Ps>") || !strings.HasSuffix(*raw, "</Ps>") {
		return errors.New("invalid package, should be <Ps>...</Ps>")
	}

	*raw = strings.TrimPrefix(*raw, "<Ps>")
	*raw = strings.TrimSuffix(*raw, "</Ps>")

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
	if len(parts) != 2 {
		return errors.New("invalid package, should contain 2 parts")
	}

	rawTimestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return errors.New("cannot convert timestamp to integer")
	}

	p.Timestamp = time.Unix(rawTimestamp, 0)
	p.Params = parseArgs(parts[1])

	return nil
}

// ToPacket is a method that converts a PaPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *PsPacket) ToPacket() *string {
	content := ""

	params := make([]string, 0)
	for key, value := range *p.Params {
		switch v := value.(type) {
		case string:
			params = append(params, fmt.Sprintf("%s:%s", key, v))
		case int:
			params = append(params, fmt.Sprintf("%s:%d", key, v))
		case float64:
			params = append(params, fmt.Sprintf("%s:%f", key, v))
		case bool:
			params = append(params, fmt.Sprintf("%s:%t", key, v))
		default:
			params = append(params, fmt.Sprintf("%s:%s", key, v))
		}
	}

	content = fmt.Sprintf("%d;%s;", p.Timestamp.Unix(), strings.Join(params, ","))
	crc := calculateCrc([]byte(content))
	content = fmt.Sprintf("<Ps>%s%04X</Ps>", content, crc)
	return &content
}
