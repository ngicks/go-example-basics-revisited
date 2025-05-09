package pkg1

import (
	"github.com/ngicks/go-example-basics-revisited/starting-projects/internal/i0"
	"github.com/ngicks/go-example-basics-revisited/starting-projects/pkg1/internal/i1"
)

var Foo = "foo"

func SayYay0() string {
	return i0.Yay0
}

func SayYay1() string {
	return i1.Yay1
}
