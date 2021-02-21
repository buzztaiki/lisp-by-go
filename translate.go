package main

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
