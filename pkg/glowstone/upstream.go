package glowstone

import (
	"errors"
	"log"
	"net"

	"google.golang.org/protobuf/proto"
)

type Upstream struct {
	conn       net.Conn
	downstream net.Conn
}

func NewUpstream(conn net.Conn, downstream net.Conn) *Upstream {
	upstream := Upstream{
		conn:       conn,
		downstream: downstream,
	}
	//listen for incoming requests and propagate them to downstream
	go upstream.down()
	return &upstream
}

func (u *Upstream) Up(payload []byte) error {
	n, err := u.conn.Write(payload)
	if err != nil {
		return err
	}
	if n == 0 {
		log.Printf("something went wrong, %d bytes received", n)
		return errors.New("read invalid amount of bytes from downstream")
	}
	return nil
}

func (u *Upstream) down() {
	buffer := make([]byte, 32*1024)
	for {
		n, err := u.conn.Read(buffer)
		if err != nil {
			log.Println(err.Error())
		}
		if n > 0 {
			tick := Tick{
				Src:     u.conn.RemoteAddr().String(),
				Dest:    "mc",
				Payload: buffer[:n],
			}
			b, err := proto.Marshal(&tick)
			if err != nil {
				log.Printf("cannot create marshall tick, %s", err.Error())
			}
			n, err = u.downstream.Write(b)
			if err != nil {
				log.Println(err.Error())
				return
			}
			log.Printf("wrote %d bytes down", n)

		}
	}
}
