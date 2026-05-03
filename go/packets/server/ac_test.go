package server_test

import (
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/definitions"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
)

func TestAc_FromPacket_ToPacket(t *testing.T) {
	tests := []struct {
		name     string
		commands []definitions.CommandDefinition
	}{
		{
			name: "single command",
			commands: []definitions.CommandDefinition{
				{CommandId: 1, CommandName: stringPtr("test"), Args: map[string]any{}},
			},
		},
		{
			name: "multiple commands with args",
			commands: []definitions.CommandDefinition{
				{CommandId: 1, CommandName: stringPtr("cmd1"), Args: map[string]any{"arg1": "value1"}},
				{CommandId: 2, CommandName: stringPtr("cmd2"), Args: map[string]any{"arg2": 42}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := server.AcPacket{Commands: tt.commands}
			encoded := *packet.ToPacket()

			raw := encoded
			decoded := server.AcPacket{}
			if err := decoded.FromPacket(&raw); err != nil {
				t.Fatalf("FromPacket failed: %v", err)
			}
			if len(decoded.Commands) != len(tt.commands) {
				t.Errorf("command count mismatch: got %d, want %d", len(decoded.Commands), len(tt.commands))
			}
			reencoded := *decoded.ToPacket()
			if reencoded != encoded {
				t.Errorf("round-trip mismatch:\n  got  %s\n  want %s", reencoded, encoded)
			}
		})
	}
}
