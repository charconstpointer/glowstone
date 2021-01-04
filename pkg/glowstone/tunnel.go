package glowstone

import (
	"context"
	"log"
	"net"

	"github.com/common-nighthawk/go-figure"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

type Tunnel struct {
	addr       string
	downAddr   string
	conn       net.Conn
	downstream net.Conn
	upstreams  map[string]*Upstream
	newUp      chan *Upstream
}

func NewTunnel(addr string, downAddr string) *Tunnel {
	tunnel := Tunnel{
		addr:      addr,
		downAddr:  downAddr,
		newUp:     make(chan *Upstream),
		upstreams: make(map[string]*Upstream),
	}
	figure.NewColorFigure("tunnel", "slant", "green", true).Print()

	return &tunnel
}

func (t *Tunnel) listenUp() error {
	conn, err := net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}
	log.Printf("listening for incoming client connections on %s", conn.Addr().String())

	for {
		client, err := conn.Accept()
		log.Printf("🦚 new upstream connected %s", client.RemoteAddr().String())
		if err != nil {
			log.Println(err.Error())
			continue
		}

		upstream := NewUpstream(client, t.downstream)
		t.upstreams[upstream.conn.RemoteAddr().String()] = upstream
	}
}

func (t *Tunnel) listenDown() error {
	conn, err := net.Listen("tcp", t.downAddr)
	if err != nil {
		return err
	}
	log.Printf("listening for incoming agent connections on %s", conn.Addr().String())

	for {
		client, err := conn.Accept()
		log.Printf("🦚 new agent connected %s", client.RemoteAddr().String())
		if err != nil {
			log.Println(err.Error())
			continue
		}
		t.downstream = client
		go t.readDown()
	}

}

func (t *Tunnel) readUp() error {
	for {
		select {
		case u := <-t.newUp:
			log.Printf("new upstream connected %s", u.conn.RemoteAddr().String())
			// go func(u *Upstream) {
			// 	buffer := make([]byte, 32*1024)

			// 	for {
			// 		n, err := t.downstream.Read(buffer)
			// 		if err != nil {
			// 			log.Println(err.Error())
			// 		}
			// 		if n > 0 {
			// 			tick := Tick{
			// 				Src:     u.conn.RemoteAddr().String(),
			// 				Dest:    "mc",
			// 				Payload: buffer[:n],
			// 			}
			// 			b, err = proto.Marshal(&tick)

			// 			if err != nil {
			// 				log.Println(err.Error())
			// 			}

			// 			err := u.down(tick)

			// 			if err != nil {
			// 				log.Println(err.Error())
			// 			}
			// 		}
			// 	}
			// }(u)
		}
	}
}

func (t *Tunnel) readDown() error {
	buffer := make([]byte, 32*1024)

	for {
		n, err := t.downstream.Read(buffer)
		if err != nil {
			return err
		}
		if n > 0 {

			var tick Tick
			err = proto.Unmarshal(buffer[:n], &tick)

			if err != nil {
				return err
			}
			log.Printf("tick : %v", tick)
			upstream := t.upstreams[tick.Dest]
			err := upstream.Up(tick.Payload)

			if err != nil {
				return err
			}
		}
	}
}

func (t *Tunnel) Listen() error {
	g, _ := errgroup.WithContext(context.Background())
	g.Go(t.listenDown)
	g.Go(t.listenUp)
	// g.Go(t.readDown)
	return g.Wait()
}
