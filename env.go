package main

type environment struct {
	funcs map[string]appliable
	vars  map[string]expr
}

func newEnvironment() *environment {
	return &environment{
		map[string]appliable{
			"cons":   builtinFunction(lispCons),
			"list":   builtinFunction(lispList),
			"car":    builtinFunction(lispCar),
			"cdr":    builtinFunction(lispCdr),
			"eq":     builtinFunction(eq),
			"+":      builtinFunction(plus),
			"-":      builtinFunction(minus),
			"quote":  specialForm(quote),
			"cond":   specialForm(cond),
			"lambda": specialForm(lambda),
			"defun":  specialForm(defun),
		},
		map[string]expr{
			symNil.String():  symNil,
			symTrue.String(): symTrue,
		},
	}
}

func (env *environment) clone() *environment {
	newEnv := &environment{map[string]appliable{}, map[string]expr{}}
	for k, v := range env.funcs {
		newEnv.funcs[k] = v
	}
	for k, v := range env.vars {
		newEnv.vars[k] = v
	}
	return newEnv
}
