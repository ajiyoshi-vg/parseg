package expr

func foldl[A, B any](f func(B, A) B, acc B, xs []A) B {
	for _, x := range xs {
		acc = f(acc, x)
	}
	return acc
}

func foldr[A, B any](f func(A, B) B, acc B, xs []A) B {
	for _, x := range xs {
		acc = f(x, acc)
	}
	return acc
}
