package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	cases := []struct {
		src  string
		want sexp
	}{
		{"(cons a b)", list(symbol("cons"), symbol("a"), symbol("b"))},
		{"(cons a (cons b c))", list(symbol("cons"), symbol("a"), list(symbol("cons"), symbol("b"), symbol("c")))},
		{"(+ 10 20)", list(symbol("+"), number(10), number(20))},
	}

	for _, c := range cases {
		t.Run("case:"+c.src, func(t *testing.T) {
			p := newParser(strings.NewReader(c.src))
			sexp, err := p.parse()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(sexp, c.want) {
				t.Errorf("want %v, got %v", c.want, sexp)
			}
		})
	}
}
