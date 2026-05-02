package layrzprotocol

import (
	"errors"
	"fmt"
	"strings"
)

func handleServerOutput(resp string) (*any, error) {
	var output any
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
	} else if strings.HasPrefix(resp, "<Ts>") && strings.HasSuffix(resp, "</Ts>") {
		ts := TsPacket{}
		err := ts.FromPacket(&resp)
		if err != nil {
			return nil, err
		}
		output = &ts
	} else if strings.HasPrefix(resp, "<Te>") && strings.HasSuffix(resp, "</Te>") {
		te := TePacket{}
		err := te.FromPacket(&resp)
		if err != nil {
			return nil, err
		}
		output = &te
	} else if strings.HasPrefix(resp, "<Im>") && strings.HasSuffix(resp, "</Im>") {
		im := ImPacket{}
		err := im.FromPacket(&resp)
		if err != nil {
			return nil, err
		}
		output = &im
	}

	if output == nil {
		return nil, errors.New("invalid packet response")
	}

	return &output, nil
}

func parsePacketToString(packet any) (*string, error) {
	var data *string
	switch p := packet.(type) {
	case *PaPacket:
		data = p.ToPacket()
	case *PbPacket:
		data = p.ToPacket()
	case *PcPacket:
		data = p.ToPacket()
	case *PdPacket:
		data = p.ToPacket()
	case *PiPacket:
		data = p.ToPacket()
	case *PsPacket:
		data = p.ToPacket()
	case *PmPacket:
		data = p.ToPacket()
	case *PrPacket:
		data = p.ToPacket()
	case *TsPacket:
		data = p.ToPacket()
	case *TePacket:
		data = p.ToPacket()
	case *ImPacket:
		data = p.ToPacket()
	default:
		return nil, fmt.Errorf("invalid packet type: %T", packet)
	}

	if data == nil {
		return nil, errors.New("invalid packet type")
	}

	return data, nil
}
