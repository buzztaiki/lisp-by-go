package main

const symNil = symbol("nil")
const symTrue = symbol("t")
const symLambda = symbol("lambda")

func length(xs expr) int {
	if xs == symNil {
		return 0
	}
	return 1 + length(cdr(xs))
}

func car(x expr) expr {
	cell, ok := x.(*cell)
	if !ok || cell == nil {
		return symNil
	}

	return cell.car
}

func cdr(x expr) expr {
	cell, ok := x.(*cell)
	if !ok || cell == nil {
		return symNil
	}

	return cell.cdr
}

func cons(a, b expr) expr {
	return &cell{a, b}
}

func list(exprs ...expr) expr {
	xs := expr(symNil)
	for i := len(exprs) - 1; i >= 0; i-- {
		xs = &cell{exprs[i], xs}
	}
	return xs
}

func nth(n int, xs expr) expr {
	for i := 0; i < n; i++ {
		xs = cdr(xs)
	}
	return car(xs)
}

func mapcar(fn func(x expr) (expr, error), xs expr) (expr, error) {
	res := expr(symNil)

	for ; xs != symNil; xs = cdr(xs) {
		x, err := fn(car(xs))
		if err != nil {
			return nil, err
		}

		res = cons(x, res)
	}

	res2 := expr(symNil)
	for ; res != symNil; res = cdr(res) {
		res2 = cons(car(res), res2)
	}

	return res2, nil
}

func consp(x expr) bool {
	_, ok := x.(*cell)
	return ok
}

func listp(x expr) bool {
	return consp(x) || x == symNil
}

func atomp(x expr) bool {
	return !consp(x)
}
