package glowstone

import (
	"context"
	"log"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

type Agent struct {
	upAddr      string
	downAddr    string
	downstreams map[string]net.Conn
	upstream    net.Conn
}

func NewAgent(addr string, downstream string) *Agent {
	agent := Agent{
		upAddr:   addr,
		downAddr: downstream,
	}
	return &agent
}

func (a *Agent) dialUp() error {
	conn, err := net.Dial("tcp", a.upAddr)
	if err != nil {
		return err
	}

	a.upstream = conn
	return nil
}
func (a *Agent) createDownstream() (net.Conn, error) {
	conn, err := net.Dial("tcp", a.downAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (a *Agent) listenDown() error {
	buffer := make([]byte, 32*1024)

	for {
		n, err := a.upstream.Read(buffer)
		if err != nil {
			return err
		}
		if n > 0 {
			var tick Tick
			err = proto.Unmarshal(buffer[:n], &tick)

			if err != nil {
				return err
			}

			if a.downstreams[tick.Src] == nil {
				downstream, err := a.createDownstream()
				if err != nil {
					return err
				}
				a.downstreams[tick.Src] = downstream
			}

			downstream := a.downstreams[tick.Dest]
			n, err := downstream.Write(tick.Payload)
			if err != nil {
				return err
			}

			log.Printf("🥕wrote %d bytes down", n)
		}
	}
}

func (a *Agent) Listen() error {
	g, _ := errgroup.WithContext(context.Background())
	g.Go(a.dialUp)
	g.Go(a.listenDown)
	return g.Wait()
}
