package main

import (
	"flag"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
)

var (
	downstream = flag.String("downstream", "144.217.199.16:25613", "downstream path")
	addr       = flag.String("addr", ":4013", "address for proxy to listen on")
)

func main() {
	flag.Parse()
	proxy := glowstone.NewProxy(*addr, *downstream)
	proxy.Listen()
}
