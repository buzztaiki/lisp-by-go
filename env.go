package main

const symNil = symbol("nil")
const symTrue = symbol("t")

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
