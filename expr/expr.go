package expr

import (
	"fmt"

	"github.com/ajiyoshi-vg/parseg"
)

func Parser() parseg.Parser[Expr] {
	var factor parseg.ParserFunc[Expr]

	// term :: factor [ ('*'|'/') factor ]*
	term := parseg.Map(
		parseg.Cons(
			parseg.Map[Expr](&factor, makeUnit(Mul)),
			parseg.Many(
				parseg.OneOf(
					parseg.Next(
						parseg.Rune('*'),
						parseg.Map[Expr](&factor, makeUnit(Mul)),
					),
					parseg.Next(
						parseg.Rune('/'),
						parseg.Map[Expr](&factor, makeUnit(Div)),
					),
				),
			),
		),
		intoTermExpr,
	)

	// expr :: term [ ('+'|'-') term ]*
	expr := parseg.Map(
		parseg.Cons(
			parseg.Map(term, makeUnit(Add)),
			parseg.Many(
				parseg.OneOf(
					parseg.Next(
						parseg.Rune('+'),
						parseg.Map(term, makeUnit(Add)),
					),
					parseg.Next(
						parseg.Rune('-'),
						parseg.Map(term, makeUnit(Sub)),
					),
				),
			),
		),
		intoExprTree,
	)

	number := parseg.Map(parseg.Natural(), func(i int) Expr { return constant(i) })
	// factor :: number | '(' expr ')'
	factor = parseg.OneOf(
		number,
		parseg.Center(parseg.Rune('('), expr, parseg.Rune(')')),
	).Func()
	return expr
}

type Expr interface {
	eval() int
	String() string
}

type Op int

type binaryOp struct {
	op  Op
	lhs Expr
	rhs Expr
}

type constant int

type unit struct {
	op      Op
	operand Expr
}

const (
	UnknownOperator Op = iota
	Add
	Sub
	Mul
	Div
)

func newUnit(op Op, operand Expr) unit {
	return unit{
		op:      op,
		operand: operand,
	}
}

func makeUnit(op Op) func(Expr) unit {
	return func(x Expr) unit {
		return newUnit(op, x)
	}
}

func intoExprTree(xs []unit) Expr {
	return foldl(
		func(acc Expr, x unit) Expr {
			return &binaryOp{
				op:  x.op,
				lhs: acc,
				rhs: x.operand,
			}
		},
		Zero,
		xs,
	)
}

func intoTermExpr(xs []unit) Expr {
	return foldl(
		func(acc Expr, x unit) Expr {
			return &binaryOp{
				op:  x.op,
				lhs: acc,
				rhs: x.operand,
			}
		},
		One,
		xs,
	)
}

var (
	_    Expr = (*binaryOp)(nil)
	Zero Expr = constant(0)
	One  Expr = constant(1)
)

func (x binaryOp) eval() int {
	switch x.op {
	case Add:
		return x.lhs.eval() + x.rhs.eval()
	case Sub:
		return x.lhs.eval() - x.rhs.eval()
	case Mul:
		return x.lhs.eval() * x.rhs.eval()
	case Div:
		return x.lhs.eval() / x.rhs.eval()
	default:
		panic(x)
	}
}
func (x *binaryOp) String() string {
	if x == nil {
		return "nil"
	}
	return fmt.Sprintf("(%s %s %s)", x.op, x.lhs, x.rhs)
}
func (x constant) eval() int {
	return int(x)
}
func (x constant) String() string {
	return fmt.Sprintf("%d", x)
}
func (x Op) String() string {
	switch x {
	case Add:
		return "+"
	case Sub:
		return "-"
	case Mul:
		return "*"
	case Div:
		return "/"
	default:
		return "UnknownOperator"
	}
}
