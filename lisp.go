package main

import (
	"fmt"
)

type sexp interface {
	Eval() (sexp, error)
	String() string
}

type cell struct {
	car sexp
	cdr sexp
}

func (cell *cell) Eval() (sexp, error) {
	sym, ok := cell.car.(symbol)
	if !ok {
		return nil, fmt.Errorf("invalid function %v", cell.car)
	}

	rest := cell.cdr

	switch sym.String() {
	case "cons":
		return cons(rest), nil
	case "car":
		return car(rest), nil
	case "cdr":
		return cdr(rest), nil
	default:
		return nil, fmt.Errorf("unknown function %v", cell.car)
	}
}

func (cell *cell) String() string {
	return fmt.Sprintf("(%v . %v)", cell.car, cell.cdr)
}

type symbol string

func (sym symbol) Eval() (sexp, error) {
	return sym, nil
}

func (sym symbol) String() string {
	return string(sym)
}

const symNil = symbol("nil")

func sexpList(sexps ...sexp) sexp {
	list := sexp(symNil)
	for i := len(sexps) - 1; i >= 0; i-- {
		list = &cell{sexps[i], list}
	}
	return list
}
