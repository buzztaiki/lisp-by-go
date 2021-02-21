package main

func boolToExpr(x bool) expr {
	if x {
		return symTrue
	}
	return symNil
}

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
	return boolToExpr(nth(0, args) == nth(1, args)), nil
}

func fnConsp(env *environment, args expr) (expr, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}
	return boolToExpr(consp(car(args))), nil
}

func fnListp(env *environment, args expr) (expr, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}
	return boolToExpr(listp(car(args))), nil
}

func fnAtom(env *environment, args expr) (expr, error) {
	if err := checkArity(args, 1); err != nil {
		return nil, err
	}
	return boolToExpr(atomp(car(args))), nil
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

func makeNumberAccum(first number, fn func(res, x number) number) function {
	return func(env *environment, args expr) (expr, error) {
		res := first
		if length(args) > 1 {
			num, err := checkNumber(car(args))
			if err != nil {
				return nil, err
			}
			res = num
			args = cdr(args)
		}

		for ; args != symNil; args = cdr(args) {
			num, err := checkNumber(car(args))
			if err != nil {
				return nil, err
			}

			res = fn(res, num)
		}
		return res, nil
	}
}

func makeNumberCmp(pred func(a, b number) bool) function {
	return func(env *environment, args expr) (expr, error) {
		if err := checkArity(args, 2); err != nil {
			return nil, err
		}

		a, err := checkNumber(car(args))
		if err != nil {
			return nil, err
		}
		b, err := checkNumber(car(cdr(args)))
		if err != nil {
			return nil, err
		}
		return boolToExpr(pred(a, b)), nil
	}
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
