package main

import (
	"fmt"
)

type uncomparableErr1 []error

func (e uncomparableErr1) Error() string {
	return "uncomparableErr1"
}

type uncomparableErr2 struct {
	errs []error
}

func (e uncomparableErr2) Error() string {
	return "uncomparableErr2"
}

func main() {
	ue1 := uncomparableErr1{fmt.Errorf("foo"), fmt.Errorf("bar")}
	ue2 := uncomparableErr2{errs: []error{fmt.Errorf("foo"), fmt.Errorf("bar")}}

	// invalid operation: ue1 == ue1 (slice can only be compared to nil) compiler(UndefinedOp)
	// if ue1 == ue1 {
	// }

	compareErr(ue1, ue2)
	// comparing err1 panicked = runtime error: comparing uncomparable type main.uncomparableErr1
	// comparing err2 panicked = runtime error: comparing uncomparable type main.uncomparableErr2

	compareErr(&ue1, &ue2)
	// err1 equal = uncomparableErr1
	// err2 equal = uncomparableErr2
}

func compareErr(err1, err2 error) {
	if err1 == err2 {
		fmt.Println("huh?")
	}

	func() {
		defer func() {
			if rec := recover(); rec != nil {
				fmt.Printf("comparing err1 panicked = %v\n", rec)
			}
		}()
		if err1 == err1 {
			fmt.Printf("err1 equal = %v\n", err1)
		}
	}()

	func() {
		defer func() {
			if rec := recover(); rec != nil {
				fmt.Printf("comparing err2 panicked = %v\n", rec)
			}
		}()
		if err2 == err2 {
			fmt.Printf("err2 equal = %v\n", err2)
		}
	}()
}
