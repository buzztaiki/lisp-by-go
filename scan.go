package main

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

type scanner struct {
	src      io.RuneScanner
	prefixes map[rune]bool
}

func newScanner(src io.Reader, prefixes []rune) *scanner {
	ps := map[rune]bool{}
	for _, k := range prefixes {
		ps[k] = true
	}

	return &scanner{bufio.NewReader(src), ps}
}

func (sc *scanner) scan() (string, error) {
	c, err := sc.skipSpaces()
	if err != nil {
		return "", err
	}

	if c == '(' || c == ')' {
		return string(c), nil
	}

	if _, ok := sc.prefixes[c]; ok {
		tok, err := sc.scan()
		if err != nil {
			return "", err
		}
		return string(c) + tok, nil

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
		_, ok := sc.prefixes[c]
		if unicode.IsSpace(c) || c == '(' || c == ')' || ok {
			sc.src.UnreadRune()
			return buf.String(), nil
		}

		buf.WriteRune(c)
	}
}
