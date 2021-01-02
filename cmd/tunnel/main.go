package main

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
)

type Player struct {
	ID         string
	upstream   net.Conn
	downstream net.Conn
}

func NewPlayer(upstream net.Conn, downstream net.Conn, id string) *Player {
	p := Player{
		ID:         id,
		upstream:   upstream,
		downstream: downstream,
	}
	go p.handle()
	return &p
}

func (p *Player) propagateUpstream() {
	for {
		var msg glowstone.Msg
		g := gob.NewDecoder(p.downstream)
		err := g.Decode(&msg)
		if err != nil {
			log.Println("err", err.Error())
		}
		log.Println(msg)
		nw, err := p.upstream.Write(msg.Payload)
		if err != nil {
			log.Println("func (p *Player) propagateUpstream() {", err.Error())
		}
		log.Println("nw", nw)
	}
}
func (p *Player) handle() {
	go p.propagateUpstream()
	size := 32 * 1024
	buf := make([]byte, size)
	for {
		nr, er := p.upstream.Read(buf)
		if er != nil {
			log.Println("could not read from upstream")
		}
		if nr > 0 {
			log.Println("received from mc ", nr)
			g := gob.NewEncoder(p.downstream)
			err := g.Encode(glowstone.NewMsg(buf[:nr], "src string", "dest string"))
			if err != nil {
				log.Println("func (p *Player) handle() {", err.Error())
			}
		}
	}

}
func (p *Player) Handle(data []byte) {
	g := gob.NewEncoder(p.downstream)
	g.Encode(glowstone.NewMsg(data, "src string", "dest string"))
}

func main() {
	//downstream
	ds, err := net.Listen("tcp", ":8889")
	if err != nil {
		log.Fatal(err.Error())
	}

	conn, err := ds.Accept()
	log.Println("tunnel connected")
	if err != nil {
		log.Fatal(err.Error())
	}

	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err.Error())
	}
	players := make([]*Player, 0)
	for {
		client, err := l.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		p := NewPlayer(client, conn, client.RemoteAddr().String())
		players = append(players, p)

	}
}
