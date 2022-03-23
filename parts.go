package parseg

import (
	"strconv"
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
	return Apply(Number(), strconv.Atoi)
}

func String(expect string) Parser[string] {
	ps := make([]Parser[rune], 0, len(expect))
	for _, r := range expect {
		ps = append(ps, Rune(r))
	}
	return IntoString(SequenceOf(ps))
}
