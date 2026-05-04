package helpers_test

import (
	"reflect"
	"testing"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/helpers"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "whitespace only",
			input: "   \n\r",
			want:  nil,
		},
		{
			name:  "plain text no tag",
			input: "hello world",
			want:  []string{"hello world"},
		},
		{
			name:  "plain text with trailing newline",
			input: "hello world\n",
			want:  []string{"hello world"},
		},
		{
			name:  "single Pa packet",
			input: "<Pa>some data</Pa>",
			want:  []string{"<Pa>some data</Pa>"},
		},
		{
			name:  "two concatenated packets",
			input: "<Pa>data1</Pa><Pb>data2</Pb>",
			want:  []string{"<Pa>data1</Pa>", "<Pb>data2</Pb>"},
		},
		{
			name:  "three concatenated packets",
			input: "<Pa>a</Pa><Pb>b</Pb><Pc>c</Pc>",
			want:  []string{"<Pa>a</Pa>", "<Pb>b</Pb>", "<Pc>c</Pc>"},
		},
		{
			name:  "trailing newline stripped from last packet",
			input: "<Pa>data</Pa>\n",
			want:  []string{"<Pa>data</Pa>"},
		},
		{
			name:  "trailing \\r\\n stripped",
			input: "<Pa>data</Pa>\r\n",
			want:  []string{"<Pa>data</Pa>"},
		},
		{
			name:  "non-P tag (As) falls through no-tag branch",
			input: "<As>server data</As>",
			want:  []string{"<As>server data</As>"},
		},
		{
			name:  "mixed: real packet tag among plain content",
			input: "<Pr>real</Pr>",
			want:  []string{"<Pr>real</Pr>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.Split(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
