package main

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

type scanner struct {
	src io.RuneScanner
}

func newScanner(src io.Reader) *scanner {
	return &scanner{bufio.NewReader(src)}
}

func (sc *scanner) scan() (string, error) {
	c, err := sc.skipSpaces()
	if err != nil {
		return "", err
	}

	if c == '(' || c == ')' {
		return string(c), nil
	}

	if c == '\'' {
		tok, err := sc.scan()
		if err != nil {
			return "", err
		}
		return "'" + tok, nil

	}

	return sc.scanToken(c)
}

func (sc *scanner) skipSpaces() (rune, error) {
	for {
		c, _, err := sc.src.ReadRune()
		if err != nil {
			return 0, err
		}
		if !unicode.IsSpace(c) {
			return c, nil
		}
	}
}

func (sc *scanner) scanToken(head rune) (string, error) {
	buf := &strings.Builder{}
	buf.WriteRune(head)

	for {
		c, _, err := sc.src.ReadRune()
		if err == io.EOF {
			return buf.String(), nil
		}
		if err != nil {
			return "", err
		}
		if unicode.IsSpace(c) || c == '(' || c == ')' || c == '\'' {
			sc.src.UnreadRune()
			return buf.String(), nil
		}

		buf.WriteRune(c)
	}
}
