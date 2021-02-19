package main

import (
	"io"
	"strconv"
)

type parser struct {
	scanner *scanner
	backlog string
}

func newParser(src io.Reader) *parser {
	return &parser{newScanner(src), ""}
}

func (p *parser) parse() (sexp, error) {
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

func (p *parser) parseSexp() (sexp, error) {
	token, err := p.scan()
	if err != nil {
		return nil, err
	}

	if token == "(" {
		return p.parseList()
	}
	if n, err := strconv.ParseFloat(token, 64); err == nil {
		return number(n), nil
	}
	return symbol(token), nil
}

func (p *parser) parseList() (sexp, error) {
	sexps := []sexp{}

	for {
		token, err := p.scanner.scan()
		if err != nil {
			return nil, err
		}

		if token == ")" {
			return list(sexps...), nil
		}

		p.unscan(token)
		sexp, err := p.parseSexp()
		if err != nil {
			return nil, err
		}
		sexps = append(sexps, sexp)
	}
}
