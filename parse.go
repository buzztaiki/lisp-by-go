package main

import (
	"io"
	"strconv"
	"strings"
)

type parser struct {
	scanner *scanner
	backlog string
}

func newParser(src io.Reader) *parser {
	return &parser{newScanner(src), ""}
}

func (p *parser) parse() (expr, error) {
	return p.parseSexp()
}

func (p *parser) scan() (string, error) {
	if p.backlog != "" {
		token := p.backlog
		p.backlog = ""
		return token, nil
	}
	return p.scanner.scan()
}

func (p *parser) unscan(token string) {

	p.backlog = token
}

func (p *parser) parseSexp() (expr, error) {
	token, err := p.scan()
	if err != nil {
		return nil, err
	}

	fn := func(token string) (expr, error) {
		if token == "(" {
			return p.parseList()
		}
		if n, err := strconv.ParseFloat(token, 64); err == nil {
			return number(n), nil
		}
		return symbol(token), nil
	}

	quoted := strings.HasPrefix(token, "'")
	if !quoted {
		return fn(token)
	}

	res, err := fn(token[1:])
	if err != nil {
		return nil, err
	}
	return list(symbol("quote"), res), nil
}

func (p *parser) parseList() (expr, error) {
	token, err := p.scanner.scan()
	if err != nil {
		return nil, err
	}

	if token == ")" {
		return symNil, nil
	}

	if token == "." {
		rest, err := p.parseList()
		if err != nil {
			return nil, err
		}
		return car(rest), nil
	}

	p.unscan(token)
	expr, err := p.parseSexp()
	if err != nil {
		return nil, err
	}

	rest, err := p.parseList()
	if err != nil {
		return nil, err
	}

	return cons(expr, rest), nil
}
