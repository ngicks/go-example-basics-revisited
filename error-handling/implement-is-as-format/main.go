package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	{
		err1 := &errBare{Msg: "is", Kind: ErrKind1}
		fmt.Printf("is = %t\n", errors.Is(err1, ErrKind1))                        // is = false
		fmt.Printf("is = %t\n", errors.Is(fmt.Errorf("wrapped; %w", err1), err1)) // is = true

		err2 := &errIs{errBare{Msg: "is", Kind: ErrKind1 | ErrKind2}}
		fmt.Printf("is = %t\n", errors.Is(err2, ErrKind1))                        // is = true
		fmt.Printf("is = %t\n", errors.Is(err2, ErrKind2))                        // is = true
		fmt.Printf("is = %t\n", errors.Is(fmt.Errorf("wrapped: %w", err2), err2)) // is = true
	}
	{
		err1 := &errBare{Msg: "is", Kind: ErrKind1}
		var kind ErrKind
		fmt.Printf("as = %t\n", errors.As(err1, &kind)) // as = false

		err2 := &errAs{errBare{Msg: "is", Kind: ErrKind1}}
		kind = 0
		fmt.Printf("as = %t, kind = %s\n", errors.As(err2, &kind), kind) // as = true, kind = kind 1
		kind = 0
		fmt.Printf("as = %t, kind = %s\n", errors.As(fmt.Errorf("wrapped: %w", err2), &kind), kind) // as = true, kind = kind 1

		err2 = &errAs{errBare{Msg: "is", Kind: ErrKind1 | ErrKind2}}
		kind = 0
		fmt.Printf("as = %t, kind = %s\n", errors.As(err2, &kind), kind) // as = true, kind = kind 1&2
	}
	{
		bare := &errBare{Msg: "is", Kind: ErrKind1}
		for _, v := range "vTtbcdoOqxXUeEfFgGsqp" {
			verb := string([]rune{v})
			fmt.Printf("verb %%%s, bare = %"+verb+", format = %"+verb+"\n", verb, bare, &errFormat{*bare})
			for _, f := range " +-#0" {
				verb := string([]rune{f, v})
				fmt.Printf("verb %%%s, bare = %"+verb+", format = %"+verb+"\n", verb, bare, &errFormat{*bare})
			}
		}
	}
}

type ErrKind int

func (e ErrKind) Error() string {
	var s strings.Builder
	s.WriteString("kind ")
	var count int
	for i := 0; i < 32; i++ {
		if e&(1<<i) > 0 {
			if count > 0 {
				s.WriteByte('&')
			}
			s.WriteString(strconv.Itoa(i + 1))
			count++
		}
	}
	return s.String()
}

const (
	ErrKind1 = ErrKind(1 << iota)
	ErrKind2
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
		return e.Kind&k > 0
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

type errFormat struct {
	errBare
}

func (e *errFormat) Format(state fmt.State, verb rune) {
	if verb == 'v' {
		if state.Flag('+') {
			_, _ = fmt.Fprintf(state, "msg = %s, kind = %s", e.Msg, e.Kind)
			return
		}
	}
	// plain does not inherit method from errFormat.
	type plain errFormat
	_, _ = fmt.Fprintf(state, fmt.FormatString(state, verb), (*plain)(e))
}
