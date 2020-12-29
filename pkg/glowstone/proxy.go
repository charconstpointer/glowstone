package glowstone

import (
	"io"
	"log"
	"net"
	"sync"
)

type TcpProxy interface {
	Listen() error
}
type Proxy struct {
	clients int
	pool    []string
	laddr   string
	m       sync.Mutex
}

func NewProxy(laddr string) TcpProxy {
	return &Proxy{
		laddr:   laddr,
		clients: 0,
		pool: []string{
			":25565",
			":25566",
		},
	}
}

func (p *Proxy) Listen() error {
	l, err := net.Listen("tcp", ":4013")
	if err != nil {
		log.Println(err.Error())
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err.Error())
		}
		p.m.Lock()
		p.clients++
		p.m.Unlock()
		go p.handleClient(c)
	}
}

func (p *Proxy) handleClient(c net.Conn) {
	log.Println("handle new client", c.RemoteAddr())
	ds, err := net.Dial("tcp", ":25565")
	if err != nil {
		log.Fatal(err.Error())

	}
	log.Println("connected to downstream server", ds)

	go io.Copy(c, ds)
	go io.Copy(ds, c)
}
