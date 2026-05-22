package parseg

import (
	"bytes"
	"testing"
)

func BenchmarkNatural(b *testing.B) {
	p := Natural()
	input := []byte("1234567890")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkString(b *testing.B) {
	p := String("hello")
	input := []byte("hello world")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkOneOf(b *testing.B) {
	p := OneOf(Rune('+'), Rune('-'), Rune('*'), Rune('/'))
	input := []byte("*123")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkOneOfNoMatch(b *testing.B) {
	p := OneOf(Rune('+'), Rune('-'), Rune('*'), Rune('/'))
	input := []byte("(123)")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkMany(b *testing.B) {
	p := Many(Digit())
	input := []byte("123456")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkSequenceOf(b *testing.B) {
	p := SequenceOf([]Parser[rune]{Rune('a'), Rune('b'), Rune('c')})
	input := []byte("abcdef")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkNext(b *testing.B) {
	p := Next(Rune('-'), Natural())
	input := []byte("-42")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}

func BenchmarkPrev(b *testing.B) {
	p := Prev(Natural(), String(";"))
	input := []byte("42;")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		_, _, _ = p.Parse(r)
	}
}
