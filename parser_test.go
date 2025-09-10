package dotenv

import (
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []pair
		error  bool
	}{
		{
			name:   "empty input",
			input:  "",
			expect: []pair{},
		},
		{
			name:   "only comments and blanks",
			input:  "# this is a comment\n\n   \n# another",
			expect: []pair{},
		},
		{
			name:  "single key value",
			input: "FOO=bar",
			expect: []pair{
				{"FOO", "bar"},
			},
		},
		{
			name:  "multiple key values",
			input: "FOO=bar\nBAZ=qux\nHELLO=WORLD",
			expect: []pair{
				{"FOO", "bar"},
				{"BAZ", "qux"},
				{"HELLO", "WORLD"},
			},
		},
		{
			name:  "multiple key values with spacing between equal sign",
			input: "FOO = bar\nBAZ = qux\nHELLO = WORLD",
			expect: []pair{
				{"FOO", "bar"},
				{"BAZ", "qux"},
				{"HELLO", "WORLD"},
			},
		},
		{
			name:  "multiple key values CRLF",
			input: "FOO=bar\r\nBAZ=qux\r\nHELLO=WORLD",
			expect: []pair{
				{"FOO", "bar"},
				{"BAZ", "qux"},
				{"HELLO", "WORLD"},
			},
		},
		{
			name:  "invalid line without equal sign",
			input: "FOO=bar\nINVALID_LINE\nBAZ=qux",
			error: true,
		},
		{
			name:  "leading comment and valid line",
			input: "# comment\nKEY=value",
			expect: []pair{
				{"KEY", "value"},
			},
		},
		{
			name:   "values with spaces should work",
			input:  "KEY=VALUE WITH SPACES",
			expect: []pair{{"KEY", "VALUE WITH SPACES"}},
		},
		{
			name:  "should not expand",
			input: "ENV=dev\nPORT=8000\nHOST=localhost:${PORT}",
			expect: []pair{
				{"ENV", "dev"},
				{"PORT", "8000"},
				{"HOST", "localhost:${PORT}"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(strings.NewReader(tt.input))
			if (err != nil) != tt.error {
				t.Errorf("unexpected error")
			}
			if !tt.error && !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("\nexpected: %s\n     got: %s", tt.expect, got)
			}
		})
	}
}
