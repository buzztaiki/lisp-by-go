package main

import (
	"fmt"
)

type sexp interface {
	Eval() (sexp, error)
	String() string
}

func must(x sexp, err error) sexp {
	if err != nil {
		panic(err)
	}
	return x
}

type symbol string

func (sym symbol) Eval() (sexp, error) {
	return sym, nil
}

func (sym symbol) String() string {
	return string(sym)
}

const symNil = symbol("nil")
const symTrue = symbol("t")

type cell struct {
	car sexp
	cdr sexp
}

func (c *cell) Eval() (sexp, error) {
	sym, ok := c.car.(symbol)
	if !ok {
		return nil, fmt.Errorf("invalid function %v", c.car)
	}

	args := c.arguments()
	funcs := map[string]appliable{
		"cons":  function(cons),
		"car":   function(car),
		"cdr":   function(cdr),
		"eq":    function(eq),
		"quote": specialForm(quote),
		"cond":  specialForm(cond),
	}

	fn := funcs[sym.String()]
	if fn == nil {
		return nil, fmt.Errorf("unknown function %v", c.car)
	}

	return fn.Apply(args)
}

func (c *cell) arguments() []sexp {
	rest := c.cdr
	args := []sexp{}
	for rest != symNil {
		args = append(args, must(car(rest)))
		rest = must(cdr(rest))
	}

	return args
}

func (c *cell) String() string {
	return fmt.Sprintf("(%v . %v)", c.car, c.cdr)
}

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
