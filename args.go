package main

import (
	"fmt"
)

type wronNumberOfArgumentError struct {
	nargs int
}

func (err wronNumberOfArgumentError) Error() string {
	return fmt.Sprintf("wrong number of argument %d", err.nargs)
}

func checkArity(args expr, n int) error {
	return checkArityX(args, func() bool { return length(args) == n })
}

func checkArityGT(args expr, n int) error {
	return checkArityX(args, func() bool { return length(args) > n })
}

func checkArityX(args expr, pred func() bool) error {
	if !pred() {
		return wronNumberOfArgumentError{length(args)}
	}
	return nil
}

func evalArgs(env *environment, args expr) (expr, error) {
	return mapcar(func(x expr) (expr, error) {
		return x.Eval(env)
	}, args)
}

func newEnvFromArgs(env *environment, varNames expr, args expr) (*environment, error) {
	newEnv := env.clone()
	nargs := length(args)
	optional := false

	for varNames != symNil {
		if car(varNames) == symbol("&rest") {
			newEnv.vars[car(cdr(varNames)).String()] = args
			return newEnv, nil
		}

		if car(varNames) == symbol("&optional") {
			optional = true
			varNames = cdr(varNames)
		}

		if args == symNil && !optional {
			return nil, wronNumberOfArgumentError{nargs}
		}

		newEnv.vars[car(varNames).String()] = car(args)
		varNames = cdr(varNames)
		args = cdr(args)
	}

	if args != symNil {
		return nil, wronNumberOfArgumentError{nargs}
	}

	return newEnv, nil
}
