package main

import (
	"fmt"
)

type appliable interface {
	Apply(env *environment, args []expr) (expr, error)
}

func evalArgs(env *environment, args []expr) ([]expr, error) {
	newArgs := []expr{}
	for i, arg := range args {
		arg2, err := arg.Eval(env)
		if err != nil {
			return nil, fmt.Errorf("args[%d]: %w", i, err)
		}

		newArgs = append(newArgs, arg2)
	}
	return newArgs, nil
}

type builtinFunction func(env *environment, args []expr) (expr, error)

func (fn builtinFunction) Apply(env *environment, args []expr) (expr, error) {
	newArgs, err := evalArgs(env, args)
	if err != nil {
		return nil, err
	}

	return fn(env, newArgs)
}

type specialForm func(env *environment, args []expr) (expr, error)

func (fn specialForm) Apply(env *environment, args []expr) (expr, error) {
	return fn(env, args)
}

type lambdaFunction struct {
	varNames []string
	body     expr
}

func newLambdaFunction(argsAndBody expr) *lambdaFunction {
	res := &lambdaFunction{[]string{}, symNil}

	varNames, body := split(argsAndBody)
	for varNames != symNil {
		res.varNames = append(res.varNames, car(varNames).String())
		varNames = cdr(varNames)
	}
	res.body = body
	return res
}

func (fn lambdaFunction) Apply(env *environment, args []expr) (expr, error) {
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

	return car(fn.body).Eval(newEnv)
}
