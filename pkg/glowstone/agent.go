package glowstone

import (
	"log"
	"net"

	"github.com/common-nighthawk/go-figure"
	"google.golang.org/protobuf/proto"
)

type Agent struct {
	upAddr      string
	downAddr    string
	downstreams map[string]net.Conn
	upstream    net.Conn
}

func NewAgent(addr string, downstream string) (*Agent, error) {
	agent := Agent{
		upAddr:   addr,
		downAddr: downstream,
	}
	figure.NewColorFigure("agent", "slant", "purple", true).Print()

	return &agent, nil
}

func (a *Agent) dialUp() error {
	conn, err := net.Dial("tcp", a.upAddr)
	if err != nil {
		return err
	}
	log.Println(conn == nil)
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
	err := a.dialUp()
	if err != nil {
		return err
	}
	err = a.listenDown()
	if err != nil {
		return err
	}
	return nil
}
