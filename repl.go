package main

import (
	"fmt"
	"io"
	"os"
)

func repl(prompt, rprompt string, r io.Reader) {
	env := newEnvironment()
	p := newParser(r)
	for {
		fmt.Print(prompt)
		sexp, err := p.parse()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "read error:", err)
		}

		res, err := sexp.Eval(env)
		if err != nil {
			fmt.Fprintln(os.Stderr, "eval error:", err)
		}
		fmt.Print(rprompt)
		fmt.Println(res)
	}
}
