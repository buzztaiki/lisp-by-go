package main

import (
	"fmt"
)

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
	return newLambdaFunction(env, symLambda.String(), args), nil
}

func defun(env *environment, args expr) (expr, error) {
	if err := checkArityGT(args, 1); err != nil {
		return nil, err
	}

	name := car(args)
	fn := newLambdaFunction(env, name.String(), cdr(args))
	env.funcs[name.String()] = fn

	return name, nil
}

func defmacro(env *environment, args expr) (expr, error) {
	if err := checkArityGT(args, 1); err != nil {
		return nil, err
	}

	name := car(args)
	fn := newMacroForm(env, name.String(), cdr(args))
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

func lispApply(env *environment, args expr) (expr, error) {
	if err := checkArityGT(args, 1); err != nil {
		return nil, err
	}

	return car(args).Apply(env, nth(1, args), false)
}
