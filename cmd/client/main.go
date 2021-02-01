package main

import (
	"log"
	"time"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
)

func main() {
	m := glowstone.NewMux()
	if err := m.Dial("ec2-3-8-238-16.eu-west-2.compute.amazonaws.com:8000"); err != nil {
		// if err := m.Dial(":8000"); err != nil {
		// if err := m.Dial("159.89.4.159:8000"); err != nil {
		log.Fatal(err.Error())
	}

	if err := m.Recv(); err != nil {
		log.Fatal(err.Error())
	}

	time.Sleep(time.Hour)
}
