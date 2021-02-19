package main

import "fmt"

func main() {
	l := func(args ...sexp) sexp {
		return joinSexps(args)
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
			l(symbol("cons"), symbol("b"), l(symbol("quote"), symbol("c"))),
		),
		l(
			symbol("quote"),
			symbol("a"),
		),
		l(),
		l(
			symbol("cons"),
			l(symbol("eq"), symbol("a"), symbol("a")),
			l(symbol("eq"), symbol("b"), symbol("c")),
		),
		l(
			symbol("cond"),
			l(l(symbol("eq"), symbol("a"), symbol("a")), symbol("b")),
			l(symbol("t"), symbol("z")),
		),
		l(l(
			symbol("lambda"),
			l(symbol("a"), symbol("b")),
			l(symbol("eq"), symbol("a"), symbol("b")),
		), l(symbol("quote"), symbol("x")), l(symbol("quote"), symbol("x"))),
	}

	for _, src := range srcs {
		fmt.Println(src)
		fmt.Println(src.Eval(newEnvironment()))
		fmt.Println()
	}
}
