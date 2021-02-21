package main

import "fmt"

type translator func(read func() (expr, error)) (expr, error)

func quoteTranslator(read func() (expr, error)) (expr, error) {
	res, err := read()
	if err != nil {
		return nil, err
	}
	return list(symbol("quote"), res), nil
}

func backquoteTranslator(read func() (expr, error)) (expr, error) {
	res, err := read()
	if err != nil {
		return nil, err
	}

	return list(symbol("backquote"), res), nil
}

func unquoteTranslator(read func() (expr, error)) (expr, error) {
	res, err := read()
	if err != nil {
		return nil, err
	}
	return list(symbol(","), res), nil
}

func hashTranslator(read func() (expr, error)) (expr, error) {
	res, err := read()
	if err != nil {
		return nil, err
	}

	if car(res) != symbol("quote") {
		return nil, fmt.Errorf("invalid syntax: #%v", res)
	}

	return cons(symbol("function"), cdr(res)), nil
}
