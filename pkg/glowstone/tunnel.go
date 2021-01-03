package glowstone

import (
	"context"
	"log"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

type Tunnel struct {
	addr       string
	downAddr   string
	conn       net.Conn
	downstream net.Conn
	upstreams  map[string]*Upstream
}

func NewTunnel(addr string, downAddr string) *Tunnel {
	tunnel := Tunnel{
		addr:     addr,
		downAddr: downAddr,
	}
	return &tunnel
}

func (t *Tunnel) listenUp() error {
	conn, err := net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}

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

			upstream := t.upstreams[tick.Dest]
			err := upstream.Up(tick.Payload)

			if err != nil {
				return err
			}
		}
	}
}

func (t *Tunnel) dialDown() error {
	conn, err := net.Dial("tcp", t.downAddr)
	if err != nil {
		return err
	}
	t.downstream = conn
	log.Printf("💰 connected to downstream %s", t.downstream.RemoteAddr())
	return nil
}

func (t *Tunnel) Listen() error {
	g, _ := errgroup.WithContext(context.Background())
	g.Go(t.dialDown)
	g.Go(t.listenUp)
	return g.Wait()
}
