package main

type environment struct {
	funcs map[string]function
	vars  map[string]expr
}

func newEnvironment() *environment {
	return &environment{
		map[string]function{
			"cons":     builtinFunction(lispCons),
			"list":     builtinFunction(lispList),
			"car":      builtinFunction(lispCar),
			"cdr":      builtinFunction(lispCdr),
			"eq":       builtinFunction(eq),
			"+":        builtinFunction(plus),
			"-":        builtinFunction(minus),
			"apply":    builtinFunction(lispApply),
			"quote":    specialForm(quote),
			"cond":     specialForm(cond),
			"lambda":   specialForm(lambda),
			"defun":    specialForm(defun),
			"defmacro": specialForm(defmacro),
		},
		map[string]expr{
			symNil.String():  symNil,
			symTrue.String(): symTrue,
		},
	}
}

func (env *environment) clone() *environment {
	newEnv := &environment{map[string]function{}, map[string]expr{}}
	for k, v := range env.funcs {
		newEnv.funcs[k] = v
	}
	for k, v := range env.vars {
		newEnv.vars[k] = v
	}
	return newEnv
}
