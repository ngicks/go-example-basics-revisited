package main

import (
	"errors"
	"fmt"
	"iter"
	"slices"
)

func reduce[V, Accum any](
	reducer func(accum Accum, next V) Accum,
	initial Accum,
	seq iter.Seq[V],
) Accum {
	accum := initial
	for v := range seq {
		accum = reducer(accum, v)
	}
	return accum
}

func reduceUsingFailableWork[V, Accum any](
	work func(accum Accum, next V) (Accum, error),
	initial Accum,
	seq iter.Seq[V],
) (a Accum, err error) {
	type wrapErr struct {
		err error
	}

	defer func() {
		rec := recover()
		if rec == nil {
			return
		}
		w, ok := rec.(wrapErr)
		if !ok {
			panic(rec)
		}
		err = w.err
	}()

	a = reduce[V, Accum](
		func(accum Accum, next V) Accum {
			var err error
			accum, err = work(accum, next)
			if err != nil {
				panic(wrapErr{err})
			}
			return accum
		},
		initial,
		seq,
	)
	return a, nil
}

func main() {
	sampleErr := errors.New("sample")
	fmt.Println(
		reduceUsingFailableWork(
			func(accum int, next int) (int, error) {
				fmt.Printf("next = %d\n", next)
				if accum > 50 {
					return accum, sampleErr
				}
				return accum + next, nil
			},
			10,
			slices.Values([]int{5, 7, 1, 2}),
		),
	)
	/*
	   next = 5
	   next = 7
	   next = 1
	   next = 2
	   25 <nil>
	*/
	fmt.Println(
		reduceUsingFailableWork(
			func(accum int, next int) (int, error) {
				fmt.Printf("next = %d\n", next)
				if accum > 50 {
					return accum, sampleErr
				}
				return accum + next, nil
			},
			40,
			slices.Values([]int{5, 7, 1, 2}),
		),
	)
	/*
		next = 5
		next = 7
		next = 1
		0 sample
	*/
}
