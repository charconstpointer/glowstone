package main

import (
	"log"
	"net"
	"regexp"
	"time"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
)

func main() {
	handlers := make([]*glowstone.Handler, 0)
	conn, err := net.Dial("tcp", ":8889")
	if err != nil {
		log.Fatal(err.Error())
	}

	if err != nil {
		log.Fatal(err.Error())

	}
	size := 32 * 1024
	buf := make([]byte, size)
	for {
		time.Sleep(1000 * time.Millisecond)
		nr, er := conn.Read(buf)
		if er != nil {
			log.Println("could not read from upstream")
		}
		log.Println(nr)
		if nr > 0 {
			//TODO replace with dynamic id
			id := "someid"
			log.Println("handling", id)
			if handlersContains(handlers, id) == nil {
				log.Printf("Creating new handler for %s", id)
				handler := glowstone.NewHandler(id, conn)
				handlers = append(handlers, handler)
			}
			handler := handlersContains(handlers, id)
			handler.Handle(buf[:nr])
		}
	}
}

func handlersContains(handlers []*glowstone.Handler, id string) *glowstone.Handler {
	r, _ := regexp.Compile("^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$")
	id = r.FindString(id)
	for _, c := range handlers {
		toMatch := r.FindString(c.ID)

		if toMatch == id {
			return c
		}
	}
	return nil
}
func contains(clients []string, id string) bool {
	r, _ := regexp.Compile("^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$")
	id = r.FindString(id)
	for _, c := range clients {
		c = r.FindString(c)

		if c == id {
			return true
		}
	}
	return false
}
