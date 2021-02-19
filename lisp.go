package main

import (
	"fmt"
)

type sexp interface {
	Eval(env *environment) (sexp, error)
	String() string
}

func joinSexps(sexps []sexp) sexp {
	xs := sexp(symNil)
	for i := len(sexps) - 1; i >= 0; i-- {
		xs = &cell{sexps[i], xs}
	}
	return xs
}

func must(x sexp, err error) sexp {
	if err != nil {
		panic(err)
	}
	return x
}

type symbol string

func (sym symbol) Eval(env *environment) (sexp, error) {
	return sym, nil
}

func (sym symbol) String() string {
	return string(sym)
}

type cell struct {
	car sexp
	cdr sexp
}

func (c *cell) Eval(env *environment) (sexp, error) {
	sym, ok := c.car.(symbol)
	if !ok {
		return nil, fmt.Errorf("invalid function %v", c.car)
	}

	args := c.arguments(env)
	fn := env.funcs[sym.String()]
	if fn == nil {
		return nil, fmt.Errorf("unknown function %v", c.car)
	}

	return fn.Apply(env, args)
}

func (c *cell) arguments(env *environment) []sexp {
	rest := c.cdr
	args := []sexp{}
	for rest != symNil {
		args = append(args, must(car(env, rest)))
		rest = must(cdr(env, rest))
	}

	return args
}

func (c *cell) String() string {
	return fmt.Sprintf("(%v . %v)", c.car, c.cdr)
}
