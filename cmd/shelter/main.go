package main

import (
	"encoding/gob"
	"log"
	"net"
	"regexp"

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
	for {
		var msg glowstone.Msg
		g := gob.NewDecoder(conn)
		g.Decode(&msg)
		if err != nil {
			log.Println("func (h *Handler) Handle(data []byte) {", err.Error())
		}
		if msg.Payload == nil {
			log.Println("if msg.Payload == nil {")
			continue
		}
		if handlersContains(handlers, "id") == nil {
			log.Printf("Creating new handler for %s", "id")
			handler := glowstone.NewHandler("id", conn)
			handlers = append(handlers, handler)
		}
		handler := handlersContains(handlers, "id")
		handler.Handle(&msg)
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
