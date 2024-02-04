package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset uint64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Uint64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Uint64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if err := Copy(from, to, int64(offset), int64(limit)); err != nil {
		fmt.Println(err)
	}
}
