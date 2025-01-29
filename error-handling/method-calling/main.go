package main

import "fmt"

type foo int

//go:noinline
func (f *foo) Bar(baz int) {
	fmt.Println(f, baz)
}

func main() {
	f := foo(0x55)
	f.Bar(0x123)
}
