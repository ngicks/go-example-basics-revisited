package main

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter/stringsiter"
	goplayground "github.com/ngicks/go-playground"
)

func main() {
	birds := []string{"🐔", "🐣", "🐧", "🐓"}
	rand.Shuffle(len(birds), func(i, j int) { birds[i], birds[j] = birds[j], birds[i] })
	fmt.Println(goplayground.Yay())
	fmt.Printf(
		"🐤< ｺﾝﾆﾁﾊ！ ₍₍⁽⁽ %v ₎₎⁾⁾\n ",
		stringsiter.Join(
			"₎₎⁾⁾ ₍₍⁽⁽",
			slices.Values(birds),
		),
	)
}
