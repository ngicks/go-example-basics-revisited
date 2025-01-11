package main

import "fmt"

type MyError struct {
	Msg string
}

func (e *MyError) Error() string {
	return e.Msg
}

func myTask() (someResult string, err *MyError) {
	return "ok", nil
}

func someTask() (string, error) {
	ret, err := myTask()
	return ret, err
}

func main() {
	ret, err := someTask()
	if err == nil {
		fmt.Println("success")
	} else {
		fmt.Println("failed") // failed
	}
	fmt.Printf("ret = %s, err = %#v\n", ret, err) // ret = ok, err = (*main.MyError)(nil)
}
