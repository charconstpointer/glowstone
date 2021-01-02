package main

import (
	"encoding/json"
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
	size := 32 * 1024
	buf := make([]byte, size)
	for {
		nr, er := p.downstream.Read(buf)
		if nr > 0 {
			//TODO remove client id
			var msg glowstone.Msg
			json.Unmarshal(buf[:nr], &msg)
			nw, _ := p.upstream.Write(msg.Payload)
			// nw, _ := p.upstream.Write(buf[:nr])
			if nw > 0 {

			}
			// if nr != nw {
			// 	log.Println("propagateUpstream nr != nw ")
			// 	break
			// }
		}
		if er != nil {
			log.Println("propagateUpstream.err")
			log.Println(er.Error())
			p.downstream.Close()
			break
		}
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
		log.Println(nr)
		if nr > 0 {
			msg := glowstone.NewMsg(buf[:nr], "src string", "dest string")
			b, err := json.Marshal(msg)
			if err != nil {
				log.Fatal(err.Error())
			}
			//TODO append client id
			p.downstream.Write(b)
		}
	}

}
func (p *Player) Handle(data []byte) {
	log.Println(string(data))
	if len(data) > 0 {
		msg := glowstone.NewMsg(data, "src string", "dest string")
		b, err := json.Marshal(msg)
		if err != nil {
			log.Fatal(err.Error())
		}
		p.downstream.Write(b)
	}
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
		go func(client net.Conn, p *Player) {
			log.Println(p.ID, "connected")
			if err != nil {
				log.Fatal(err.Error())

			}
			// go copyClient(conn, client)
			// go copyServer(client, conn)
		}(client, p)
	}
}
