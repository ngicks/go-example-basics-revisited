package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"syscall"
)

func main() {
	f, err := os.Create("path/to/file")
	if err != nil {
		switch {
		case errors.Is(err, fs.ErrNotExist):
			// errors.Is(err, syscall.ENOENT)と同じ効果
			fmt.Println(fs.ErrNotExist.Error()) // file does not exist
		case errors.Is(err, fs.ErrPermission):
			// errors.Is(err, syscall.EPERM)と同じ効果
			fmt.Println(fs.ErrPermission.Error())
		case errors.Is(err, syscall.EROFS):
		}
	}
	_ = f.Close()
}
