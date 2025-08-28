package layrzprotocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Based on the `Layrz Protocol v2` specification. AbPacket is the Blueooth Low Energy packet
// sent from the server to the device
type AbPacket struct {
	// Is the list of devices that the server wants to send from the device
	Devices *[]BleData
}

// FromPacket is a method that converts a raw packet to a PaPacket
// based on the `Layrz Protocol v2` specification
//
// Returns an error if the packet is invalid
func (p *AbPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Ab>") || !strings.HasSuffix(*raw, "</Ab>") {
		return errors.New("invalid package, should be <Ab>...</Ab>")
	}

	*raw = strings.TrimPrefix(*raw, "<Ab>")
	*raw = strings.TrimSuffix(*raw, "</Ab>")

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

	devices := make([]BleData, 0)

	for _, part := range parts {
		subparts := strings.Split(part, ":")
		if len(subparts) != 2 {
			return errors.New("invalid device definition")
		}

		macAddress := ""

		for j := 0; j < len(subparts[0]); j += 2 {
			macAddress += subparts[0][j : j+2]
			if j != len(subparts[0])-2 {
				macAddress += ":"
			}

		}

		devices = append(devices, BleData{
			MacAddress: &macAddress,
			Model:      &subparts[1],
		})

		p.Devices = &devices
	}

	return nil
}

// ToPacket is a method that converts a PaPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *AbPacket) ToPacket() *string {
	devices := make([]string, 0)

	for _, device := range *p.Devices {
		macAddress := strings.ReplaceAll(*device.MacAddress, ":", "")
		devices = append(devices, macAddress+":"+(*device.Model))
	}

	content := strings.Join(devices, ";")
	content += ";"

	crc := calculateCrc([]byte(content))
	content += fmt.Sprintf("%04X", crc)

	content = fmt.Sprintf("<Ab>%s</Ab>", content)
	return &content
}

type BleData struct {
	// Is the MAC address of the device
	MacAddress *string

	// Is the Model identifier of the device
	Model *string
}
