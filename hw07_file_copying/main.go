package main

import (
	"errors"
	"flag"
	"log"
)

var (
	from, to      string
	limit, offset int64
)

var (
	ErrorFromFlagRequired = errors.New("file to read from is not defined")
	ErrorToFlagRequired   = errors.New("file to write to is not defined")
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	if from == "" {
		log.Fatal(ErrorFromFlagRequired)
	}
	if to == "" {
		log.Fatal(ErrorToFlagRequired)
	}

	if err := Copy(from, to, offset, limit); err != nil {
		log.Fatal(err)
	}
}
