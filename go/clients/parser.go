package clients

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/ai"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/trips"
)

// DecodeServerOutput decodes a server output string into a packet
func DecodeServerOutput(resp string) (*any, error) {
	var output any
	if strings.HasPrefix(resp, "<Ab>") && strings.HasSuffix(resp, "</Ab>") {
		p := server.AbPacket{}
		if err := p.FromPacket(&resp); err != nil {
			return nil, err
		}
		output = &p
	} else if strings.HasPrefix(resp, "<Ac>") && strings.HasSuffix(resp, "</Ac>") {
		p := server.AcPacket{}
		if err := p.FromPacket(&resp); err != nil {
			return nil, err
		}
		output = &p
	} else if strings.HasPrefix(resp, "<Ao>") && strings.HasSuffix(resp, "</Ao>") {
		p := server.AoPacket{}
		if err := p.FromPacket(&resp); err != nil {
			return nil, err
		}
		output = &p
	} else if strings.HasPrefix(resp, "<Ar>") && strings.HasSuffix(resp, "</Ar>") {
		p := server.ArPacket{}
		if err := p.FromPacket(&resp); err != nil {
			return nil, err
		}
		output = &p
	} else if strings.HasPrefix(resp, "<As>") && strings.HasSuffix(resp, "</As>") {
		p := server.AsPacket{}
		if err := p.FromPacket(&resp); err != nil {
			return nil, err
		}
		output = &p
	} else if strings.HasPrefix(resp, "<Au>") && strings.HasSuffix(resp, "</Au>") {
		p := server.AuPacket{} //nolint:staticcheck
		if err := p.FromPacket(&resp); err != nil {
			return nil, err
		}
		output = &p
	} else if strings.HasPrefix(resp, "<Ts>") && strings.HasSuffix(resp, "</Ts>") {
		p := trips.TsPacket{}
		if err := p.FromPacket(&resp); err != nil {
			return nil, err
		}
		output = &p
	} else if strings.HasPrefix(resp, "<Te>") && strings.HasSuffix(resp, "</Te>") {
		p := trips.TePacket{}
		if err := p.FromPacket(&resp); err != nil {
			return nil, err
		}
		output = &p
	} else if strings.HasPrefix(resp, "<Im>") && strings.HasSuffix(resp, "</Im>") {
		p := ai.ImPacket{}
		if err := p.FromPacket(&resp); err != nil {
			return nil, err
		}
		output = &p
	}
	if output == nil {
		return nil, errors.New("invalid packet response")
	}
	return &output, nil
}

// EncodeClientPacket encodes a client packet into a string
func EncodeClientPacket(packet any) (*string, error) {
	var data *string
	switch p := packet.(type) {
	case *client.PaPacket:
		data = p.ToPacket()
	case *client.PbPacket:
		data = p.ToPacket()
	case *client.PcPacket:
		data = p.ToPacket()
	case *client.PdPacket:
		data = p.ToPacket()
	case *client.PiPacket:
		data = p.ToPacket()
	case *client.PsPacket:
		data = p.ToPacket()
	case *client.PmPacket:
		data = p.ToPacket()
	case *client.PrPacket:
		data = p.ToPacket()
	case *trips.TsPacket:
		data = p.ToPacket()
	case *trips.TePacket:
		data = p.ToPacket()
	case *ai.ImPacket:
		data = p.ToPacket()
	default:
		return nil, fmt.Errorf("invalid packet type: %T", packet)
	}
	if data == nil {
		return nil, errors.New("invalid packet type")
	}
	return data, nil
}
