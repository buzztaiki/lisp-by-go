package main

import (
	"fmt"
)

func fnCar(env *environment, args expr) (expr, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}

	return car(car(args)), nil
}

func fnCdr(env *environment, args expr) (expr, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}

	return cdr(car(args)), nil
}

func fnCons(env *environment, args expr) (expr, error) {
	if err := checkArity(args, 2); err != nil {
		return nil, err
	}
	return cons(nth(0, args), nth(1, args)), nil
}

func fnList(env *environment, args expr) (expr, error) {
	return args, nil
}

func fnQuote(env *environment, args expr) (expr, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}
	return car(args), nil
}

func fnEq(env *environment, args expr) (expr, error) {
	if err := checkArity(args, 2); err != nil {
		return nil, err
	}

	if nth(0, args) != nth(1, args) {
		return symNil, nil
	}

	return symTrue, nil
}

func fnCond(env *environment, args expr) (expr, error) {
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

func fnLambda(env *environment, args expr) (expr, error) {
	return newLambdaFunction(env, symLambda.String(), args), nil
}

func fnDefun(env *environment, args expr) (expr, error) {
	if err := checkArityGT(args, 1); err != nil {
		return nil, err
	}

	name := car(args)
	fn := newLambdaFunction(env, name.String(), cdr(args))
	env.funcs[name.String()] = fn

	return name, nil
}

func fnDefmacro(env *environment, args expr) (expr, error) {
	if err := checkArityGT(args, 1); err != nil {
		return nil, err
	}

	name := car(args)
	fn := newMacroForm(env, name.String(), cdr(args))
	env.funcs[name.String()] = fn

	return name, nil
}

func fnPlus(env *environment, args expr) (expr, error) {
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
func fnMinus(env *environment, args expr) (expr, error) {
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

func fnLispApply(env *environment, args expr) (expr, error) {
	if err := checkArityGT(args, 1); err != nil {
		return nil, err
	}

	return car(args).Apply(env, nth(1, args), false)
}

func fnBackquote(env *environment, args expr) (expr, error) {
	return backquoteDoUnquote(env, car(args))
}

func backquoteDoUnquote(env *environment, x expr) (expr, error) {
	if atomp(x) {
		return x, nil
	}
	if x == symNil {
		return x, nil
	}

	xs := []expr{}
	for ; x != symNil; x = cdr(x) {
		head := car(x)
		// `(,'a) => '((, 'a)) => '(a)
		if car(head) == symbol(",") && nth(1, head) != symbol("@") {
			x1, err := nth(1, head).Eval(env)
			if err != nil {
				return nil, err
			}
			xs = append(xs, x1)
			continue
		}
		// `(a ,@'(b)) => `(a (, @) '(b)) => '(a b)
		if car(head) == symbol(",") && nth(1, head) == symbol("@") {
			x = cdr(x)
			x1, err := car(x).Eval(env)
			if err != nil {
				return nil, err
			}
			for ; x1 != symNil; x1 = cdr(x1) {
				xs = append(xs, car(x1))
			}
			continue
		}

		x1, err := backquoteDoUnquote(env, head)
		if err != nil {
			return nil, err
		}
		xs = append(xs, x1)
	}

	return list(xs...), nil
}
