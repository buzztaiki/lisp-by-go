package main

import (
	"fmt"
)

type function func(env *environment, args expr) (expr, error)

type functionExpr struct {
	name           string
	fn             function
	shouldEvalArgs bool
}

func (fn *functionExpr) Eval(env *environment) (expr, error) {
	return fn, nil
}

func (fn *functionExpr) String() string {
	return fmt.Sprintf("#<function %v>", fn.name)
}

func (fn *functionExpr) Apply(env *environment, args expr, shouldEvalArgs bool) (expr, error) {
	newArgs, err := mapcar(func(x expr) (expr, error) {
		if fn.shouldEvalArgs && shouldEvalArgs {
			return x.Eval(env)
		}
		return x, nil
	}, args)
	if err != nil {
		return nil, fmt.Errorf("argument evaluation failed: %w", err)
	}

	return fn.fn(env, newArgs)
}

func newBuiltinFunction(name string, fn function) *functionExpr {
	return &functionExpr{name, fn, true}
}

func newSpecialForm(name string, fn function) *functionExpr {
	return &functionExpr{name, fn, false}
}

func newLambdaFunction(env *environment, name string, argsAndBody expr) *functionExpr {
	varNames, body := car(argsAndBody), cdr(argsAndBody)

	return &functionExpr{
		name,
		func(_ *environment, args expr) (expr, error) {
			argsEnv, err := newEnvFromArgs(env, varNames, args)
			if err != nil {
				return nil, err
			}

			return car(body).Eval(argsEnv)
		},
		true,
	}
}

func newMacroForm(env *environment, name string, argsAndBody expr) *functionExpr {
	varNames, body := car(argsAndBody), cdr(argsAndBody)

	return &functionExpr{
		name,
		func(evalEnv *environment, args expr) (expr, error) {
			argsEnv, err := newEnvFromArgs(env, varNames, args)
			if err != nil {
				return nil, err
			}

			expanded, err := car(body).Eval(argsEnv)
			if err != nil {
				return nil, fmt.Errorf("macro expantion: %w", err)
			}

			return expanded.Eval(evalEnv)
		},
		false,
	}
}
