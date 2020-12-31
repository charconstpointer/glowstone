package main

import (
	"io"
	"log"
	"net"
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
			id := buf[:len(buf)-64]
			log.Println("ID?", string(id))
			//TODO remove client id
			nw, _ := p.upstream.Write(buf[:nr])
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
			//TODO append client id
			p.downstream.Write(buf[:nr])
		}
	}

}
func (p *Player) Handle(data []byte) {
	log.Println(string(data))
	if len(data) > 0 {
		nw, _ := p.downstream.Write(data)
		log.Printf("handler %s wrote %d bytes", p.ID, nw)
		if nw > 0 {
			// log.Println(nw)
		}
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

func copyClient(c io.Writer, ds net.Conn) {
	size := 32 * 1024
	buf := make([]byte, size)
	for {
		nr, er := ds.Read(buf)
		if er != nil {
			log.Println("er != nil ")
			ds.Close()
		}
		if nr > 0 {
			nw, ew := c.Write(buf[:nr])
			if nw > 0 {
				log.Println(nw)
			}
			if ew != nil {
				log.Println("ew != nil")
				break
			}
			if nr != nw {
				log.Println("nr != nw ", nr, nw)
				break
			}
		}
		if er != nil {
			log.Println(er.Error())
			ds.Close()
			break
		}
	}
}

func copyServer(c io.Writer, ds net.Conn) {
	size := 32 * 1024
	buf := make([]byte, size)
	for {
		nr, er := ds.Read(buf)
		if er != nil {
			log.Println("er != nil ")
			ds.Close()
		}
		if nr > 0 {
			nw, ew := c.Write(buf[:nr])
			if nw > 0 {
				log.Println("sent up", nw)
			}
			if ew != nil {
				log.Println("ew != nil")
				break
			}
			if nr != nw {
				// if nr != nw {
				log.Println("nr != nw ")
				break
			}
		}
		if er != nil {
			log.Println(er.Error())
			ds.Close()
			break
		}
	}
}
