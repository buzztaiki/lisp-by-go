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
	return apply(env, c.car, c.cdr)
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
