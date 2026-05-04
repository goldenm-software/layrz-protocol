package server

import (
	"fmt"
	"strings"
)

func Decode(dataBytes []byte) (ServerPackets, error) {
	data := string(dataBytes)
	if strings.HasPrefix(data, "<Ab>") && strings.HasSuffix(data, "</Ab>") {
		packet := &AbPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	if strings.HasPrefix(data, "<Ac>") && strings.HasSuffix(data, "</Ac>") {
		packet := &AcPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	if strings.HasPrefix(data, "<Ao>") && strings.HasSuffix(data, "</Ao>") {
		packet := &AoPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	if strings.HasPrefix(data, "<Ar>") && strings.HasSuffix(data, "</Ar>") {
		packet := &ArPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	if strings.HasPrefix(data, "<As>") && strings.HasSuffix(data, "</As>") {
		packet := &AsPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	if strings.HasPrefix(data, "<Au>") && strings.HasSuffix(data, "</Au>") {
		packet := &AuPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	return nil, fmt.Errorf("invalid packet: %s", data)
}
