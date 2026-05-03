package client_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
)

func TestPr_FromPacket_ToPacket(t *testing.T) {
	packet := client.PrPacket{}
	encoded := *packet.ToPacket()

	raw := encoded
	decoded := client.PrPacket{}
	if err := decoded.FromPacket(&raw); err != nil {
		t.Fatalf("FromPacket failed: %v", err)
	}
	if *decoded.ToPacket() != encoded {
		t.Errorf("round-trip mismatch")
	}
}

func TestPr_FromPacket_Errors(t *testing.T) {
	prPkt := client.PrPacket{}
	validPacket := *prPkt.ToPacket()

	tests := []struct {
		name  string
		input string
	}{
		{"wrong prefix/suffix", "<Xx>" + validPacket[4:len(validPacket)-5] + "</Xx>"},
		{"not 2 parts", "<Pr>only_one_part</Pr>"},
		{"bad CRC hex", "<Pr>;XXXX</Pr>"},
		{"CRC mismatch", "<Pr>;0000</Pr>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := tt.input
			p := client.PrPacket{}
			if err := p.FromPacket(&raw); err == nil {
				t.Errorf("expected error for %q", tt.name)
			}
		})
	}
}
