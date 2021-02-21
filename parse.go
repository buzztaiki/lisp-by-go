package main

import (
	"io"
	"strconv"
)

type parser struct {
	scanner     *scanner
	backlog     string
	translators map[rune]translator
}

func newParser(src io.Reader) *parser {
	translators := map[rune]translator{
		'\'': quoteTranslator,
		'`':  backquoteTranslator,
		',':  unquoteTranslator,
	}

	return &parser{newScanner(src, []rune{'\'', '`', ','}), "", translators}
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

	tokenRunes := []rune(token)
	if len(tokenRunes) == 1 {
		prefix := tokenRunes[0]
		if trans, ok := p.translators[prefix]; ok {
			return trans(p.parseSexp)
		}
	}

	if token == "(" {
		return p.parseList()
	}
	if n, err := strconv.ParseFloat(token, 64); err == nil {
		return number(n), nil
	}
	return symbol(token), nil
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
