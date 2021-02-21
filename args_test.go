package main

import (
	"reflect"
	"testing"
)

func TestNewEnvFromArgs(t *testing.T) {
	cases := []struct {
		name           string
		varNames, args expr
		vars           map[string]expr
		err            error
	}{
		{"ok",
			list(symbol("x")),
			list(symbol("a")),
			map[string]expr{"x": symbol("a")},
			nil},
		{"too long",
			list(symbol("x")),
			list(symbol("a"), symbol("b")),
			nil,
			wronNumberOfArgumentError(2)},
		{"too short",
			list(symbol("x")),
			list(),
			nil,
			wronNumberOfArgumentError(0)},
		{"optional",
			list(symbol("x"), symbol("&optional"), symbol("y"), symbol("z")),
			list(symbol("a"), symbol("b")),
			map[string]expr{"x": symbol("a"), "y": symbol("b"), "z": symNil},
			nil},
		{"rest",
			list(symbol("x"), symbol("&rest"), symbol("args")),
			list(symbol("a"), symbol("b"), symbol("c")),
			map[string]expr{"x": symbol("a"), "args": list(symbol("b"), symbol("c"))},
			nil},
		{"empty rest",
			list(symbol("x"), symbol("&rest"), symbol("args")),
			list(symbol("a")),
			map[string]expr{"x": symbol("a"), "args": list()},
			nil},
		{"rest only",
			list(symbol("&rest"), symbol("args")),
			list(symbol("a")),
			map[string]expr{"args": list(symbol("a"))},
			nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			env := newEnvironment()
			env.vars = map[string]expr{}
			newEnv, err := newEnvFromArgs(env, c.varNames, c.args)
			if err != nil {
				if err.Error() != c.err.Error() {
					t.Errorf("got error %v, want %v", err, c.err)
				}
				return
			}

			if !reflect.DeepEqual(newEnv.vars, c.vars) {
				t.Errorf("got %v, want %v", newEnv.vars, c.vars)
			}
		})
	}
}
