package client

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/definitions"
	"github.com/goldenm-software/layrz-protocol/go/v3/internal/wire"
)

// Based on the `Layrz Protocol v2` specification. PbPacket is the Bluetooth low energy packet
// sent from the device to the server.
type PbPacket struct {
	// Is the list of BLE advertisements detected by the device.
	Advertisements *[]definitions.BleAdvertisement
}

// FromPacket is a method that converts a raw packet to a PbPacket
// based on the `Layrz Protocol v2` specification
//
// Returns an error if the packet is invalid
func (p *PbPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Pb>") || !strings.HasSuffix(*raw, "</Pb>") {
		return errors.New("invalid package, should be <Pb>...</Pb>")
	}

	*raw = strings.TrimPrefix(*raw, "<Pb>")
	*raw = strings.TrimSuffix(*raw, "</Pb>")

	// Get the last 4 characters
	rawCrc := (*raw)[len(*raw)-4:]
	*raw = (*raw)[:len(*raw)-4]

	receivedCrc, err := strconv.ParseUint(rawCrc, 16, 16)

	if err != nil {
		return errors.New("cannot convert CRC to integer")
	}

	calculatedCrc := wire.Calculate([]byte(*raw))

	if calculatedCrc != uint16(receivedCrc) {
		return fmt.Errorf("invalid CRC, received: %04X, calculated: %04X", receivedCrc, calculatedCrc)
	}

	*raw = strings.TrimSuffix(*raw, ";")

	parts := strings.Split(*raw, ";")
	if len(parts)%12 != 0 {
		return errors.New("invalid advertisement definition")
	}

	var advertisements []definitions.BleAdvertisement
	for i := 0; i < len(parts); i += 12 {
		rawMacAddress := parts[i]
		rawTimestamp := parts[i+1]
		rawLatitude := parts[i+2]
		rawLongitude := parts[i+3]
		rawAltitude := parts[i+4]
		rawModel := parts[i+5]
		rawDeviceName := parts[i+6]
		rawRssi := parts[i+7]
		rawTxPower := parts[i+8]
		rawManufacturerData := parts[i+9]
		rawServiceData := parts[i+10]
		rawCrc := parts[i+11]

		receivedCrc, err := strconv.ParseUint(rawCrc, 16, 16)
		if err != nil {
			return errors.New("cannot convert CRC to integer")
		}

		content := fmt.Sprintf("%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;", rawMacAddress, rawTimestamp, rawLatitude, rawLongitude, rawAltitude, rawModel, rawDeviceName, rawRssi, rawTxPower, rawManufacturerData, rawServiceData)
		calculatedCrc := wire.Calculate([]byte(content))

		if calculatedCrc != uint16(receivedCrc) {
			return fmt.Errorf("invalid CRC, received: %04X, calculated: %04X", receivedCrc, calculatedCrc)
		}

		timestampInt, err := strconv.ParseInt(rawTimestamp, 10, 64)
		if err != nil {
			return errors.New("cannot convert timestamp to integer")
		}

		timestamp := time.Unix(timestampInt, 0)

		var latitude *float64
		if rawLatitude != "" {
			v, err := strconv.ParseFloat(rawLatitude, 64)
			if err != nil {
				return errors.New("cannot convert latitude to float")
			}
			latitude = &v
		}

		var longitude *float64
		if rawLongitude != "" {
			v, err := strconv.ParseFloat(rawLongitude, 64)
			if err != nil {
				return errors.New("cannot convert longitude to float")
			}
			longitude = &v
		}

		var altitude *float64
		if rawAltitude != "" {
			v, err := strconv.ParseFloat(rawAltitude, 64)
			if err != nil {
				return errors.New("cannot convert altitude to float")
			}
			altitude = &v
		}

		rssi, err := strconv.Atoi(rawRssi)
		if err != nil {
			return errors.New("cannot convert rssi to integer")
		}

		txPower := -999
		if rawTxPower != "" {
			txPower, err = strconv.Atoi(rawTxPower)
			if err != nil {
				return errors.New("cannot convert tx power to integer")
			}
		}

		manufacturerData := make([]definitions.BleManufacturerData, 0)
		for _, data := range strings.Split(rawManufacturerData, ",") {
			if data == "" {
				continue
			}
			subparts := strings.Split(data, ":")

			companyId, err := strconv.ParseUint(subparts[0], 16, 16)
			if err != nil {
				return errors.New("cannot convert company id to integer")
			}

			data := make([]byte, 0)
			for i := 0; i < len(subparts[1]); i += 2 {
				value, err := strconv.ParseInt(subparts[1][i:i+2], 16, 64)
				if err != nil {
					return errors.New("cannot convert data to byte")
				}

				data = append(data, byte(value))
			}

			manufacturerData = append(manufacturerData, definitions.BleManufacturerData{
				CompanyId: int(companyId),
				Data:      []byte(data),
			})
		}

		serviceData := make([]definitions.BleServiceData, 0)
		for _, data := range strings.Split(rawServiceData, ",") {
			subparts := strings.Split(data, ":")
			if subparts[0] == "" {
				continue
			}
			uuid, err := strconv.ParseUint(subparts[0], 16, 16)
			if err != nil {
				return fmt.Errorf("cannot convert uuid \"%s\" to integer", subparts[0])
			}

			data := make([]byte, 0)
			for i := 0; i < len(subparts[1]); i += 2 {
				value, err := strconv.ParseInt(subparts[1][i:i+2], 16, 64)
				if err != nil {
					return errors.New("cannot convert data to byte")
				}

				data = append(data, byte(value))
			}

			serviceData = append(serviceData, definitions.BleServiceData{
				Uuid: int(uuid),
				Data: []byte(data),
			})
		}

		advertisements = append(advertisements, definitions.BleAdvertisement{
			MacAddress:       normalizeMac(rawMacAddress),
			Timestamp:        timestamp,
			Latitude:         latitude,
			Longitude:        longitude,
			Altitude:         altitude,
			Rssi:             rssi,
			TxPower:          txPower,
			Model:            rawModel,
			DeviceName:       rawDeviceName,
			ManufacturerData: manufacturerData,
			ServiceData:      serviceData,
		})
	}

	p.Advertisements = &advertisements
	return nil
}

