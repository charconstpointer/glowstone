package main

import (
	"flag"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
)

var (
	addr = flag.String("addr", ":9000", "address to listen on for http requests")
)

func main() {
	flag.Parse()
	s := glowstone.NewServer()
	s.Listen(*addr)
}
