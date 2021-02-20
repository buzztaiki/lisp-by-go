# Lisp By Go

Lisp implementation by golang for my interest.

## TODO
- backquote
- string
- implement more functions

## Example

```console
go run .
>
(defun fib (n)
  (cond ((eq n 0) 0)
        ((eq n 1) 1)
        (t (+ (fib (- n 1)) (fib (- n 2))))))
==> fib
> (fib 10)
==> 55
```

```console
% go run .
>
(defmacro my-if (cond then else)
  (list 'cond
        (list cond then)
        (list t else)))
==> my-if
> (my-if (eq 1 1) 'ok moo)
==> ok
> (my-if (eq 1 2) moo 'ng)
==> ng
```

```console
% go run .
>
(defun my-append (&rest seqs)
  (cond ((eq seqs nil) nil)
        ((eq (car seqs) nil) (apply 'my-append (cdr seqs)))
        (t (cons (car (car seqs))
                 (apply 'my-append (cons (cdr (car seqs))
                                         (cdr seqs)))))))
==> my-append
> (my-append '(1 2 3) '(4 5 6) '(7 8 9))
==> (1 2 3 4 5 6 7 8 9)
```
