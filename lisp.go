package main

import (
	"fmt"
	"strconv"
)

type expr interface {
	Eval(env *environment) (expr, error)
	String() string
}

type symbol string

func (sym symbol) Eval(env *environment) (expr, error) {
	x := env.vars[sym.String()]
	if x == nil {
		return nil, fmt.Errorf("variable %v not found", sym)
	}

	return x, nil
}

func (sym symbol) String() string {
	return string(sym)
}

const symNil = symbol("nil")
const symTrue = symbol("t")
const symLambda = symbol("lambda")

type number float64

func (num number) Eval(env *environment) (expr, error) {
	return num, nil
}

func (num number) String() string {
	return strconv.FormatFloat(float64(num), 'g', 16, 64)
}

type cell struct {
	car expr
	cdr expr
}

func (c *cell) Eval(env *environment) (expr, error) {
	args := c.arguments(env)

	switch x := c.car.(type) {
	case symbol:
		fn := env.funcs[x.String()]
		if fn == nil {
			return nil, fmt.Errorf("unknown function %v", c.car)
		}

		return fn.Apply(env, args)
	case *cell:
		if x.car != symLambda {
			return nil, fmt.Errorf("invalid function %v", c.car)
		}
		return newLambdaFunction(x.cdr).Apply(env, args)
	}
	return nil, fmt.Errorf("invalid function %v", c.car)
}

func (c *cell) arguments(env *environment) []expr {
	rest := c.cdr
	args := []expr{}
	for rest != symNil {
		args = append(args, car(rest))
		rest = cdr(rest)
	}

	return args
}

func (c *cell) stringNoParen() string {
	switch x := c.cdr.(type) {
	case *cell:
		return fmt.Sprintf("%v %v", c.car, x.stringNoParen())
	default:
		if x == symNil {
			return fmt.Sprint(c.car)
		}
		return fmt.Sprintf("%v . %v", c.car, x)
	}
}

func (c *cell) String() string {
	return "(" + c.stringNoParen() + ")"
}
