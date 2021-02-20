package main

import (
	"fmt"
)

func checkArity(args expr, n int) error {
	return checkArityX(args, func() bool { return length(args) == n })
}

func checkArityX(args expr, pred func() bool) error {
	if !pred() {
		return fmt.Errorf("wrong number of argument %d", length(args))
	}
	return nil
}

func length(xs expr) int {
	if xs == symNil {
		return 0
	}
	return 1 + length(cdr(xs))
}

func car(x expr) expr {
	cell, ok := x.(*cell)
	if !ok || cell == nil {
		return symNil
	}

	return cell.car
}

func cdr(x expr) expr {
	cell, ok := x.(*cell)
	if !ok || cell == nil {
		return symNil
	}

	return cell.cdr
}

func cons(a, b expr) expr {
	return &cell{a, b}
}

func list(exprs ...expr) expr {
	xs := expr(symNil)
	for i := len(exprs) - 1; i >= 0; i-- {
		xs = &cell{exprs[i], xs}
	}
	return xs
}

func nth(n int, xs expr) expr {
	for i := 0; i < n; i++ {
		xs = cdr(xs)
	}
	return car(xs)
}

func mapcar(fn func(x expr) (expr, error), xs expr) (expr, error) {
	res := expr(symNil)

	for ; xs != symNil; xs = cdr(xs) {
		x, err := fn(car(xs))
		if err != nil {
			return nil, err
		}

		res = cons(x, res)
	}

	res2 := expr(symNil)
	for ; res != symNil; res = cdr(res) {
		res2 = cons(car(res), res2)
	}

	return res2, nil
}

func lispCar(env *environment, args expr) (expr, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}

	return car(car(args)), nil
}

func lispCdr(env *environment, args expr) (expr, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}

	return cdr(car(args)), nil
}

func lispCons(env *environment, args expr) (expr, error) {
	if err := checkArity(args, 2); err != nil {
		return nil, err
	}
	return cons(nth(0, args), nth(1, args)), nil
}

func lispList(env *environment, args expr) (expr, error) {
	return args, nil
}

func quote(env *environment, args expr) (expr, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}
	return car(args), nil
}

func eq(env *environment, args expr) (expr, error) {
	if err := checkArity(args, 2); err != nil {
		return nil, err
	}

	if nth(0, args) != nth(1, args) {
		return symNil, nil
	}

	return symTrue, nil
}

func cond(env *environment, args expr) (expr, error) {
	for ; args != symNil; args = cdr(args) {
		clause := car(args)
		cond, body := car(clause), cdr(clause)

		res, err := cond.Eval(env)
		if err != nil {
			return nil, err
		}

		if res != symNil {
			return car(body).Eval(env)
		}
	}

	return symNil, nil
}

func lambda(env *environment, args expr) (expr, error) {
	return cons(symLambda, args), nil
}

func defun(env *environment, args expr) (expr, error) {
	if err := checkArityX(args, func() bool { return length(args) > 1 }); err != nil {
		return nil, err
	}

	name := car(args)
	fn := newLambdaFunction(cdr(args))
	env.funcs[name.String()] = fn

	return name, nil
}

func defmacro(env *environment, args expr) (expr, error) {
	if err := checkArityX(args, func() bool { return length(args) > 1 }); err != nil {
		return nil, err
	}

	name := car(args)
	fn := newMacroForm(cdr(args))
	env.funcs[name.String()] = fn

	return name, nil
}

func plus(env *environment, args expr) (expr, error) {
	res := float64(0)
	for ; args != symNil; args = cdr(args) {
		num, ok := car(args).(number)
		if !ok {
			return nil, fmt.Errorf("wrong number type argument %v", car(args))
		}
		res += float64(num)
	}
	return number(res), nil
}

// 引数が 1 つの場合はその負数を返す。
// 二つ以上の場合は引き算する。
func minus(env *environment, args expr) (expr, error) {
	res := float64(0)

	for i := 0; args != symNil; args = cdr(args) {
		num, ok := car(args).(number)
		if !ok {
			return nil, fmt.Errorf("wrong number type argument %v", car(args))
		}

		// 引数が複数ある場合は最初の値をひっくりかえす
		if i == 1 {
			res *= -1
		}
		res -= float64(num)
		i++
	}
	return number(res), nil
}
