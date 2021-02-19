package main

import (
	"fmt"
)

func checkArity(args []sexp, n int) error {
	return checkArityX(args, func() bool { return len(args) == n })
}

func checkArityX(args []sexp, pred func() bool) error {
	if !pred() {
		return fmt.Errorf("wrong number of argument %d", len(args))
	}
	return nil
}

func car(x sexp) sexp {
	cell, ok := x.(*cell)
	if !ok || cell == nil {
		return symNil
	}

	return cell.car
}

func cdr(x sexp) sexp {
	cell, ok := x.(*cell)
	if !ok || cell == nil {
		return symNil
	}

	return cell.cdr
}

func split(x sexp) (sexp, sexp) {
	return car(x), cdr(x)
}

func lispCar(env *environment, args []sexp) (sexp, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}

	return car(args[0]), nil
}

func lispCdr(env *environment, args []sexp) (sexp, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}

	return cdr(args[0]), nil
}

func cons(a, b sexp) sexp {
	return &cell{a, b}
}

func lispCons(env *environment, args []sexp) (sexp, error) {
	if err := checkArity(args, 2); err != nil {
		return nil, err
	}
	return cons(args[0], args[1]), nil
}

func list(sexps ...sexp) sexp {
	xs := sexp(symNil)
	for i := len(sexps) - 1; i >= 0; i-- {
		xs = &cell{sexps[i], xs}
	}
	return xs
}

func lispList(env *environment, args []sexp) (sexp, error) {
	return list(args...), nil
}

func quote(env *environment, args []sexp) (sexp, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}
	return args[0], nil
}

func eq(env *environment, args []sexp) (sexp, error) {
	if err := checkArity(args, 2); err != nil {
		return nil, err
	}

	if args[0] != args[1] {
		return symNil, nil
	}

	return symTrue, nil
}

func cond(env *environment, args []sexp) (sexp, error) {
	for i, clause := range args {
		cond, body := split(clause)

		res, err := cond.Eval(env)
		if err != nil {
			return nil, fmt.Errorf("clauses[%d]: %w", i, err)
		}

		if res != symNil {
			return car(body).Eval(env)
		}
	}

	return symNil, nil
}

func lambda(env *environment, args []sexp) (sexp, error) {
	return cons(symLambda, list(args...)), nil
}

func defun(env *environment, args []sexp) (sexp, error) {
	if err := checkArityX(args, func() bool { return len(args) > 1 }); err != nil {
		return nil, err
	}

	name := args[0]
	fn := newLambdaFunction(list(args[1:]...))
	env.funcs[name.String()] = fn

	return name, nil
}

func plus(env *environment, args []sexp) (sexp, error) {
	res := float64(0)
	for i, x := range args {
		num, ok := x.(number)
		if !ok {
			return nil, fmt.Errorf("args[%d]: wrong number type argument %v", i, x)
		}
		res += float64(num)
	}
	return number(res), nil
}

// 引数が 1 つの場合はその負数を返す。
// 二つ以上の場合は引き算する。
func minus(env *environment, args []sexp) (sexp, error) {
	res := float64(0)

	for i, x := range args {
		num, ok := x.(number)
		if !ok {
			return nil, fmt.Errorf("args[%d]: wrong number type argument %v", i, x)
		}

		// 引数が複数ある場合は最初の値をひっくりかえす
		if i == 1 {
			res *= -1
		}
		res -= float64(num)
	}
	return number(res), nil
}
