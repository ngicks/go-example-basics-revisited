package main

import (
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/mapper"
	"github.com/ngicks/go-iterator-helper/hiter/mathiter"
)

func main() {
	fmt.Println("Hello world", Foo)
	fmt.Println(hiter.Sum(mapper.Sprintf("%x", hiter.Limit(8, mathiter.Rng(256)))))
}
