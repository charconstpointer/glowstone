package main

import (
	"flag"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
)

var (
	downstream = flag.String("downstream", "http://google.com", "downstream path")
	laddr      = flag.String("laddr", "localhost:4013", "address for proxy to listen on")
)

func main() {
	proxy := glowstone.NewProxy(":4013")
	proxy.Listen()
}
