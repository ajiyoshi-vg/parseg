package stream

import (
	"bufio"
	"bytes"
	"io"
)

type Stream interface {
	io.ReadSeeker
	ReadRune() (rune, int, error)
	UnreadRune() error
}

var (
	_ Stream = (*bytes.Reader)(nil)
	_ Stream = (*Reader)(nil)
)

type Reader struct {
	rs  io.ReadSeeker
	buf *bufio.Reader
}

func New(rs io.ReadSeeker) *Reader {
	return &Reader{
		rs:  rs,
		buf: bufio.NewReader(rs),
	}
}

func (r *Reader) Read(b []byte) (n int, err error) {
	return r.buf.Read(b)
}

func (r *Reader) ReadRune() (ch rune, size int, err error) {
	return r.buf.ReadRune()
}

func (r *Reader) UnreadRune() error {
	return r.buf.UnreadRune()
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	ret, err := r.rs.Seek(offset, whence)
	r.buf = bufio.NewReader(r.rs)
	return ret, err
}

func (r *Reader) WriteTo(w io.Writer) (int64, error) {
	return r.buf.WriteTo(w)
}
