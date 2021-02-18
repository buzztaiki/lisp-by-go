package main

func car(sexp sexp) sexp {
	cell := sexp.(*cell)
	if cell == nil {
		return symNil
	}

	return cell.car
}

func cdr(sexp sexp) sexp {
	cell := sexp.(*cell)
	if cell == nil {
		return symNil
	}

	return cell.cdr
}

func cons(sexp sexp) sexp {
	return &cell{car(sexp), cdr(sexp)}
}
