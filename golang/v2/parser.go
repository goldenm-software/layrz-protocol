package layrzprotocol

import (
	"errors"
	"fmt"
	"strings"
)

func handleServerOutput(resp string) (*interface{}, error) {
	var output interface{}
	if strings.HasPrefix(resp, "<Ab>") && strings.HasSuffix(resp, "</Ab>") {
		ao := AbPacket{}
		err := ao.FromPacket(&resp)
		if err != nil {
			return nil, err
		}
		output = &ao
	} else if strings.HasPrefix(resp, "<Ac>") && strings.HasSuffix(resp, "</Ac>") {
		ac := AcPacket{}
		err := ac.FromPacket(&resp)
		if err != nil {
			return nil, err
		}
		output = &ac
	} else if strings.HasPrefix(resp, "<Ao>") && strings.HasSuffix(resp, "</Ao>") {
		ao := AoPacket{}
		err := ao.FromPacket(&resp)
		if err != nil {
			return nil, err
		}
		output = &ao
	} else if strings.HasPrefix(resp, "<Ar>") && strings.HasSuffix(resp, "</Ar>") {
		ar := ArPacket{}
		err := ar.FromPacket(&resp)
		if err != nil {
			return nil, err
		}
		output = &ar
	} else if strings.HasPrefix(resp, "<As>") && strings.HasSuffix(resp, "</As>") {
		as := AsPacket{}
		err := as.FromPacket(&resp)
		if err != nil {
			return nil, err
		}
		output = &as
	} else if strings.HasPrefix(resp, "<Au>") && strings.HasSuffix(resp, "</Au>") {
		au := AuPacket{}
		err := au.FromPacket(&resp)
		if err != nil {
			return nil, err
		}
		output = &au
	}

	if output == nil {
		return nil, errors.New("invalid packet response")
	}

	return &output, nil
}

func parsePacketToString(packet interface{}) (*string, error) {
	var data *string
	switch packet.(type) {
	case *PaPacket:
		data = packet.(*PaPacket).ToPacket()
	case *PbPacket:
		data = packet.(*PbPacket).ToPacket()
	case *PcPacket:
		data = packet.(*PcPacket).ToPacket()
	case *PdPacket:
		data = packet.(*PdPacket).ToPacket()
	case *PiPacket:
		data = packet.(*PiPacket).ToPacket()
	case *PsPacket:
		data = packet.(*PsPacket).ToPacket()
	case *PmPacket:
		data = packet.(*PmPacket).ToPacket()
	case *PrPacket:
		data = packet.(*PrPacket).ToPacket()
	default:
		return nil, fmt.Errorf("invalid packet type, should be PaPacket, PbPacket, PcPacket, PdPacket, PiPacket, PsPacket, PmPacket or PrPacket, not %T", packet)
	}

	if data == nil {
		return nil, errors.New("invalid packet type, should be PaPacket, PbPacket, PcPacket, PdPacket, PiPacket, PsPacket, PmPacket or PrPacket")
	}

	return data, nil
}
