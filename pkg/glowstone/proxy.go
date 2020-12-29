package glowstone

import (
	"fmt"
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

func NewProxy(laddr string, ds ...string) TcpProxy {
	return &Proxy{
		laddr:   laddr,
		clients: 0,
		pool:    ds,
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
func (p *Proxy) choseServer() string {
	i := p.clients % len(p.pool)
	log.Println("chosen server", p.pool[i])
	return p.pool[i]
}
func (p *Proxy) handleClient(c net.Conn) {
	server := p.choseServer()
	log.Println("handle new client", c.RemoteAddr(), "connecting to server", server)
	ds, err := net.Dial("tcp", server)
	log.Println(ds.LocalAddr())
	if err != nil {
		log.Fatal(err.Error())

	}
	log.Println("connected to downstream server", ds.RemoteAddr())
	go p.handleClientStream(ds, c)
	go p.handleServerStream(c, ds)

}
func (p *Proxy) decClients() {
	defer p.m.Unlock()
	p.m.Lock()
	p.clients--
}

func (p *Proxy) handleServerStream(c io.Writer, ds io.Reader) {
	if _, err := io.Copy(c, ds); err != nil {
		log.Println("stream closed")
		p.decClients()
		fmt.Println(err)
	}
}

func (p *Proxy) handleClientStream(ds io.Writer, c io.Reader) {
	if _, err := io.Copy(ds, c); err != nil {
		log.Println("stream closed")
		p.decClients()
		fmt.Println(err)
	}
}
