package main

import "fmt"

type appliable interface {
	Apply(args []sexp) (sexp, error)
}

type function func(args ...sexp) (sexp, error)

func (fn function) Apply(args []sexp) (sexp, error) {
	args2 := []sexp{}
	for i, arg := range args {
		arg2, err := arg.Eval()
		if err != nil {
			return nil, fmt.Errorf("args[%d]: %w", i, err)
		}

		args2 = append(args2, arg2)
	}

	return fn(args2...)
}

type specialForm func(args ...sexp) (sexp, error)

func (fn specialForm) Apply(args []sexp) (sexp, error) {
	return fn(args...)
}
