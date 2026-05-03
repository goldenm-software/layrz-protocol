package wire

import (
	"fmt"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected uint16
	}{
		{
			name:     "empty",
			input:    []byte{},
			expected: 0x0000,
		},
		{
			name:     "single zero byte",
			input:    []byte{0x00},
			expected: Calculate([]byte{0x00}),
		},
		{
			name:     "single byte 0xFF",
			input:    []byte{0xFF},
			expected: Calculate([]byte{0xFF}),
		},
		{
			name:     "semicolon separator",
			input:    []byte(";"),
			expected: Calculate([]byte(";")),
		},
		{
			name:     "multi byte ascii",
			input:    []byte("hello"),
			expected: Calculate([]byte("hello")),
		},
		{
			name:     "known protocol content",
			input:    []byte("1700000000;19.430000;-99.180000;2240.000000;;;;;;;"),
			expected: Calculate([]byte("1700000000;19.430000;-99.180000;2240.000000;;;;;;;")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Calculate(tt.input)
			if got != tt.expected {
				t.Errorf("Calculate(%q) = %04X, want %04X", tt.input, got, tt.expected)
			}
		})
	}

	t.Run("deterministic", func(t *testing.T) {
		data := []byte("test data for crc")
		first := Calculate(data)
		second := Calculate(data)
		if first != second {
			t.Errorf("Calculate is not deterministic: %04X != %04X", first, second)
		}
	})

	t.Run("empty is 0x0000", func(t *testing.T) {
		got := Calculate([]byte{})
		if got != 0x0000 {
			t.Errorf("expected 0x0000 for empty input, got %04X", got)
		}
	})

	t.Run("hex string representation", func(t *testing.T) {
		data := []byte("abc")
		got := Calculate(data)
		hex := fmt.Sprintf("%04X", got)
		if len(hex) != 4 {
			t.Errorf("CRC hex should be 4 chars, got %q", hex)
		}
	})
}
