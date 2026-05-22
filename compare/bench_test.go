// Package compare benchmarks parseg against goparsec on equivalent tasks.
//
// parseg:   byte-stream (io.ReadSeeker), Parser[T] returns (*T, int, error)
// goparsec: string cursor (Input),       Parser[T] returns (T, Input, error)
package compare

import (
	"bytes"
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
	"github.com/ajiyoshi-vg/parseg"
	parsegexpr "github.com/ajiyoshi-vg/parseg/expr"
)

// ---- Natural ----------------------------------------------------------------

func BenchmarkNatural_parseg(b *testing.B) {
	p := parseg.Natural()
	input := []byte("1234567890")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkNatural_goparsec(b *testing.B) {
	p := parsec.Natural()
	input := "1234567890"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parsec.Run(p, input)
	}
}

// ---- String matching --------------------------------------------------------

func BenchmarkString_parseg(b *testing.B) {
	p := parseg.String("hello")
	input := []byte("hello world")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkString_goparsec(b *testing.B) {
	p := parsec.String("hello")
	input := "hello world"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parsec.Run(p, input)
	}
}

// ---- Choice / OneOf (match) -------------------------------------------------

func BenchmarkChoice_parseg(b *testing.B) {
	p := parseg.OneOf(parseg.Rune('+'), parseg.Rune('-'), parseg.Rune('*'), parseg.Rune('/'))
	input := []byte("*123")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkChoice_goparsec(b *testing.B) {
	p := parsec.Choice(parsec.Char('+'), parsec.Char('-'), parsec.Char('*'), parsec.Char('/'))
	input := "*123"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parsec.Run(p, input)
	}
}

// ---- Choice / OneOf (no match) ----------------------------------------------

func BenchmarkChoiceNoMatch_parseg(b *testing.B) {
	p := parseg.OneOf(parseg.Rune('+'), parseg.Rune('-'), parseg.Rune('*'), parseg.Rune('/'))
	input := []byte("(123)")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkChoiceNoMatch_goparsec(b *testing.B) {
	p := parsec.Choice(parsec.Char('+'), parsec.Char('-'), parsec.Char('*'), parsec.Char('/'))
	input := "(123)"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parsec.Run(p, input)
	}
}

// ---- Many -------------------------------------------------------------------

func BenchmarkMany_parseg(b *testing.B) {
	p := parseg.Many(parseg.Digit())
	input := []byte("123456")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkMany_goparsec(b *testing.B) {
	p := parsec.Many(parsec.Digit())
	input := "123456"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parsec.Run(p, input)
	}
}

// ---- Sequence ---------------------------------------------------------------

func BenchmarkSequence_parseg(b *testing.B) {
	p := parseg.SequenceOf([]parseg.Parser[rune]{parseg.Rune('a'), parseg.Rune('b'), parseg.Rune('c')})
	input := []byte("abcdef")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkSequence_goparsec(b *testing.B) {
	p := parsec.Count(3, parsec.Choice(parsec.Char('a'), parsec.Char('b'), parsec.Char('c')))
	input := "abcdef"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parsec.Run(p, input)
	}
}

// ---- Arithmetic expression: simple (1+2) ------------------------------------

func BenchmarkExprSimple_parseg(b *testing.B) {
	p := parsegexpr.Parser()
	input := []byte("1+2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkExprSimple_goparsec(b *testing.B) {
	p := buildGoparsecExpr()
	input := "1+2"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parsec.Run(p, input)
	}
}

// ---- Arithmetic expression: complex (1+2*6/(10-7)) --------------------------

func BenchmarkExprComplex_parseg(b *testing.B) {
	p := parsegexpr.Parser()
	input := []byte("1+2*6/(10-7)")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkExprComplex_goparsec(b *testing.B) {
	p := buildGoparsecExpr()
	input := "1+2*6/(10-7)"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parsec.Run(p, input)
	}
}

// ---- Arithmetic expression: nested parens ((1+2)*(3+4))/(5-2) ---------------

func BenchmarkExprNested_parseg(b *testing.B) {
	p := parsegexpr.Parser()
	input := []byte("((1+2)*(3+4))/(5-2)")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkExprNested_goparsec(b *testing.B) {
	p := buildGoparsecExpr()
	input := "((1+2)*(3+4))/(5-2)"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parsec.Run(p, input)
	}
}

// buildGoparsecExpr constructs an arithmetic expression parser using goparsec.
// Grammar: expr = chainl1(term, addOp), term = chainl1(factor, mulOp)
//
//	factor = natural | '(' expr ')'
func buildGoparsecExpr() parsec.Parser[int] {
	var expr parsec.Parser[int]

	opParser := func(c rune, fn func(int, int) int) parsec.Parser[func(int, int) int] {
		return parsec.Map(parsec.Char(c), func(rune) func(int, int) int { return fn })
	}

	factor := func(in parsec.Input) (int, parsec.Input, error) {
		paren := parsec.Between(
			parsec.Char('('),
			parsec.Char(')'),
			parsec.Parser[int](func(in parsec.Input) (int, parsec.Input, error) {
				return expr(in)
			}),
		)
		return parsec.Choice(parsec.Natural(), paren)(in)
	}

	mulOp := parsec.Choice(
		opParser('*', func(a, b int) int { return a * b }),
		opParser('/', func(a, b int) int { return a / b }),
	)
	addOp := parsec.Choice(
		opParser('+', func(a, b int) int { return a + b }),
		opParser('-', func(a, b int) int { return a - b }),
	)

	term := parsec.Chainl1(parsec.Parser[int](factor), mulOp)
	expr = parsec.Chainl1(term, addOp)

	return expr
}
