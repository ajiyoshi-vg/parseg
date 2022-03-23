package parseg

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTryParser(t *testing.T) {
	cases := map[string]struct {
		input  string
		parser Parser[rune]
		expect *rune
		rest   string
		nRead  int
	}{
		"Rune": {
			input:  "abc",
			parser: Rune('a'),
			expect: Ptr('a'),
			rest:   "bc",
			nRead:  1,
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			r := bytes.NewReader([]byte(c.input))
			actual, n, err := TryParser(c.parser).Parse(r)
			assert.NoError(t, err)
			assert.Equal(t, c.expect, actual)
			assert.Equal(t, c.nRead, n)

			rest, err := io.ReadAll(r)
			assert.NoError(t, err)
			assert.Equal(t, c.rest, string(rest))
		})
	}
}

func TestSequenceOf(t *testing.T) {
	cases := map[string]struct {
		input  string
		parser []Parser[rune]
		expect *[]rune
		rest   string
	}{
		"SequenceOf": {
			input: "abc",
			parser: []Parser[rune]{
				Rune('a'),
				Rune('b'),
			},
			expect: Ptr([]rune{'a', 'b'}),
			rest:   "c",
		},
		"SequenceOf fail": {
			input: "abc",
			parser: []Parser[rune]{
				Rune('a'),
				Rune('z'),
			},
			expect: nil,
			rest:   "abc",
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			r := bytes.NewReader([]byte(c.input))
			actual, _, err := SequenceOf(c.parser).Parse(r)
			assert.NoError(t, err)
			assert.Equal(t, c.expect, actual)

			rest, err := io.ReadAll(r)
			assert.NoError(t, err)
			assert.Equal(t, c.rest, string(rest))
		})
	}
}

func TestOneOf(t *testing.T) {
	op := OneOf(
		Rune('+'),
		Rune('-'),
		Rune('*'),
		Rune('/'),
	)
	cases := map[string]struct {
		input  string
		expect *rune
		rest   string
	}{
		"OneOf OK+": {"+123", Ptr('+'), "123"},
		"OneOf OK-": {"-123", Ptr('-'), "123"},
		"OneOf OK*": {"*123", Ptr('*'), "123"},
		"OneOf OK/": {"/123", Ptr('/'), "123"},
		"OneOf NG":  {"(123)", nil, "(123)"},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			r := bytes.NewReader([]byte(c.input))
			actual, _, err := op.Parse(r)
			assert.NoError(t, err)
			assert.Equal(t, c.expect, actual)

			rest, err := io.ReadAll(r)
			assert.NoError(t, err)
			assert.Equal(t, c.rest, string(rest))
		})
	}
}

func TestNextPrev(t *testing.T) {
	cases := map[string]struct {
		input  string
		parser Parser[int]
		expect *int
		rest   string
	}{
		"next": {
			input:  "-123*456",
			parser: Next(Rune('-'), Natural()),
			expect: Ptr(123),
			rest:   "*456",
		},
		"next NG": {
			input:  "-abc*456",
			parser: Next(Rune('-'), Natural()),
			expect: nil,
			rest:   "-abc*456",
		},
		"prev": {
			input:  "123;",
			parser: Prev(Natural(), String(";")),
			expect: Ptr(123),
			rest:   "",
		},
		"prev NG": {
			input:  "123.",
			parser: Prev(Natural(), String(";")),
			expect: nil,
			rest:   "123.",
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
