package main

import (
	"fmt"
	"runtime"

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
		pc := make([]uintptr, 100)
		// skip runtime.Callers, this closure, runtime.gopanic
		n := runtime.Callers(3, pc)
		pc = pc[:n]

		fmt.Printf("panicked: %v\n", rec)
		frames := runtime.CallersFrames(pc)
		for {
			f, ok := frames.Next()
			if !ok {
				break
			}
			fmt.Printf("    %s(%s:%d)\n", f.Function, f.File, f.Line)
		}
		fmt.Printf("caused by\n")
		for f := range serr.Frames(rec.(error)) {
			fmt.Printf("    %s(%s:%d)\n", f.Function, f.File, f.Line)
		}
	}()
	example()
	//nolint
	// panicked: panicked: yay
	//     main.frames(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:39)
	//     main.calling(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:19)
	//     main.deep(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:15)
	//     main.example(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:11)
	//     main.main(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:80)
	//     runtime.main(runtime/proc.go:272)
	// caused by
	//     main.frames.func1.1(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:33)
	//     runtime.gopanic(runtime/panic.go:785)
	//     main.frames2(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:55)
	//     main.calling2(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:51)
	//     main.deep2(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:47)
	//     main.example2(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:43)
	//     main.frames.func1(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:36)
}
