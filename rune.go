package parseg

import (
	"unicode"

	"github.com/ajiyoshi-vg/parseg/stream"
)

func Satisfy(pred func(rune) bool) Parser[rune] {
	return ParserFunc[rune](func(r stream.Stream) (*rune, int, error) {
		ret, n, err := r.ReadRune()
		if err != nil {
			return nil, n, err
		}
		if !pred(ret) {
			return nil, 0, r.UnreadRune()
		}
		return &ret, n, nil
	})
}

func AnyRune() Parser[rune] {
	return Satisfy(func(rune) bool { return true })
}

func Rune(expect rune) Parser[rune] {
	return Satisfy(func(actual rune) bool { return expect == actual })
}

func Digit() Parser[rune] {
	return Satisfy(func(actual rune) bool { return unicode.IsDigit(actual) })
}
