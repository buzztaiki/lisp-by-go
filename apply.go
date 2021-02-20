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
	newArgs, err := evalArgs(env, args)
	if err != nil {
		return nil, err
	}

	newEnv, err := newEnvFromArgs(env, fn.varNames, newArgs)
	if err != nil {
		return nil, err
	}

	return car(fn.body).Eval(newEnv)
}

type macroForm struct {
	varNames expr
	body     expr
}

func newMacroForm(argsAndBody expr) *macroForm {
	return &macroForm{car(argsAndBody), cdr(argsAndBody)}
}

func (fn macroForm) Apply(env *environment, args expr) (expr, error) {
	newEnv, err := newEnvFromArgs(env, fn.varNames, args)
	if err != nil {
		return nil, err
	}

	expanded, err := car(fn.body).Eval(newEnv)
	if err != nil {
		return nil, fmt.Errorf("macro expantion: %w", err)
	}

	return expanded.Eval(env)
}
