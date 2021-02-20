package main

type environment struct {
	funcs map[string]expr
	vars  map[string]expr
}

func newEnvironment() *environment {
	funcs := map[string]expr{}
	addFunc := func(name string, fn function) { funcs[name] = newBuiltinFunction(name, fn) }
	addSpForm := func(name string, fn function) { funcs[name] = newSpecialForm(name, fn) }

	addFunc("cons", lispCons)
	addFunc("list", lispList)
	addFunc("car", lispCar)
	addFunc("cdr", lispCdr)
	addFunc("eq", eq)
	addFunc("+", plus)
	addFunc("-", minus)
	addFunc("apply", lispApply)
	addSpForm("quote", quote)
	addSpForm("cond", cond)
	addSpForm("lambda", lambda)
	addSpForm("defun", defun)
	addSpForm("defmacro", defmacro)

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
