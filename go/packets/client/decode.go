package client

import (
	"fmt"
	"strings"
)

func Decode(dataBytes []byte) (ClientPackets, error) {
	data := string(dataBytes)
	if strings.HasPrefix(data, "<Pa>") && strings.HasSuffix(data, "</Pa>") {
		packet := &PaPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	if strings.HasPrefix(data, "<Pb>") && strings.HasSuffix(data, "</Pb>") {
		packet := &PbPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	if strings.HasPrefix(data, "<Pc>") && strings.HasSuffix(data, "</Pc>") {
		packet := &PcPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	if strings.HasPrefix(data, "<Pd>") && strings.HasSuffix(data, "</Pd>") {
		packet := &PdPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	if strings.HasPrefix(data, "<Pi>") && strings.HasSuffix(data, "</Pi>") {
		packet := &PiPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	if strings.HasPrefix(data, "<Pm>") && strings.HasSuffix(data, "</Pm>") {
		packet := &PmPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	if strings.HasPrefix(data, "<Pr>") && strings.HasSuffix(data, "</Pr>") {
		packet := &PrPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	if strings.HasPrefix(data, "<Ps>") && strings.HasSuffix(data, "</Ps>") {
		packet := &PsPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	return nil, fmt.Errorf("invalid packet: %s", data)
}
