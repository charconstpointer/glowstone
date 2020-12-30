package main

import (
	"log"
	"net"
	"regexp"
	"time"
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
		time.Sleep(1000 * time.Millisecond)
		nr, er := conn.Read(buf)
		if er != nil {
			log.Println("could not read from upstream")
		}
		log.Println(nr)
		if nr > 0 {
			id := "marysia"
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
		// if er != nil {
		// 	log.Println(er.Error())
		// 	// h.downstream.Close()
		// 	continue
		// }
		if nr > 0 {
			// id := make([]byte, 64)
			// copy(id, h.ID)
			// toWrite := append(buf[0:nr], id...)
			// toWrite := append(buf[0:nr], id...)
			nw, _ := h.upstream.Write(buf[:nr])
			if nw > 0 {

			}
			// if ew != nil {
			// 	log.Println("ew != nil")
			// 	break
			// }
			if nr != nw {
				log.Println("propagateUpstream nr != nw ")
				break
			}
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
	log.Println(string(data))
	if len(data) > 0 {
		nw, _ := h.downstream.Write(data)
		log.Printf("handler %s wrote %d bytes", h.ID, nw)
		if nw > 0 {
			// log.Println(nw)
		}
		// if ew != nil {
		// 	log.Println("ew != nil")
		// }
		// if len(data) != nw {
		// 	log.Println("nr != nw ")
		// }
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
