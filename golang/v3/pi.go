package layrzprotocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type FirmwareBranch string

const (
	Development FirmwareBranch = "1"
	Stable      FirmwareBranch = "0"
)

// Based on the `Layrz Protocol v2` specification. PiPacket is the identification packet
// sent from the device to the server
type PiPacket struct {
	// [ident] is the Unique identifier, sent as part of the package as `IMEI`
	Ident string

	// [firmwareId] is the firmware internal ID, this is in newer versions is a Layrz ID, otherwise is a unique
	// identifier set by the hardware department of Layrz LTD. Also, is idenfified in the package as `FW_ID`
	FirmwareId string

	// [firmwareBuild] is the firmware version, is an incremental number that is increased in each release.
	// This is identified in the package as `FW_BUILD`
	FirmwareBuild int

	// [deviceId] is the device internal ID, this is in newer versions is a Layrz ID, otherwise is a unique
	// identifier set by the hardware department of Layrz LTD. Also, is idenfified in the package as `SYS_DEV_ID`
	DeviceId int

	// [hardwareId] is the hardware internal ID, this is in newer versions is a Layrz ID, otherwise is a unique
	// identifier set by the hardware department of Layrz LTD. Also, is idenfified in the package as `SYS_DEV_HW_ID`
	HardwareId int

	// [modelId] is the model internal ID, this is in newer versions is a Layrz ID, otherwise is a unique
	// identifier set by the hardware department of Layrz LTD. Also, is idenfified in the package
	// as `SYS_DEV_MODEL_ID`
	ModelId int

	// [firmwareBranch] is the branch of the firmware, this is identified in the package as `SYS_DEV_FW_BRANCH`
	FirmwareBranch FirmwareBranch

	// [fotaEnabled] is a boolean that indicates if the device is capable of receiving FOTA updates.
	// This is identified in the package as `FOTA_ENABLED`
	FotaEnabled bool
}

// FromPacket is a method that converts a raw packet to a PaPacket
// based on the `Layrz Protocol v2` specification
//
// Returns an error if the packet is invalid
func (p *PiPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Pi>") || !strings.HasSuffix(*raw, "</Pi>") {
		return errors.New("invalid package, should be <Pi>...</Pi>")
	}

	*raw = strings.TrimPrefix(*raw, "<Pi>")
	*raw = strings.TrimSuffix(*raw, "</Pi>")

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

	if len(parts) != 8 {
		return errors.New("invalid package, should contain 8 parts")
	}

	p.Ident = parts[0]
	p.FirmwareId = parts[1]

	firmwareBuild, err := strconv.Atoi(parts[2])
	if err != nil {
		return errors.New("cannot convert firmware build to integer")
	}
	p.FirmwareBuild = firmwareBuild

	deviceId, err := strconv.Atoi(parts[3])
	if err != nil {
		return errors.New("cannot convert device id to integer")
	}
	p.DeviceId = deviceId

	hardwareId, err := strconv.Atoi(parts[4])
	if err != nil {
		return errors.New("cannot convert hardware id to integer")
	}
	p.HardwareId = hardwareId

	modelId, err := strconv.Atoi(parts[5])
	if err != nil {
		return errors.New("cannot convert model id to integer")
	}
	p.ModelId = modelId

	p.FirmwareBranch = FirmwareBranch(parts[6])
	p.FotaEnabled = parts[7] == "true" || parts[7] == "1"

	return nil
}

// ToPacket is a method that converts a PaPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *PiPacket) ToPacket() *string {
	content := ""
	content += p.Ident + ";"
	content += p.FirmwareId + ";"
	content += fmt.Sprintf("%d;", p.FirmwareBuild)
	content += fmt.Sprintf("%d;", p.DeviceId)
	content += fmt.Sprintf("%d;", p.HardwareId)
	content += fmt.Sprintf("%d;", p.ModelId)
	content += string(p.FirmwareBranch) + ";"
	if p.FotaEnabled {
		content += "true"
	} else {
		content += "false"
	}
	content += ";"

	crc := calculateCrc([]byte(content))
	content = fmt.Sprintf("<Pi>%s%04X</Pi>", content, crc)
	return &content
}
