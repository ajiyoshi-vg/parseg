package expr

import (
	"bytes"
	"io"
	"testing"

	"github.com/ajiyoshi-vg/parseg"
	"github.com/stretchr/testify/assert"
)

func TestExpr(t *testing.T) {
	cases := map[string]struct {
		input  string
		parser parseg.Parser[Expr]
		expect int
		rest   string
	}{
		"1+2": {
			input:  "1+2",
			parser: Parser(),
			expect: 3,
			rest:   "",
		},
		"1+2*6/(10-7)": {
			input:  "1+2*6/(10-7)",
			parser: Parser(),
			expect: 5,
			rest:   "",
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			r := bytes.NewReader([]byte(c.input))
			actual, _, err := c.parser.Parse(r)
			assert.NoError(t, err)
			assert.Equal(t, c.expect, (*actual).eval())

			rest, err := io.ReadAll(r)
			assert.NoError(t, err)
			assert.Equal(t, c.rest, string(rest))
		})
	}

}
