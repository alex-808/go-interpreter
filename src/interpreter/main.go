package main

import (
	"fmt"
	"github.com/alex-davis-808/go-interpreter/src/interpreter/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s. Welcome to Chlorophyll.\n", user.Username)
	fmt.Printf("Type commands here\n")
	repl.Start(os.Stdin, os.Stdout)
}
