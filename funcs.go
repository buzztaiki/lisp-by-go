package main

import (
	"fmt"
)

func checkArity(args []sexp, n int) error {
	if len(args) != n {
		return fmt.Errorf("wrong number of argument %d", len(args))
	}
	return nil
}

func car(env *environment, args ...sexp) (sexp, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}

	cell, ok := args[0].(*cell)
	if !ok || cell == nil {
		return symNil, nil
	}

	return cell.car, nil
}

func cdr(env *environment, args ...sexp) (sexp, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}

	cell, ok := args[0].(*cell)
	if !ok || cell == nil {
		return symNil, nil
	}

	return cell.cdr, nil
}

func cons(env *environment, args ...sexp) (sexp, error) {
	if err := checkArity(args, 2); err != nil {
		return nil, err
	}

	return &cell{args[0], args[1]}, nil
}

func list(env *environment, args ...sexp) (sexp, error) {
	return joinSexps(args), nil
}

func quote(env *environment, args ...sexp) (sexp, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}
	return args[0], nil
}

func eq(env *environment, args ...sexp) (sexp, error) {
	if err := checkArity(args, 2); err != nil {
		return nil, err
	}

	if args[0] != args[1] {
		return symNil, nil
	}

	return symTrue, nil
}

func cond(env *environment, args ...sexp) (sexp, error) {
	for i, clause := range args {
		cond := must(car(env, clause))
		body := must(cdr(env, clause))

		res, err := cond.Eval(env)
		if err != nil {
			return nil, fmt.Errorf("clauses[%d]: %w", i, err)
		}

		if res != symNil {
			return must(car(env, body)).Eval(env)
		}
	}

	return symNil, nil
}
