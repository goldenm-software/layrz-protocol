package trips

import (
	"fmt"
	"strings"
)

func Decode(dataBytes []byte) (TripsPackets, error) {
	data := string(dataBytes)
	if strings.HasPrefix(data, "<Te>") && strings.HasSuffix(data, "</Te>") {
		packet := &TePacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	if strings.HasPrefix(data, "<Ts>") && strings.HasSuffix(data, "</Ts>") {
		packet := &TsPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	return nil, fmt.Errorf("invalid packet: %s", data)
}
