package layrzprotocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Based on the `Layrz Protocol v2` specification. PbPacket is the Bluetooth low energy packet
// sent from the device to the server.
type PbPacket struct {
	// Is the list of BLE advertisements detected by the device.
	Advertisements *[]BleAdvertisement
}

// FromPacket is a method that converts a raw packet to a PaPacket
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

	calculatedCrc := calculateCrc([]byte(*raw))

	if calculatedCrc != uint16(receivedCrc) {
		return fmt.Errorf("invalid CRC, received: %04X, calculated: %04X", receivedCrc, calculatedCrc)
	}

	*raw = strings.TrimSuffix(*raw, ";")

	parts := strings.Split(*raw, ";")
	if len(parts)%12 != 0 {
		return errors.New("invalid advertisement definition")
	}

	var advertisements []BleAdvertisement
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
		calculatedCrc := calculateCrc([]byte(content))

		if calculatedCrc != uint16(receivedCrc) {
			return fmt.Errorf("invalid CRC, received: %04X, calculated: %04X", receivedCrc, calculatedCrc)
		}

		timestampInt, err := strconv.ParseInt(rawTimestamp, 10, 64)
		if err != nil {
			return errors.New("cannot convert timestamp to integer")
		}

		timestamp := time.Unix(timestampInt, 0)

		latitude, err := strconv.ParseFloat(rawLatitude, 64)
		if err != nil {
			return errors.New("cannot convert latitude to float")
		}

		longitude, err := strconv.ParseFloat(rawLongitude, 64)
		if err != nil {
			return errors.New("cannot convert longitude to float")
		}

		altitude, err := strconv.ParseFloat(rawAltitude, 64)
		if err != nil {
			return errors.New("cannot convert altitude to float")
		}

		rssi, err := strconv.Atoi(rawRssi)
		if err != nil {
			return errors.New("cannot convert rssi to integer")
		}

		var txPower int = -999
		if rawTxPower != "" {
			txPower, err = strconv.Atoi(rawTxPower)
			if err != nil {
				return errors.New("cannot convert tx power to integer")
			}
		}

		manufacturerData := make([]BleManufacturerData, 0)
		for _, data := range strings.Split(rawManufacturerData, ",") {
			subparts := strings.Split(data, ":")

			companyId, err := strconv.ParseInt(subparts[0], 16, 16)
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

			manufacturerData = append(manufacturerData, BleManufacturerData{
				CompanyId: int(companyId),
				Data:      []byte(data),
			})
		}

		serviceData := make([]BleServiceData, 0)
		for _, data := range strings.Split(rawServiceData, ",") {
			subparts := strings.Split(data, ":")
			if subparts[0] == "" {
				continue
			}
			uuid, err := strconv.ParseInt(subparts[0], 16, 64)
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

			serviceData = append(serviceData, BleServiceData{
				Uuid: int(uuid),
				Data: []byte(data),
			})
		}

		advertisements = append(advertisements, BleAdvertisement{
			MacAddress:       rawMacAddress,
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

// ToPacket is a method that converts a PaPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *PbPacket) ToPacket() *string {
	packets := make([]string, 0)
	for _, advertisement := range *p.Advertisements {
		packets = append(packets, p.composeAdvertisement(advertisement))
	}

	content := strings.Join(packets, ";")
	content += ";"

	crc := calculateCrc([]byte(content))
	content += fmt.Sprintf("%04X", crc)

	content = fmt.Sprintf("<Pb>%s</Pb>", content)
	return &content
}

func (p *PbPacket) composeAdvertisement(advertisement BleAdvertisement) string {
	content := ""
	content += advertisement.MacAddress + ";"
	content += strconv.FormatInt(advertisement.Timestamp.Unix(), 10) + ";"
	content += fmt.Sprintf("%f", advertisement.Latitude) + ";"
	content += fmt.Sprintf("%f", advertisement.Longitude) + ";"
	content += fmt.Sprintf("%f", advertisement.Altitude) + ";"
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

	crc := calculateCrc([]byte(content))
	content += fmt.Sprintf("%04X", crc)

	return content
}

// Defines a detection of a BLE advertisement.
type BleAdvertisement struct {
	// Is the detected Mac Address. This Mac Adress came from the detected device, not the
	// device that is detecting.
	MacAddress string
	// Is when the device was detected.
	Timestamp time.Time

	// Is the closest latitude of the device. Defined by the device that is detecting.
	// This value is optional
	Latitude float64

	// Is the closest longitude of the device. Defined by the device that is detecting.
	// This value is optional
	Longitude float64

	// Is the closest altitude of the device. Defined by the device that is detecting.
	// This value is optional
	Altitude float64

	// Is the signal strength of the detected device.
	Rssi int

	// Is the transmission power of the detected device.
	// This value is optional
	TxPower int

	// Is the model of the detected device. This model should be equals to the model of the device
	// and the model defined by Layrz.
	Model string

	// Is the list of manufacturer data advertised by the device.
	ManufacturerData []BleManufacturerData

	// Is the list of service data advertised by the device.
	ServiceData []BleServiceData

	// Is the name of the device. This name is optional.
	DeviceName string
}

// Defines the manufacturer data advertised by the device.
type BleManufacturerData struct {
	// Is the manufacturer identifier.
	CompanyId int

	// Is the manufacturer data.
	Data []byte
}

// Defines the service data advertised by the device.
type BleServiceData struct {
	// Is the service UUID.
	Uuid int

	// Is the service data.
	Data []byte
}
