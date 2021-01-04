package main

import (
	"flag"
	"log"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
)

var (
	tunnel = flag.String("tunnel", ":8889", "tunnel address")
	mc     = flag.String("mc", ":25565", "minecraft server downstream address")
)

func main() {
	flag.Parse()
	agent, err := glowstone.NewAgent(*tunnel, *mc)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = agent.Listen()
	if err != nil {
		log.Fatal(err.Error())
	}
}
