package main

import (
	"github.com/charconstpointer/glowstone/pkg/glowstone"
)

func main() {
	s := glowstone.NewServer()
	s.Listen(":9999")
}
