package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	original := "Hello, OTUS!"

	fmt.Println(reverse.String(original))
}
