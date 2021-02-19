package main

import (
	"fmt"
)

type appliable interface {
	Apply(env *environment, args []sexp) (sexp, error)
}

func evalArgs(env *environment, args []sexp) ([]sexp, error) {
	newArgs := []sexp{}
	for i, arg := range args {
		arg2, err := arg.Eval(env)
		if err != nil {
			return nil, fmt.Errorf("args[%d]: %w", i, err)
		}

		newArgs = append(newArgs, arg2)
	}
	return newArgs, nil
}

type function func(env *environment, args ...sexp) (sexp, error)

func (fn function) Apply(env *environment, args []sexp) (sexp, error) {
	newArgs, err := evalArgs(env, args)
	if err != nil {
		return nil, err
	}

	return fn(env, newArgs...)
}

type specialForm func(env *environment, args ...sexp) (sexp, error)

func (fn specialForm) Apply(env *environment, args []sexp) (sexp, error) {
	return fn(env, args...)
}

type lambdaFunction struct {
	varNames []string
	body     sexp
}

func newLambdaFunction(env *environment, x sexp) *lambdaFunction {
	res := &lambdaFunction{[]string{}, symNil}

	argsAndBody, ok := x.(*cell)
	if ok {
		varNames, body := must(car(env, argsAndBody)), must(cdr(env, argsAndBody))
		for varNames != symNil {
			res.varNames = append(res.varNames, must(car(env, varNames)).String())
			varNames = must(cdr(env, varNames))
		}
		res.body = body
	}
	return res
}

func (fn lambdaFunction) Apply(env *environment, args []sexp) (sexp, error) {
	if err := checkArity(args, len(fn.varNames)); err != nil {
		return nil, err
	}

	newArgs, err := evalArgs(env, args)
	if err != nil {
		return nil, err
	}

	newEnv := env.clone()
	for i := range args {
		newEnv.vars[fn.varNames[i]] = newArgs[i]
	}

	return must(car(env, fn.body)).Eval(newEnv)
}
