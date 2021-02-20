package main

import (
	"fmt"
)

func findFunction(env *environment, fnExpr expr) (appliable, error) {
	switch x := fnExpr.(type) {
	case symbol:
		fn := env.funcs[x.String()]
		if fn == nil {
			return nil, fmt.Errorf("unknown function %v", fnExpr)
		}

		return fn, nil
	case *cell:
		if x.car != symLambda {
			return nil, fmt.Errorf("invalid function %v", fnExpr)
		}

		return newLambdaFunction(x.cdr), nil
	default:
		return nil, fmt.Errorf("invalid function %v", fnExpr)
	}
}

func apply(env *environment, fnExpr expr, args expr, shouldEvalArgs bool) (expr, error) {
	fn, err := findFunction(env, fnExpr)
	if err != nil {
		return nil, err
	}

	newArgs, err := mapcar(func(x expr) (expr, error) {
		if shouldEvalArgs && fn.ShouldEvalArgs() {
			return x.Eval(env)
		}
		return x, nil
	}, args)
	if err != nil {
		return nil, fmt.Errorf("argument evaluation failed: %w", err)
	}

	res, err := fn.Apply(env, newArgs)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", fnExpr, err)
	}
	return res, nil
}

type appliable interface {
	Apply(env *environment, args expr) (expr, error)
	ShouldEvalArgs() bool
}

type builtinFunction func(env *environment, args expr) (expr, error)

func (fn builtinFunction) Apply(env *environment, args expr) (expr, error) {
	return fn(env, args)
}

func (fn builtinFunction) ShouldEvalArgs() bool {
	return true
}

type specialForm func(env *environment, args expr) (expr, error)

func (fn specialForm) Apply(env *environment, args expr) (expr, error) {
	return fn(env, args)
}

func (fn specialForm) ShouldEvalArgs() bool {
	return false
}

type lambdaFunction struct {
	varNames expr
	body     expr
}

func newLambdaFunction(argsAndBody expr) *lambdaFunction {
	return &lambdaFunction{car(argsAndBody), cdr(argsAndBody)}
}

func (fn lambdaFunction) Apply(env *environment, args expr) (expr, error) {
	newEnv, err := newEnvFromArgs(env, fn.varNames, args)
	if err != nil {
		return nil, err
	}

	return car(fn.body).Eval(newEnv)
}

func (fn lambdaFunction) ShouldEvalArgs() bool {
	return true
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

func (fn macroForm) ShouldEvalArgs() bool {
	return false
}
