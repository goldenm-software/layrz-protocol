package layrzprotocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Based on the `Layrz Protocol v2` specification. PdPacket is the Data packet
// sent from the device to the server
type PdPacket struct {
	// Is the timestamp of the response packet
	Timestamp time.Time

	// Is the position of the device
	Position *Position

	// Is the extra data sent by the device
	ExtraData *map[string]interface{}
}

// FromPacket is a method that converts a raw packet to a PaPacket
// based on the `Layrz Protocol v2` specification
//
// Returns an error if the packet is invalid
func (p *PdPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Pd>") || !strings.HasSuffix(*raw, "</Pd>") {
		return errors.New("invalid package, should be <Pd>...</Pd>")
	}

	*raw = strings.TrimPrefix(*raw, "<Pd>")
	*raw = strings.TrimSuffix(*raw, "</Pd>")

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
	if len(parts) != 9 {
		return errors.New("invalid package, should contain 9 parts")
	}

	rawTimestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return errors.New("cannot convert timestamp to integer")
	}
	p.Timestamp = time.Unix(rawTimestamp, 0)

	var latitude, longitude, altitude, speed, direction, hdop *float64
	var satelliteCount *int

	if parts[1] != "" {
		converted, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return errors.New("cannot convert latitude to float")
		}
		latitude = &converted
	}

	if parts[2] != "" {
		converted, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return errors.New("cannot convert longitude to float")
		}
		longitude = &converted
	}

	if parts[3] != "" {
		converted, err := strconv.ParseFloat(parts[3], 64)
		if err != nil {
			return errors.New("cannot convert altitude to float")
		}
		altitude = &converted
	}

	if parts[4] != "" {
		converted, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			return errors.New("cannot convert speed to float")
		}
		speed = &converted
	}

	if parts[5] != "" {
		converted, err := strconv.ParseFloat(parts[5], 64)
		if err != nil {
			return errors.New("cannot convert direction to float")
		}
		direction = &converted
	}

	if parts[6] != "" {
		converted, err := strconv.Atoi(parts[6])
		if err != nil {
			return errors.New("cannot convert satellite count to integer")
		}
		satelliteCount = &converted
	}

	if parts[7] != "" {
		converted, err := strconv.ParseFloat(parts[7], 64)
		if err != nil {
			return errors.New("cannot convert hdop to float")
		}
		hdop = &converted
	}

	p.Position = &Position{
		Latitude:       latitude,
		Longitude:      longitude,
		Altitude:       altitude,
		Speed:          speed,
		Direction:      direction,
		SatelliteCount: satelliteCount,
		Hdop:           hdop,
	}

	p.ExtraData = parseArgs(parts[8])

	return nil
}

// ToPacket is a method that converts a PaPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *PdPacket) ToPacket() *string {
	content := ""
	content += fmt.Sprintf("%d;", p.Timestamp.Unix())
	if p.Position != nil {
		if p.Position.Latitude != nil {
			content += fmt.Sprintf("%f;", *p.Position.Latitude)
		} else {
			content += ";"
		}
		if p.Position.Longitude != nil {
			content += fmt.Sprintf("%f;", *p.Position.Longitude)
		} else {
			content += ";"
		}
		if p.Position.Altitude != nil {
			content += fmt.Sprintf("%f;", *p.Position.Altitude)
		} else {
			content += ";"
		}
		if p.Position.Speed != nil {
			content += fmt.Sprintf("%f;", *p.Position.Speed)
		} else {
			content += ";"
		}
		if p.Position.Direction != nil {
			content += fmt.Sprintf("%f;", *p.Position.Direction)
		} else {
			content += ";"
		}
		if p.Position.SatelliteCount != nil {
			content += fmt.Sprintf("%d;", *p.Position.SatelliteCount)
		} else {
			content += ";"
		}
		if p.Position.Hdop != nil {
			content += fmt.Sprintf("%f;", *p.Position.Hdop)
		} else {
			content += ";"
		}
	} else {
		content += ";;;;;;;"
	}

	if p.ExtraData != nil {
		args := make([]string, 0)

		for key, value := range *p.ExtraData {
			switch v := value.(type) {
			case string:
				args = append(args, fmt.Sprintf("%s:%s", key, strings.TrimSpace(v)))
			case int:
			case int8:
			case int16:
			case int32:
			case int64:
			case uint:
			case uint8:
			case uint16:
			case uint32:
			case uint64:
				args = append(args, fmt.Sprintf("%s:%d", key, v))
			case float32:
			case float64:
				args = append(args, fmt.Sprintf("%s:%f", key, v))
			case bool:
				args = append(args, fmt.Sprintf("%s:%t", key, v))
			default:
				args = append(args, fmt.Sprintf("%s:%s", key, v))
			}
		}
		content += strings.Join(args, ",") + ";"
	} else {
		content += ";"
	}

	crc := calculateCrc([]byte(content))
	content = fmt.Sprintf("<Pd>%s%04X</Pd>", content, crc)
	return &content
}

type Position struct {
	// Is the latitude of the device
	Latitude *float64

	// Is the longitude of the device
	Longitude *float64

	// Is the altitude of the device
	Altitude *float64

	// Is the speed of the device
	Speed *float64

	// Is the direction of the device
	Direction *float64

	// Is the satellite count of the device
	SatelliteCount *int

	// Is the horizontal accuracy of the device
	Hdop *float64
}
