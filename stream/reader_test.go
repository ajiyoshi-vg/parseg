package stream

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestDupe(t *testing.T) {
	cases := map[string]struct {
		input Stream
		check func(Stream) (int, error)
	}{
		"normal": {
			input: bytes.NewReader([]byte("日本語")),
			check: checkUnread,
		},
	}
	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			n, err := c.check(c.input)
			assert.NoError(t, err)

			_, err = c.input.Seek(int64(-n), io.SeekCurrent)
			assert.NoError(t, err)

			y, err := c.check(c.input)
			assert.Equal(t, n, y)
			assert.NoError(t, err)
		})
	}
}

func checkUnread(r Stream) (n int, _ error) {
	x, num, err := r.ReadRune()
	n += num
	if err != nil {
		return n, err
	}
	if x != '日' {
		return n, fmt.Errorf("want %c got %c", '日', x)
	}

	y, num, err := r.ReadRune()
	n += num
	if err != nil {
		return n, err
	}
	if y != '本' {
		return n, fmt.Errorf("want %c got %c", '本', y)
	}

	if err := r.UnreadRune(); err != nil {
		return n, err
	}
	n -= utf8.RuneLen(y)

	z, num, err := r.ReadRune()
	n += num
	if err != nil {
		return n, err
	}
	if z != '本' {
		return n, fmt.Errorf("want %c got %c", '本', z)
	}

	return n, nil
}
