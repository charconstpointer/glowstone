package glowstone

import (
	"encoding/gob"
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
			g := gob.NewEncoder(h.upstream)
			err := g.Encode(NewMsg(buf[:nr], "src string", "dest string"))
			if err != nil {
				log.Println("func (h *Handler) propagateUpstream() {", err.Error())
			}
			log.Println(nr)
		}
		if er != nil {
			log.Println("propagateUpstream.err")
			log.Println(er.Error())
			h.downstream.Close()
			break
		}
	}
}

func (h *Handler) Handle(msg *Msg) {
	nw, err := h.downstream.Write(msg.Payload)
	if nw > 0 {
		log.Println("nw", nw)
	}

	if err != nil {
		log.Println("func (h *Handler) Handle(data []byte) {", err.Error())
	}
}

func NewHandler(id string, upstream net.Conn) *Handler {
	server := ":25565"
	ds, err := net.Dial("tcp", server)
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
