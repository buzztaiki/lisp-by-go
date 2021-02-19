package main

import (
	"fmt"
)

func apply(env *environment, fnExpr expr, argsExpr expr) (expr, error) {
	switch x := fnExpr.(type) {
	case symbol:
		fn := env.funcs[x.String()]
		if fn == nil {
			return nil, fmt.Errorf("unknown function %v", fnExpr)
		}

		res, err := fn.Apply(env, argsExpr)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", x, err)
		}
		return res, nil
	case *cell:
		if x.car != symLambda {
			return nil, fmt.Errorf("invalid function %v", fnExpr)
		}

		res, err := newLambdaFunction(x.cdr).Apply(env, argsExpr)
		if err != nil {
			return nil, fmt.Errorf("lambda: %w", err)
		}
		return res, nil
	}
	return nil, fmt.Errorf("invalid function %v", fnExpr)
}

type appliable interface {
	Apply(env *environment, args expr) (expr, error)
}

func evalArgs(env *environment, args expr) (expr, error) {
	return mapcar(func(x expr) (expr, error) {
		return x.Eval(env)
	}, args)
}

type builtinFunction func(env *environment, args expr) (expr, error)

func (fn builtinFunction) Apply(env *environment, args expr) (expr, error) {
	newArgs, err := evalArgs(env, args)

	if err != nil {
		return nil, err
	}

	return fn(env, newArgs)
}

type specialForm func(env *environment, args expr) (expr, error)

func (fn specialForm) Apply(env *environment, args expr) (expr, error) {
	return fn(env, args)
}

type lambdaFunction struct {
	varNames expr
	body     expr
}

func newLambdaFunction(argsAndBody expr) *lambdaFunction {
	return &lambdaFunction{car(argsAndBody), cdr(argsAndBody)}
}

func (fn lambdaFunction) Apply(env *environment, args expr) (expr, error) {
	if err := checkArity(args, length(fn.varNames)); err != nil {
		return nil, err
	}

	newArgs, err := evalArgs(env, args)
	if err != nil {
		return nil, err
	}

	newEnv := env.clone()
	varNames := fn.varNames
	for newArgs != symNil && varNames != symNil {
		newEnv.vars[car(varNames).String()] = car(newArgs)

		newArgs = cdr(newArgs)
		varNames = cdr(varNames)
	}

	return car(fn.body).Eval(newEnv)
}
