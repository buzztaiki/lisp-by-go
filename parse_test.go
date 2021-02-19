package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	cases := []struct {
		src  string
		want expr
	}{
		{"(cons a b)", list(symbol("cons"), symbol("a"), symbol("b"))},
		{"(cons a (cons b c))", list(symbol("cons"), symbol("a"), list(symbol("cons"), symbol("b"), symbol("c")))},
		{"(+ 10 20)", list(symbol("+"), number(10), number(20))},
		{"'a", list(symbol("quote"), symbol("a"))},
		{"'(a b)", list(symbol("quote"), list(symbol("a"), symbol("b")))},
		{"(a b . c)", cons(symbol("a"), cons(symbol("b"), symbol("c")))},
	}

	for _, c := range cases {
		t.Run("case:"+c.src, func(t *testing.T) {
			p := newParser(strings.NewReader(c.src))
			expr, err := p.parse()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(expr, c.want) {
				t.Errorf("want %v, got %v", c.want, expr)
			}
		})
	}
}
