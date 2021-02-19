package main

import "fmt"

func main() {
	l := func(args ...sexp) sexp { return list(args...) }
	s := func(x string) symbol { return symbol(x) }
	n := func(x float64) number { return number(x) }
	q := func(x string) sexp { return l(s("quote"), s(x)) }

	srcs := []sexp{
		&cell{
			s("cons"),
			&cell{
				q("a"),
				&cell{q("b"), symNil},
			},
		},
		cons(s("cons"), cons(q("a"), cons(q("b"), symNil))),
		l(
			s("cons"),
			q("a"),
			l(s("cons"), q("b"), q("c")),
		),
		q("a"),
		s("a"),
		l(),
		l(
			s("cons"),
			l(s("eq"), q("a"), q("a")),
			l(s("eq"), q("b"), q("c")),
		),
		l(
			s("cond"),
			l(l(s("eq"), q("a"), q("a")), q("b")),
			l(s("t"), q("z")),
		),
		l(l(
			s("lambda"),
			l(s("a"), s("b")),
			l(s("eq"), s("a"), s("b")),
		), l(s("quote"), s("x")), l(s("quote"), s("x"))),
		l(s("list"), n(10), n(20)),
		l(s("eq"), n(10), n(10)),
		l(s("+"), n(10), n(20)),
		l(s("+"), n(10), l(s("-"), n(10), n(30))),
		l(s("+"), n(10), l(s("quote"), s("x"))),
	}

	env := newEnvironment()
	for _, src := range srcs {
		fmt.Println(src)
		fmt.Println(src.Eval(env))
		fmt.Println()
	}
}
