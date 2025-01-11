package main

import (
	"errors"
	"fmt"
)

func main() {
	wrapSingleVal()
	wrapMultipleVal()
	wrapType()
	wrapByCustomError()
	wrapByCustomErrorMult()
}

func wrapSingleVal() {
	err := fmt.Errorf("foo")

	wrapped := fmt.Errorf("bar: %w", err)

	fmt.Printf("not same = %t\n", err == wrapped)                // not same = false
	fmt.Printf("but is wrapped = %t\n", errors.Is(wrapped, err)) // but is wrapped = true
}

func wrapMultipleVal() {
	err1 := fmt.Errorf("foo")
	err2 := fmt.Errorf("bar")
	err3 := fmt.Errorf("baz")

	wrapped := fmt.Errorf("multiple: %w, %w, %w", err1, err2, err3)

	fmt.Printf(
		"wraps all = %t, %t, %t\n",
		errors.Is(wrapped, err1),
		errors.Is(wrapped, err2),
		errors.Is(wrapped, err3),
	) // wraps all = true, true, true
}

func wrapType() {
	type wrappeeErr struct {
		error
	}

	wrapped := fmt.Errorf("bar: %w", wrappeeErr{fmt.Errorf("foo")})

	var tgt wrappeeErr
	fmt.Printf("wrapped = %t\n", errors.As(wrapped, &tgt)) // wrapped = true
}

type customErr struct {
	Reason string
	Param  any
	Err    error
}

func (e *customErr) Error() string {
	return fmt.Sprintf("*customErr: %s, param = %v", e.Reason, e.Param)
}

func (e *customErr) Unwrap() error {
	return e.Err
}

func wrapByCustomError() {
	err := fmt.Errorf("foo")

	wrapped := &customErr{
		Reason: "bar",
		Param:  "baz",
		Err:    err,
	}

	fmt.Printf("err = %v\n", wrapped)
	// err = *customErr: bar, param = baz
	fmt.Printf("wrapped = %t\n", errors.Is(wrapped, err))
	// wrapped = true
}

type customErr2 struct {
	customErr
}

func (e *customErr2) Error() string {
	return fmt.Sprintf("*customErr2: %s, param = %v", e.Reason, e.Param)
}

func (e *customErr2) Unwrap() []error {
	return []error{e.Err}
}

func wrapByCustomErrorMult() {
	err := fmt.Errorf("foo")

	wrapped := &customErr2{
		customErr{
			Reason: "bar",
			Param:  "baz",
			Err:    err,
		},
	}

	fmt.Printf("err = %v\n", wrapped)
	// err = *customErr2: bar, param = baz
	fmt.Printf("wrapped = %t\n", errors.Is(wrapped, err))
	// wrapped = true
}
