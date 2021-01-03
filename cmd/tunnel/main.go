package main

import (
	"flag"
	"log"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
)

var (
	up   = flag.String("up", ":8888", "upstream port")
	down = flag.String("down", ":8889", "downstream port")
)

func main() {
	flag.Parse()
	tunnel := glowstone.NewTunnel(*up, *down)
	err := tunnel.Listen()
	if err != nil {
		log.Fatal(err.Error())
	}
}
