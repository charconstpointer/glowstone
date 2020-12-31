package glowstone

import (
	"log"
	"net"
)

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
			toWrite := buf
			//TODO append client id
			nw, _ := h.upstream.Write(toWrite[:nr])
			if nw > 0 {

			}
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
		//TODO remove client id
		nw, _ := h.downstream.Write(data)
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
