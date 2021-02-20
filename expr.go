package main

import (
	"fmt"
	"strconv"
)

type expr interface {
	Eval(env *environment) (expr, error)
	Apply(env *environment, args expr, shouldEvalArgs bool) (expr, error)
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

func (sym symbol) Apply(env *environment, args expr, shouldEvalArgs bool) (expr, error) {
	fn := env.funcs[sym.String()]
	if fn == nil {
		return nil, fmt.Errorf("function %v not found", sym)
	}

	res, err := fn.Apply(env, args, shouldEvalArgs)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", sym, err)
	}

	return res, nil
}

func (sym symbol) String() string {
	return string(sym)
}

type number float64

func (num number) Eval(env *environment) (expr, error) {
	return num, nil
}

func (num number) Apply(env *environment, args expr, shouldEvalArgs bool) (expr, error) {
	return nil, fmt.Errorf("number cannot be apply")
}

func (num number) String() string {
	return strconv.FormatFloat(float64(num), 'g', 16, 64)
}

type cell struct {
	car expr
	cdr expr
}

func (c *cell) Eval(env *environment) (expr, error) {
	return c.car.Apply(env, c.cdr, true)
}

func (c *cell) stringNoParen() string {
	switch cdr := c.cdr.(type) {
	case *cell:
		return fmt.Sprintf("%v %v", c.car, cdr.stringNoParen())
	default:
		if cdr == symNil {
			return fmt.Sprint(c.car)
		}
		return fmt.Sprintf("%v . %v", c.car, cdr)
	}
}

func (c *cell) Apply(env *environment, args expr, shouldEvalArgs bool) (expr, error) {
	if c.car != symLambda {
		return nil, fmt.Errorf("invalid function %v", c)
	}

	return newLambdaFunction(symLambda.String(), c.cdr).Apply(env, args, shouldEvalArgs)
}

func (c *cell) String() string {
	return "(" + c.stringNoParen() + ")"
}
