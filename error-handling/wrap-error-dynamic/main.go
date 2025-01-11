package main

import (
	"errors"
	"fmt"
	"strings"
)

func main() {
	var (
		err1 = errors.New("1")
		err2 = errors.New("2")
		err3 = errors.New("3")
	)

	errs := []any{err1, err2, err3}

	const sep = ", "
	format, _ := strings.CutSuffix(strings.Repeat("%w"+sep, len(errs)), sep)
	wrapped := fmt.Errorf(format, errs...)

	fmt.Printf("err = %v\n", wrapped) // err = 1, 2, 3
}
