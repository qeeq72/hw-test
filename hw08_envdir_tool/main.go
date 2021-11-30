package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Too few arguments to run envdir tool!")
	} else {
		env, err := ReadDir(args[0])
		if err != nil {
			os.Exit(1)
			return
		}
		os.Exit(RunCmd(args[1:], env))
	}
}
