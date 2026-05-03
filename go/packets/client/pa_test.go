package client_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
)

func TestPa_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name     string
		ident    string
		password string
	}{
		{
			name:     "valid credentials",
			ident:    "IMEI123",
			password: "secret",
		},
		{
			name:     "empty password",
			ident:    "IMEI123",
			password: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := client.PaPacket{Ident: stringPtr(tt.ident), Password: stringPtr(tt.password)}
			encoded := *packet.ToPacket()

			raw := encoded
			decoded := client.PaPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}
			if *decoded.Ident != tt.ident {
				t.Errorf("ident mismatch: got %s, want %s", *decoded.Ident, tt.ident)
			}
			if *decoded.Password != tt.password {
				t.Errorf("password mismatch: got %s, want %s", *decoded.Password, tt.password)
			}
			if *decoded.ToPacket() != encoded {
				t.Errorf("round-trip mismatch")
			}
		})
	}
}
