package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrInvalidArgs = errors.New("not enough arguments, a minimum of 2 required")

func main() {
	if len(os.Args) < 2 {
		fmt.Println(ErrInvalidArgs)
		os.Exit(1)
	}

	envPath := os.Args[1]
	cmdArgs := os.Args[2:]

	env, err := ReadDir(envPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	os.Exit(RunCmd(cmdArgs, env))
}
