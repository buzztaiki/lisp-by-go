package main

type environment struct {
	funcs map[string]expr
	vars  map[string]expr
}

func newEnvironment() *environment {
	funcs := map[string]expr{}
	addFunc := func(name string, fn function) { funcs[name] = newBuiltinFunction(name, fn) }
	addSpForm := func(name string, fn function) { funcs[name] = newSpecialForm(name, fn) }

	addFunc("cons", fnCons)
	addFunc("list", fnList)
	addFunc("car", fnCar)
	addFunc("cdr", fnCdr)
	addFunc("eq", fnEq)
	addFunc("consp", fnConsp)
	addFunc("listp", fnListp)
	addFunc("atom", fnAtom)
	addFunc("+", makeNumberAccum(0, func(res, x number) number { return res + x }))
	addFunc("-", makeNumberAccum(0, func(res, x number) number { return res - x }))
	addFunc("*", makeNumberAccum(1, func(res, x number) number { return res * x }))
	addFunc("/", makeNumberAccum(1, func(res, x number) number { return res / x }))
	addFunc("=", makeNumberCmp(func(a, b number) bool { return a == b }))
	addFunc("<", makeNumberCmp(func(a, b number) bool { return a < b }))
	addFunc("<=", makeNumberCmp(func(a, b number) bool { return a <= b }))
	addFunc(">", makeNumberCmp(func(a, b number) bool { return a > b }))
	addFunc(">=", makeNumberCmp(func(a, b number) bool { return a >= b }))
	addFunc("apply", fnLispApply)
	addSpForm("quote", fnQuote)
	addSpForm("backquote", fnBackquote)
	addSpForm("cond", fnCond)
	addSpForm("lambda", fnLambda)
	addSpForm("defun", fnDefun)
	addSpForm("defmacro", fnDefmacro)
	addSpForm("and", fnAnd)
	addSpForm("or", fnOr)
	addSpForm("function", fnFunction)
	addSpForm("let", fnLet)

	return &environment{
		funcs,
		map[string]expr{
			symNil.String():  symNil,
			symTrue.String(): symTrue,
		},
	}
}

func (env *environment) clone() *environment {
	newEnv := &environment{map[string]expr{}, map[string]expr{}}
	for k, v := range env.funcs {
		newEnv.funcs[k] = v
	}
	for k, v := range env.vars {
		newEnv.vars[k] = v
	}
	return newEnv
}
