package main

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"runtime"
)

const maxDepth = 100

type withStack struct {
	err error
	pc  []uintptr
}

func (e *withStack) Error() string {
	return e.err.Error()
}

func (e *withStack) Unwrap() error {
	return e.err
}

func WithStack(err error) error {
	var pc [maxDepth]uintptr
	// skip runtime.Callers, WithStack
	n := runtime.Callers(2, pc[:])
	return &withStack{
		err: err,
		pc:  pc[:n],
	}
}

func Frames(err error) iter.Seq[runtime.Frame] {
	return func(yield func(runtime.Frame) bool) {
		var ws *withStack
		if !errors.As(err, &ws) {
			return
		}

		frames := runtime.CallersFrames(ws.pc)
		for {
			f, ok := frames.Next()
			if !ok {
				return
			}
			if !yield(f) {
				return
			}
		}
	}
}

func PrintStack(w io.Writer, err error) error {
	for f := range Frames(err) {
		_, err := fmt.Fprintf(w, "%s(%s:%d)\n", f.Function, f.File, f.Line)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	sample := errors.New("sample")

	wrapped := example(sample)

	fmt.Printf("%v\n", wrapped) // sample
	err := PrintStack(os.Stdout, wrapped)
	if err != nil {
		panic(err)
	}
	/*
		main.frames(github.com/ngicks/go-example-basics-revisited/error-handling/with-stack/main.go:91)
		main.calling(github.com/ngicks/go-example-basics-revisited/error-handling/with-stack/main.go:87)
		main.deep(github.com/ngicks/go-example-basics-revisited/error-handling/with-stack/main.go:83)
		main.example(github.com/ngicks/go-example-basics-revisited/error-handling/with-stack/main.go:79)
		main.main(github.com/ngicks/go-example-basics-revisited/error-handling/with-stack/main.go:70)
		runtime.main(runtime/proc.go:272)
	*/
}

func example(err error) error {
	return deep(err)
}

func deep(err error) error {
	return calling(err)
}

func calling(err error) error {
	return frames(err)
}

func frames(err error) error {
	return WithStack(err)
}
