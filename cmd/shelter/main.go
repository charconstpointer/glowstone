package main

import (
	"encoding/json"
	"log"
	"net"
	"regexp"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
)

func main() {
	handlers := make([]*Handler, 0)
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
				handler := NewHandler(id, conn)
				handlers = append(handlers, handler)
			}
			handler := handlersContains(handlers, id)
			handler.Handle(buf[:nr])
		}
	}
}

type Handler struct {
	ID         string
	upstream   net.Conn
	downstream net.Conn
}

func (h *Handler) propagateUpstream() {
	size := 32 * 1024
	buf := make([]byte, size)
	for {
		nr, er := h.downstream.Read(buf)
		if nr > 0 {
			msg := glowstone.NewMsg(buf[:nr], "src string", "dest string")
			b, _ := json.Marshal(&msg)
			nw, _ := h.upstream.Write(b)
			log.Println(nw)
		}
		if er != nil {
			log.Println("propagateUpstream.err")
			log.Println(er.Error())
			h.downstream.Close()
			break
		}
	}
}

func (h *Handler) Handle(data []byte) {
	if len(data) > 0 {
		var msg glowstone.Msg
		json.Unmarshal(data, &msg)
		nw, _ := h.downstream.Write(msg.Payload)
		//TODO remove client id
		log.Printf("handler %s wrote %d bytes", h.ID, nw)
		if nw > 0 {
			// log.Println(nw)
		}
	}
}

func NewHandler(id string, upstream net.Conn) *Handler {
	server := ":25565"
	ds, err := net.Dial("tcp", server)
	log.Println(ds.LocalAddr())
	if err != nil {
		log.Fatal(err.Error())

	}

	handler := Handler{
		ID:         id,
		upstream:   upstream,
		downstream: ds,
	}
	go handler.propagateUpstream()
	return &handler
}

func handlersContains(handlers []*Handler, id string) *Handler {
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
