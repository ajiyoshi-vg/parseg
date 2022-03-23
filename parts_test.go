package parseg

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumber(t *testing.T) {
	cases := map[string]struct {
		input  string
		expect *string
		rest   string
	}{
		"Number1": {
			input:  "123abc",
			expect: Ptr("123"),
			rest:   "abc",
		},
		"Number2": {
			input:  "1",
			expect: Ptr("1"),
			rest:   "",
		},
		"Number fail": {
			input:  "abc",
			expect: nil,
			rest:   "abc",
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			r := bytes.NewReader([]byte(c.input))
			actual, _, err := Number().Parse(r)
			assert.NoError(t, err)
			assert.Equal(t, c.expect, actual)

			rest, err := io.ReadAll(r)
			assert.NoError(t, err)
			assert.Equal(t, c.rest, string(rest))
		})
	}
}
