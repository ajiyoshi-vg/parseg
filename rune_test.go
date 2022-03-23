package parseg

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRune(t *testing.T) {
	cases := map[string]struct {
		input  string
		parser Parser[rune]
		expect *rune
		rest   string
	}{
		"AnyRune": {
			input:  "abc",
			parser: AnyRune(),
			expect: Ptr('a'),
			rest:   "bc",
		},
		"Rune ok": {
			input:  "abc",
			parser: Rune('a'),
			expect: Ptr('a'),
			rest:   "bc",
		},
		"Rune ng": {
			input:  "abc",
			parser: Rune('z'),
			expect: nil,
			rest:   "abc",
		},
		"Digit ok": {
			input:  "1bc",
			parser: Digit(),
			expect: Ptr('1'),
			rest:   "bc",
		},
		"Digit ng": {
			input:  "abc",
			parser: Digit(),
			expect: nil,
			rest:   "abc",
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			r := bytes.NewReader([]byte(c.input))
			actual, _, err := c.parser.Parse(r)
			assert.NoError(t, err)
			assert.Equal(t, c.expect, actual)

			rest, err := io.ReadAll(r)
			assert.NoError(t, err)
			assert.Equal(t, c.rest, string(rest))
		})
	}
}
