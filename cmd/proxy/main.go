package main

import (
	"log"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
)

func main() {
	m := glowstone.NewMux()
	if err := m.ListenMux(":8000"); err != nil {
		log.Println(err.Error())
	}

	if err := m.Listen(":9000"); err != nil {
		log.Println(err.Error())
	}
}