// ToPacket is a method that converts a PbPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *PbPacket) ToPacket() *string {
	packets := make([]string, 0)
	for _, advertisement := range *p.Advertisements {
		packets = append(packets, p.composeAdvertisement(advertisement))
	}

	content := strings.Join(packets, ";")
	content += ";"

	crc := wire.Calculate([]byte(content))
	content += fmt.Sprintf("%04X", crc)

	content = fmt.Sprintf("<Pb>%s</Pb>", content)
	return &content
}

// normalizeMac converts a MAC address to uppercase colon-separated format (AA:BB:CC:DD:EE:FF).
func normalizeMac(mac string) string {
	mac = strings.ToUpper(strings.ReplaceAll(mac, ":", ""))
	if len(mac) != 12 {
		return mac
	}
	return mac[0:2] + ":" + mac[2:4] + ":" + mac[4:6] + ":" + mac[6:8] + ":" + mac[8:10] + ":" + mac[10:12]
}

func formatCoord(v *float64) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%f", *v)
}

func (p *PbPacket) composeAdvertisement(advertisement definitions.BleAdvertisement) string {
	content := ""
	content += advertisement.MacAddress + ";"
	content += strconv.FormatInt(advertisement.Timestamp.Unix(), 10) + ";"
	content += formatCoord(advertisement.Latitude) + ";"
	content += formatCoord(advertisement.Longitude) + ";"
	content += formatCoord(advertisement.Altitude) + ";"
	content += advertisement.Model + ";"
	content += advertisement.DeviceName + ";"
	content += strconv.Itoa(advertisement.Rssi) + ";"
	content += strconv.Itoa(advertisement.TxPower) + ";"

	manufacturer := make([]string, 0)
	for _, data := range advertisement.ManufacturerData {
		manufacturer = append(manufacturer, fmt.Sprintf("%04X", data.CompanyId)+":"+fmt.Sprintf("%X", data.Data))
	}
	content += strings.Join(manufacturer, ",") + ";"

	service := make([]string, 0)
	for _, data := range advertisement.ServiceData {
		service = append(service, fmt.Sprintf("%04X", data.Uuid)+":"+fmt.Sprintf("%X", data.Data))
	}

	content += strings.Join(service, ",") + ";"

	crc := wire.Calculate([]byte(content))
	content += fmt.Sprintf("%04X", crc)

	return content
}
