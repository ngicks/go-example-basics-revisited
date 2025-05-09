package pkg2

import (
	"fmt"

	"github.com/ngicks/go-example-basics-revisited/starting-projects/internal/i0"
	"github.com/ngicks/go-example-basics-revisited/starting-projects/pkg1"
	"github.com/ngicks/go-example-basics-revisited/starting-projects/pkg2/internal/i2"
)

func SayDouble() string {
	return fmt.Sprintf("%q%q", pkg1.Foo, pkg1.Foo)
}

func SayYay0() string {
	return i0.Yay0
}

func SayYay2() string {
	return i2.Yay2
}
