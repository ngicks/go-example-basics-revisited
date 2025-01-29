package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/ngicks/go-common/serr"
)

//go:noinline
func example(ctx context.Context) {
	deep(ctx)
}

func deep(ctx context.Context) {
	calling(ctx)
}

func calling(ctx context.Context) {
	frames(ctx)
}

func frames(ctx context.Context) {
	var (
		panicVal  any
		panicOnce sync.Once
		wg        sync.WaitGroup
	)
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			panicOnce.Do(func() {
				// In case there's many goroutines running and
				// you are going to capture only first panic value recovered.
				panicVal = serr.WithStack(fmt.Errorf("panicked: %v", rec))
			})
			cancel(panicVal.(error))
		}()
		example2(ctx)
	}()
	wg.Wait()
	if panicVal != nil {
		panic(panicVal)
	}
}

//go:noinline
func example2(ctx context.Context) {
	deep2(ctx)
}

func deep2(ctx context.Context) {
	calling2(ctx)
}

func calling2(ctx context.Context) {
	frames2(ctx)
}

func frames2(_ context.Context) {
	s := make([]int, 2)
	_ = s[4]
}

func main() {
	defer func() {
		rec := recover()
		if rec == nil {
			return
		}
		// skip runtime.Callers, inner func, WithStackOpt.
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
	example(context.Background())
	//nolint
	// panicked: panicked: runtime error: index out of range [4] with length 2
	//     main.main.func1(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:80)
	//     runtime.gopanic(runtime/panic.go:785)
	//     main.frames(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:51)
	//     main.calling(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:21)
	//     main.deep(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:17)
	//     main.example(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:13)
	//     main.main(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:93)
	//     runtime.main(runtime/proc.go:272)
	// caused by
	//     main.frames.func1.1.1(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:43)
	//     sync.(*Once).doSlow(sync/once.go:76)
	//     sync.(*Once).Do(sync/once.go:67)
	//     main.frames.func1.1(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:40)
	//     runtime.gopanic(runtime/panic.go:785)
	//     runtime.goPanicIndex(runtime/panic.go:115)
	//     main.frames2(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:70)
	//     main.calling2(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:65)
	//     main.deep2(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:61)
	//     main.example2(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:57)
	//     main.frames.func1(github.com/ngicks/go-example-basics-revisited/error-handling/log-stacktrace/main.go:47)
}
