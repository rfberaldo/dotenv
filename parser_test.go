package dotenv

import (
	"reflect"
	"testing"
)

type kv = map[string]string

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect kv
	}{
		{
			name:   "simple key-value should work",
			input:  "KEY=VALUE",
			expect: kv{"KEY": "VALUE"},
		},
		{
			name:   "LF multiple key-values should work",
			input:  "KEY=VALUE\nFOO=BAR\nBAR=BAZ\n",
			expect: kv{"KEY": "VALUE", "FOO": "BAR", "BAR": "BAZ"},
		},
		{
			name:   "CRLF multiple key-values should work",
			input:  "KEY=VALUE\r\nFOO=BAR\r\nBAR=BAZ\r\n",
			expect: kv{"KEY": "VALUE", "FOO": "BAR", "BAR": "BAZ"},
		},
		{
			name:   "values with spaces should work",
			input:  "KEY=VALUE WITH SPACES",
			expect: kv{"KEY": "VALUE WITH SPACES"},
		},
		{
			name:   "comment should skip",
			input:  "KEY=VALUE\nFOO=BAR\n#BAR=BAZ",
			expect: kv{"KEY": "VALUE", "FOO": "BAR"},
		},
		{
			name:   "LF empty lines should skip",
			input:  "KEY=VALUE\nFOO=BAR\n\nBAR=BAZ\n\n",
			expect: kv{"KEY": "VALUE", "FOO": "BAR", "BAR": "BAZ"},
		},
		{
			name:   "CRLF empty lines should skip",
			input:  "KEY=VALUE\r\nFOO=BAR\r\n\r\nBAR=BAZ\r\n\r\n",
			expect: kv{"KEY": "VALUE", "FOO": "BAR", "BAR": "BAZ"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Parse(tt.input); err != nil {
				t.Errorf("error not expected: %s", err)
			} else if !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("\nexpected: %s\n     got: %s", tt.expect, got)
			}
		})
	}
}
