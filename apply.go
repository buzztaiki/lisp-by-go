package main

import "fmt"

type appliable interface {
	Apply(env *environment, args []sexp) (sexp, error)
}

type function func(env *environment, args ...sexp) (sexp, error)

func (fn function) Apply(env *environment, args []sexp) (sexp, error) {
	args2 := []sexp{}
	for i, arg := range args {
		arg2, err := arg.Eval(env)
		if err != nil {
			return nil, fmt.Errorf("args[%d]: %w", i, err)
		}

		args2 = append(args2, arg2)
	}

	return fn(env, args2...)
}

type specialForm func(env *environment, args ...sexp) (sexp, error)

func (fn specialForm) Apply(env *environment, args []sexp) (sexp, error) {
	return fn(env, args...)
}
