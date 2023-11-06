package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		panic("Wrong arguments count: min 2 args")
	}

	dir := os.Args[1]
	command := os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		log.Fatalf("Err: %s. Can't read env", err)
	}

	os.Exit(RunCmd(command, env))
}
