package parseg

import (
	"unicode"

	"github.com/ajiyoshi-vg/parseg/stream"
)

func IntoString(p Parser[[]rune]) Parser[string] {
	return Map(p, intoString)
}

func intoString(xs []rune) string {
	return string(xs)
}

func Number() Parser[string] {
	return IntoString(Many1(Digit()))
}

func Natural() Parser[int] {
	return ParserFunc[int](func(r stream.Stream) (*int, int, error) {
		ch, n, err := r.ReadRune()
		if err != nil || !unicode.IsDigit(ch) {
			if n > 0 {
				_ = r.UnreadRune()
			}
			return nil, 0, err
		}
		val := int(ch - '0')
		total := n
		for {
			ch, n, err = r.ReadRune()
			if err != nil || !unicode.IsDigit(ch) {
				if n > 0 {
					_ = r.UnreadRune()
				}
				break
			}
			val = val*10 + int(ch-'0')
			total += n
		}
		return new(val), total, nil
	})
}

func String(expect string) Parser[string] {
	runes := []rune(expect)
	return ParserFunc[string](func(r stream.Stream) (*string, int, error) {
		total := 0
		for _, want := range runes {
			ch, n, err := r.ReadRune()
			total += n
			if err != nil {
				return nil, total, err
			}
			if ch != want {
				return nil, total, nil
			}
		}
		return new(expect), total, nil
	}).TryParser()
}
