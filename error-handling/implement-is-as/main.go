package main

import (
	"errors"
	"fmt"
)

func main() {
	{
		err1 := &errBare{Msg: "is", Kind: ErrKind1}
		fmt.Printf("is = %t\n", errors.Is(err1, ErrKind1))                        // is = false
		fmt.Printf("is = %t\n", errors.Is(fmt.Errorf("wrapped; %w", err1), err1)) // is = true

		err2 := &errIs{errBare{Msg: "is", Kind: ErrKind1}}
		fmt.Printf("is = %t\n", errors.Is(err2, ErrKind1))                        // is = true
		fmt.Printf("is = %t\n", errors.Is(err2, ErrKind2))                        // is = false
		fmt.Printf("is = %t\n", errors.Is(fmt.Errorf("wrapped: %w", err2), err2)) // is = true
	}
	{
		err1 := &errBare{Msg: "is", Kind: ErrKind1}
		var kind ErrKind
		fmt.Printf("as = %t\n", errors.As(err1, &kind)) // as = false

		err2 := &errAs{errBare{Msg: "is", Kind: ErrKind1}}
		kind = ""
		fmt.Printf("as = %t, kind = %s\n", errors.As(err2, &kind), kind) // as = true, kind = kind 1
		kind = ""
		fmt.Printf("as = %t, kind = %s\n", errors.As(fmt.Errorf("wrapped: %w", err2), &kind), kind) // as = true, kind = kind 1
	}
}

type ErrKind string

func (e ErrKind) Error() string {
	return string(e)
}

var (
	ErrKind1 = ErrKind("kind 1")
	ErrKind2 = ErrKind("kind 2")
)

type errBare struct {
	Msg  string
	Kind ErrKind
}

func (e *errBare) Error() string {
	return e.Msg
}

type errIs struct {
	errBare
}

func (e *errIs) Is(err error) bool {
	if k, ok := err.(ErrKind); ok {
		return e.Kind == k
	}
	return false
}

type errAs struct {
	errBare
}

func (e *errAs) As(tgt any) bool {
	if k, ok := tgt.(*ErrKind); ok {
		*k = e.Kind
		return true
	}
	return false
}
