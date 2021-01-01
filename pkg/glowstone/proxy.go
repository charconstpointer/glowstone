package glowstone

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
	downstreams := make([]string, len(ds))
	for _, d := range ds {
		log.Println(d)
		if strings.Contains(d, "localhost") {
			downstreams = append(downstreams, d)
			continue
		}
		if net.ParseIP(d) == nil {
			ip, err := net.LookupIP(d)
			if err != nil {
				log.Printf("could not complete DNS look up for %s", d)
				continue
			}
			for _, ipaddr := range ip {
				downstreams = append(downstreams, ipaddr.String())
			}
		}
	}
	return &Proxy{
		laddr:   laddr,
		clients: 0,
		pool:    downstreams,
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
	go p.copy(ds, c)
	go p.copy(c, ds)
	// go p.copy(c, ds)

}
func (p *Proxy) decClients() {
	defer p.m.Unlock()
	p.m.Lock()
	p.clients--
}

func (p *Proxy) copy(c io.Writer, ds io.Reader) {
	size := 32 * 1024
	buf := make([]byte, size)
	for {
		nr, er := ds.Read(buf)
		msg := "hello world"
		l := len(msg)
		buf = append(buf, []byte(msg)...)
		buf = buf[:len(buf)-l]
		if er != nil {
			log.Println("er != nil ")
		}
		if nr > 0 {
			log.Println(nr)
			nw, ew := c.Write(buf[0:nr])
			if nw > 0 {
				log.Println(nw)
			}
			if ew != nil {
				log.Println("ew != nil")
				break
			}
			if nr != nw {
				log.Println("nr != nw ")
				break
			}
		}
		if er != nil {
			log.Println(er.Error())
			break
		}
	}
}

func (p *Proxy) handleServerStream(c io.Writer, ds io.Reader) {
	buffer := make([]byte, 4096)
	for {
		n, err := ds.Read(buffer)
		log.Println(n)
		if err != nil {
			p.decClients()
			log.Println(err.Error())
		}
	}
}

func (p *Proxy) handleClientStream(ds io.Writer, c io.Reader) {
	if _, err := io.Copy(ds, c); err != nil {
		log.Println("stream closed")
		p.decClients()
		fmt.Println(err)
	}
}
