package parseg

import (
	"unicode"

	"github.com/ajiyoshi-vg/parseg/stream"
)

// runeCache holds pre-allocated pointers for ASCII runes to avoid heap allocation in Satisfy.
var runeCache = func() [128]*rune {
	var t [128]*rune
	for i := range t {
		r := rune(i)
		t[i] = &r
	}
	return t
}()

func cachedRune(r rune) *rune {
	if r >= 0 && int(r) < len(runeCache) {
		return runeCache[r]
	}
	return new(r)
}

func Satisfy(pred func(rune) bool) Parser[rune] {
	return ParserFunc[rune](func(r stream.Stream) (*rune, int, error) {
		ret, n, err := r.ReadRune()
		if err != nil {
			return nil, n, err
		}
		if !pred(ret) {
			return nil, 0, r.UnreadRune()
		}
		return cachedRune(ret), n, nil
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
