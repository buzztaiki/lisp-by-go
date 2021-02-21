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
		expr, err := p.parse()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "read error:", err)
			continue
		}

		res, err := expr.Eval(env)
		if err != nil {
			fmt.Fprintln(os.Stderr, "eval error:", err)
			continue
		}
		fmt.Print(rprompt)
		fmt.Println(res)
	}
}
