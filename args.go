package main

import (
	"fmt"
)

func wronNumberOfArgumentError(nargs int) error {
	return fmt.Errorf("wrong number of argument %d", nargs)
}

func wrongNumberTypeArgumentError(arg expr) error {
	return fmt.Errorf("wrong number type argument %v", arg)
}

func checkArity(args expr, n int) error {
	return checkArityX(args, func() bool { return length(args) == n })
}

func checkArityGT(args expr, n int) error {
	return checkArityX(args, func() bool { return length(args) > n })
}

func checkArityX(args expr, pred func() bool) error {
	if !pred() {
		return wronNumberOfArgumentError(length(args))
	}
	return nil
}

func checkNumber(arg expr) (number, error) {
	num, ok := arg.(number)
	if !ok {
		return number(0), wrongNumberTypeArgumentError(arg)
	}
	return num, nil
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
			return nil, wronNumberOfArgumentError(nargs)
		}

		newEnv.vars[car(varNames).String()] = car(args)
		varNames = cdr(varNames)
		args = cdr(args)
	}

	if args != symNil {
		return nil, wronNumberOfArgumentError(nargs)
	}
	return newEnv, nil
}
