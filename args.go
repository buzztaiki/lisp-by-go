package main

func checkArity(args expr, n int) error {
	return checkArityX(args, func() bool { return length(args) == n })
}

func checkArityX(args expr, pred func() bool) error {
	if !pred() {
		return fmt.Errorf("wrong number of argument %d", length(args))
	}
	return nil
}

func evalArgs(env *environment, args expr) (expr, error) {
	return mapcar(func(x expr) (expr, error) {
		return x.Eval(env)
	}, args)
}

func newEnvFromArgs(env *environment, varNames expr, args expr) *environment {
	newEnv := env.clone()
	for args != symNil && varNames != symNil {
		newEnv.vars[car(varNames).String()] = car(args)

		args = cdr(args)
		varNames = cdr(varNames)
	}

	return newEnv
}
