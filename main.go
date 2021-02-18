package main

import "fmt"

func main() {
	sexp := &cell{
		symbol("cons"),
		&cell{
			symbol("a"),
			&cell{symbol("b"), symNil},
		},
	}
	sexp2 := sexpList(
		symbol("cons"),
		symbol("a"),
		symbol("b"),
	)

	fmt.Println(sexp)
	fmt.Println(sexp.Eval())

	fmt.Println(sexp2)
	fmt.Println(sexp2.Eval())
	fmt.Println(sexpList().Eval())
}
