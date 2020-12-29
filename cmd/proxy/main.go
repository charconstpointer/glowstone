package main

import (
	"flag"
	"log"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
)

var (
	downstream = flag.String("downstream", ":25565", "downstream path")
	addr       = flag.String("addr", ":4013", "address for proxy to listen on")
)

func main() {
	flag.Parse()
	log.Println(*downstream, *addr)
	proxy := glowstone.NewProxy(*addr, *downstream)
	proxy.Listen()
}
