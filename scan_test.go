package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	cases := []struct {
		src      string
		expected []string
	}{
		{"(a word 101)", []string{"(", "a", "word", "101", ")"}},
		{"(a (word 101))", []string{"(", "a", "(", "word", "101", ")", ")"}},
		{`(a (word 101 "moo"))`, []string{"(", "a", "(", "word", "101", `"moo"`, ")", ")"}},
		{`a<b a-z :name`, []string{"a<b", "a-z", ":name"}},
		{`"a b c"`, []string{`"a`, `b`, `c"`}}, // とりあえず文字列はなし
		{"'(a b) 'a 'x'y", []string{"'(", "a", "b", ")", "'a", "'x", "'y"}},
	}

	for _, c := range cases {
		t.Run("case:"+c.src, func(t *testing.T) {
			sc := newScanner(strings.NewReader(c.src))
			toks := []string{}
			for {
				tok, err := sc.scan()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatal(err)
				}
				toks = append(toks, tok)
			}
			if !reflect.DeepEqual(c.expected, toks) {
				t.Errorf("want %q, got %q", c.expected, toks)
			}
		})
	}
}
