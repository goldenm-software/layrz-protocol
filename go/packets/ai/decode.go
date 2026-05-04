package ai

import (
	"fmt"
	"strings"
)

func Decode(dataBytes []byte) (AiPackets, error) {
	data := string(dataBytes)
	if strings.HasPrefix(data, "<Im>") && strings.HasSuffix(data, "</Im>") {
		packet := &ImPacket{}
		if err := packet.FromPacket(&data); err != nil {
			return nil, err
		}
		return packet, nil
	}

	return nil, fmt.Errorf("invalid packet: %s", data)
}
