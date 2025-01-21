package main

import (
	"fmt"

	"github.com/ngicks/go-common/serr"
)

//go:noinline
func example() {
	deep()
}

func deep() {
	calling()
}

func calling() {
	frames()
}

func frames() {
	var (
		panicVal any
		done     = make(chan struct{})
	)
	go func() {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			panicVal = serr.WithStack(fmt.Errorf("panicked: %v", rec))
			close(done)
		}()
		example2()
	}()
	<-done
	panic(panicVal)
}

//go:noinline
func example2() {
	deep2()
}

func deep2() {
	calling2()
}

func calling2() {
	frames2()
}

func frames2() {
	s := make([]int, 2)
	_ = s[4]
}

func main() {
	defer func() {
		rec := recover()
		if rec == nil {
			return
		}
		// skip runtime.Callers, inner func, WithStackOpt, gopanic, this func.
		err := serr.WithStackOpt(rec.(error), &serr.WrapStackOpt{Override: true, Skip: 3})
		fmt.Printf("panicked: %v\n", rec)
		var i int
		for seq := range serr.DeepFrames(err) {
			if i > 0 {
				fmt.Printf("caused by\n")
			}
			i++
			for f := range seq {
				fmt.Printf("    %s(%s:%d)\n", f.Function, f.File, f.Line)
			}
		}
	}()
	example()
	//nolint
	// panicked: panicked: runtime error: index out of range [4] with length 2
	//     main.main.func1(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:67)
	//     runtime.gopanic(runtime/panic.go:785)
	//     main.frames(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:39)
	//     main.calling(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:19)
	//     main.deep(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:15)
	//     main.example(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:11)
	//     main.main(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:80)
	//     runtime.main(runtime/proc.go:272)
	// caused by
	//     main.frames.func1.1(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:33)
	//     runtime.gopanic(runtime/panic.go:785)
	//     runtime.goPanicIndex(runtime/panic.go:115)
	//     main.frames2(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:57)
	//     main.calling2(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:52)
	//     main.deep2(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:48)
	//     main.example2(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:44)
	//     main.frames.func1(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:36)
}
