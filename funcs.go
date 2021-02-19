package main

import "fmt"

func checkArity(args []sexp, n int) error {
	if len(args) != n {
		return fmt.Errorf("wrong number of argument %d", len(args))
	}
	return nil
}

func car(args ...sexp) (sexp, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}

	cell, ok := args[0].(*cell)
	if !ok || cell == nil {
		return symNil, nil
	}

	return cell.car, nil
}

func cdr(args ...sexp) (sexp, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}

	cell, ok := args[0].(*cell)
	if !ok || cell == nil {
		return symNil, nil
	}

	return cell.cdr, nil
}

func cons(args ...sexp) (sexp, error) {
	if err := checkArity(args, 2); err != nil {
		return nil, err
	}

	return &cell{args[0], args[1]}, nil
}

func list(sexps ...sexp) (sexp, error) {
	xs := sexp(symNil)
	for i := len(sexps) - 1; i >= 0; i-- {
		xs = &cell{sexps[i], xs}
	}
	return xs, nil
}

func quote(args ...sexp) (sexp, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}
	return args[0], nil
}
