package main

const symNil = symbol("nil")
const symTrue = symbol("t")
const symLambda = symbol("lambda")

type environment struct {
	funcs map[string]appliable
	vars  map[string]sexp
}

func newEnvironment() *environment {
	return &environment{
		map[string]appliable{
			"cons":  function(cons),
			"car":   function(car),
			"cdr":   function(cdr),
			"eq":    function(eq),
			"quote": specialForm(quote),
			"cond":  specialForm(cond),
		},
		map[string]sexp{
			symNil.String():  symNil,
			symTrue.String(): symTrue,
		},
	}
}

func (env *environment) clone() *environment {
	newEnv := &environment{map[string]appliable{}, map[string]sexp{}}
	for k, v := range env.funcs {
		newEnv.funcs[k] = v
	}
	for k, v := range env.vars {
		newEnv.vars[k] = v
	}
	return newEnv
}
