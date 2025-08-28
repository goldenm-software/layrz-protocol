package layrzprotocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Based on the `Layrz Protocol v2` specification. AcPacket is the Command queue packet
// sent from the server to the device
type AcPacket struct {
	// Is the list of commands that the server wants to send to the device
	Commands []CommandDefinition
}

// FromPacket is a method that converts a raw packet to a PaPacket
// based on the `Layrz Protocol v2` specification
//
// Returns an error if the packet is invalid
func (p *AcPacket) FromPacket(raw *string) error {
	if !strings.HasPrefix(*raw, "<Ac>") || !strings.HasSuffix(*raw, "</Ac>") {
		return errors.New("invalid package, should be <Ac>...</Ac>")
	}

	*raw = strings.TrimPrefix(*raw, "<Ac>")
	*raw = strings.TrimSuffix(*raw, "</Ac>")

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

	if len(parts) == 0 {
		return errors.New("invalid package, should contain at least one command")
	}

	if len(parts)%4 != 0 {
		return errors.New("invalid package, should contain a multiple of 4 parts")
	}

	commands := make([]CommandDefinition, 0)

	for i := 0; i < len(parts); i += 4 {
		rawCommandId := parts[i]
		rawCommandName := parts[i+1]
		rawArgs := parts[i+2]
		rawCrc := parts[i+3]

		receivedCrc, err := strconv.ParseUint(rawCrc, 16, 16)
		if err != nil {
			return errors.New("cannot convert CRC to integer")
		}

		calculatedCrc := calculateCrc([]byte(fmt.Sprintf("%s;%s;%s;", rawCommandId, rawCommandName, rawArgs)))

		if calculatedCrc != uint16(receivedCrc) {
			return fmt.Errorf("invalid CRC, received: %04X, calculated: %04X", receivedCrc, calculatedCrc)
		}

		commandId, err := strconv.Atoi(rawCommandId)
		if err != nil {
			return errors.New("cannot convert command id to integer")
		}

		args := parseArgs(rawArgs)

		commands = append(commands, CommandDefinition{
			CommandId:   commandId,
			CommandName: &rawCommandName,
			Args:        args,
		})
	}

	p.Commands = commands
	return nil
}

// ToPacket is a method that converts a PaPacket to a raw packet
// based on the `Layrz Protocol v2` specification
func (p *AcPacket) ToPacket() *string {
	content := ""

	commands := make([]string, 0)

	for _, command := range p.Commands {
		args := make([]string, 0)

		for key, value := range *command.Args {
			args = append(args, fmt.Sprintf("%s:%v", key, value))
		}

		cmd := fmt.Sprintf(
			"%d;%s;%s;",
			command.CommandId,
			*command.CommandName,
			strings.Join(args, ","),
		)

		crc := calculateCrc([]byte(cmd))
		cmd += fmt.Sprintf("%04X", crc)

		commands = append(commands, cmd)
	}

	content += strings.Join(commands, ";")
	content += ";"

	crc := calculateCrc([]byte(content))
	content += fmt.Sprintf("%04X", crc)

	content = fmt.Sprintf("<Ac>%s</Ac>", content)

	return &content
}

// Defines the command structure based on the `Layrz Protocol v2` specification
type CommandDefinition struct {
	// Is the command id, this value is unique and should be used to
	// send the ACK packet PdPacket to the server
	CommandId int

	// Is the command name, this value is used to identify the command
	CommandName *string

	// Is the command arguments, may contain any value depending of
	// the command definition
	Args *map[string]interface{}
}
