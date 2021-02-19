package main

import "fmt"

func main() {
	l := func(args ...sexp) sexp {
		return must(list(args...))
	}

	srcs := []sexp{
		&cell{
			symbol("cons"),
			&cell{
				symbol("a"),
				&cell{symbol("b"), symNil},
			},
		},
		l(
			symbol("cons"),
			symbol("a"),
			l(symbol("cons"), symbol("b"), symbol("c")),
		),
		l(),
		l(
			symbol("cons"),
			l(symbol("eq"), symbol("a"), symbol("a")),
			l(symbol("eq"), symbol("b"), symbol("c")),
		),
	}

	for _, src := range srcs {
		fmt.Println(src)
		fmt.Println(src.Eval())
		fmt.Println()
	}
}
