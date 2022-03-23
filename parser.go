package parseg

import (
	"io"

	"github.com/ajiyoshi-vg/parseg/stream"
)

type Parser[T any] interface {
	Parse(stream.Stream) (*T, int, error)
	TryParser() Parser[T]
	Or(g Parser[T]) Parser[T]
}

var (
	_ Parser[int] = (ParserFunc[int])(nil)
	_ Parser[int] = (*ParserFunc[int])(nil)
)

type ParserFunc[T any] func(stream.Stream) (*T, int, error)

func (f ParserFunc[T]) Parse(r stream.Stream) (*T, int, error) {
	return f(r)
}

func (f ParserFunc[T]) TryParser() Parser[T] {
	return TryParser[T](f)
}

func (f ParserFunc[T]) Or(g Parser[T]) Parser[T] {
	return Or[T](f, g)
}

func Lazy[T any](t *T) T {
	return *t
}

func Map[T, S any](p Parser[T], f func(T) S) Parser[S] {
	return ParserFunc[S](func(r stream.Stream) (*S, int, error) {
		x, n, err := p.Parse(r)
		if isError(err) {
			return nil, n, err
		}
		if x == nil {
			return nil, n, nil
		}
		return Ptr(f(*x)), n, nil
	})
}

func Apply[T, S any](p Parser[T], f func(T) (S, error)) Parser[S] {
	return ParserFunc[S](func(r stream.Stream) (*S, int, error) {
		x, n, err := p.Parse(r)
		if isError(err) {
			return nil, n, err
		}
		if x == nil {
			return nil, n, nil
		}
		ret, err := f(*x)
		if isError(err) {
			return nil, n, err
		}
		return Ptr(ret), n, nil
	})
}

func SequenceOf[T any](ps []Parser[T]) Parser[[]T] {
	return ParserFunc[[]T](func(r stream.Stream) (*[]T, int, error) {
		var parsed []T
		nRead := 0
		for _, p := range ps {
			x, n, err := p.Parse(r)
			nRead += n
			if isError(err) {
				return nil, nRead, err
			}
			if x == nil {
				return nil, nRead, nil
			}
			parsed = append(parsed, *x)
		}
		return &parsed, nRead, nil
	}).TryParser()
}

func Sequence[T any](ps ...Parser[T]) Parser[[]T] {
	return SequenceOf(ps)
}

func TryParser[T any](p Parser[T]) Parser[T] {
	return ParserFunc[T](func(r stream.Stream) (*T, int, error) {
		parsed, n, err := p.Parse(r)
		if isError(err) {
			return nil, n, err
		}
		if parsed == nil {
			_, err := r.Seek(int64(-n), io.SeekCurrent)
			return nil, 0, err
		}
		return parsed, n, nil
	})
}

func Or[T any](a, b Parser[T]) ParserFunc[T] {
	return func(r stream.Stream) (*T, int, error) {
		x, n, err := a.Parse(r)
		if isError(err) {
			return nil, n, err
		}
		if x != nil {
			return x, n, nil
		}
		y, m, err := b.Parse(r)
		m += n
		if isError(err) {
			return nil, n, err
		}
		return y, n, err
	}
}

func OneOf[T any](ps ...Parser[T]) Parser[T] {
	return ParserFunc[T](func(r stream.Stream) (*T, int, error) {
		for _, p := range ps {
			x, n, err := p.Parse(r)
			if isError(err) {
				return nil, n, err
			}
			if x != nil {
				return x, n, nil
			}
		}
		return nil, 0, nil
	})
}

func Cons[T any](car Parser[T], cdr Parser[[]T]) Parser[[]T] {
	return ParserFunc[[]T](func(r stream.Stream) (*[]T, int, error) {
		x, n, err := car.Parse(r)
		if isError(err) || x == nil {
			return nil, n, err
		}
		ret := []T{*x}
		y, m, err := cdr.Parse(r)
		n += m
		if isError(err) || y == nil {
			return Ptr(ret), n, err
		}
		return Ptr(append(ret, *y...)), n, nil
	})
}

func Many[T any](p Parser[T]) Parser[[]T] {
	return ParserFunc[[]T](func(r stream.Stream) (*[]T, int, error) {
		var ret []T
		n := 0
		for {
			x, num, err := p.Parse(r)
			n += num
			if isError(err) {
				return nil, n, err
			}
			if x == nil {
				return &ret, n, nil
			}
			ret = append(ret, *x)
		}
	})
}

func Many1[T any](p Parser[T]) Parser[[]T] {
	return Cons(p, Many(p))
}

func Next[T, S any](a Parser[T], b Parser[S]) Parser[S] {
	return ParserFunc[S](func(r stream.Stream) (*S, int, error) {
		x, n, err := a.Parse(r)
		if isError(err) {
			return nil, n, err
		}
		if x == nil {
			return nil, n, nil
		}
		ret, m, err := b.Parse(r)
		n += m
		if isError(err) {
			return nil, n, err
		}
		return ret, n, nil
	}).TryParser()
}

func Prev[T, S any](a Parser[T], b Parser[S]) Parser[T] {
	return ParserFunc[T](func(r stream.Stream) (*T, int, error) {
		ret, n, err := a.Parse(r)
		if isError(err) {
			return nil, n, err
		}
		if ret == nil {
			return nil, n, nil
		}
		x, m, err := b.Parse(r)
		n += m
		if isError(err) {
			return nil, n, err
		}
		if x == nil {
			return nil, n, nil
		}
		return ret, n, nil
	}).TryParser()
}

func Center[A, B, C any](a Parser[A], b Parser[B], c Parser[C]) Parser[B] {
	return Prev(Next(a, b), c).TryParser()
}

func Ptr[T any](x T) *T {
	return &x
}

func isError(err error) bool {
	return err != nil && err != io.EOF
}
