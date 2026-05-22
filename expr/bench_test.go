package expr

import (
	"bytes"
	"testing"
)

var benchCases = []struct {
	name  string
	input []byte
}{
	{"simple_add", []byte("1+2")},
	{"simple_mul", []byte("3*4")},
	{"mixed", []byte("1+2*3")},
	{"parens", []byte("(1+2)*3")},
	{"complex", []byte("1+2*6/(10-7)")},
	{"nested_parens", []byte("((1+2)*(3+4))/(5-2)")},
}

func BenchmarkParser(b *testing.B) {
	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			p := Parser()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r := bytes.NewReader(bc.input)
				_, _, _ = p.Parse(r)
			}
		})
	}
}

func BenchmarkParserBuild(b *testing.B) {
	input := []byte("1+2*6/(10-7)")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := Parser()
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}
