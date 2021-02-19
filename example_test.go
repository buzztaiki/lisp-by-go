package main

import "fmt"

func ExampleEval() {
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
		l(s("defun"), s("f"), l(s("x")), l(s("+"), s("x"), s("x"))),
		l(s("f"), n(3)),
		l(s("f"), n(3), n(4)),
	}

	env := newEnvironment()
	for _, src := range srcs {
		fmt.Println(src)
		if res, err := src.Eval(env); err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("==>", res)
		}
		fmt.Println()
	}

	// Output:
	// (cons (quote a) (quote b))
	// ==> (a . b)
	//
	// (cons (quote a) (quote b))
	// ==> (a . b)
	//
	// (cons (quote a) (cons (quote b) (quote c)))
	// ==> (a b . c)
	//
	// (quote a)
	// ==> a
	//
	// a
	// error: variable a not found
	//
	// nil
	// ==> nil
	//
	// (cons (eq (quote a) (quote a)) (eq (quote b) (quote c)))
	// ==> (t)
	//
	// (cond ((eq (quote a) (quote a)) (quote b)) (t (quote z)))
	// ==> b
	//
	// ((lambda (a b) (eq a b)) (quote x) (quote x))
	// ==> t
	//
	// (list 10 20)
	// ==> (10 20)
	//
	// (eq 10 10)
	// ==> t
	//
	// (+ 10 20)
	// ==> 30
	//
	// (+ 10 (- 10 30))
	// ==> -30
	//
	// (+ 10 (quote x))
	// error: args[1]: wrong number type argument x
	//
	// (defun f (x) (+ x x))
	// ==> f
	//
	// (f 3)
	// ==> 6
	//
	// (f 3 4)
	// error: wrong number of argument 2
}
