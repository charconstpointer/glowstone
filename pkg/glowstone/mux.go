package glowstone

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Mux struct {
	clients map[int]net.Conn
	next    net.Conn
}

func NewMux() *Mux {
	return &Mux{
		clients: make(map[int]net.Conn),
	}
}

func (m *Mux) Dial(addr string) error {
	if m.next != nil {
		return fmt.Errorf("cannot dial new mux as this mux is already connected with a different one on addr %s", m.next.RemoteAddr().String())
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "cannot dial another mutex")
	}

	m.next = conn
	return err
}

func (m *Mux) ListenMux(addr string) error {
	l, err := net.Listen("tcp", addr)
	conn, err := l.Accept()
	m.next = conn
	go m.Recv()
	log.Printf("new mux connected %s", conn.RemoteAddr().String())
	return err
}

func (m *Mux) Listen(addr string) error {
	g, _ := errgroup.WithContext(context.Background())
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		id := getID(conn.RemoteAddr().String())
		m.clients[int(id)] = conn
		g.Go(func() error {
			return m.handleConn(conn, int(id))
		})
	}

}

func (m *Mux) Recv() error {
	for {

		h := Header(make([]byte, HeaderSize))
		n, err := io.ReadFull(m.next, h)
		if h.Len() == 0 {
			m.next.Close()
			return errors.New("dc?")
		}
		id := h.ID()
		if err != nil {
			log.Fatal(err.Error())
		}

		b := make([]byte, int(h.Len()))
		n, err = io.ReadFull(m.next, b)
		if err != nil {
			return err
		}

		if n > 0 {
			c, e := m.clients[int(id)]
			if !e {
				// conn, err := net.Dial("tcp", ":25565")
				conn, err := net.Dial("tcp", ":64884")
				if err != nil {
					log.Println("cannot dial minecraft", err.Error())
					return err
				}
				m.clients[int(id)] = conn
				go m.handleConn(conn, int(id))
				time.Sleep(time.Millisecond * 100)
			}
			c = m.clients[int(id)]
			io.Copy(c, bytes.NewBuffer(b))
			if err != nil {
				log.Println("cannot write payload to minecraft")
			}

			if n == 0 {
				log.Println("cannot write payload to minecraft")
			}
		}
	}
}

func (m *Mux) handleConn(conn net.Conn, id int) error {
	log.Printf("handling new conn %s", conn.RemoteAddr().String())
	for {
		b := make([]byte, 1024)
		nr, err := conn.Read(b)

		if err != nil {
			return errors.Wrap(err, "cannot read from handled connection")
		}
		if nr == 0 {
			conn.Close()
			return errors.New("dc??")
		}
		if nr > 0 {
			h := make(Header, HeaderSize)
			h.Encode(PASS, int32(id), int32(nr))
			sent := 0
			for sent < HeaderSize {
				n, _ := m.next.Write(h)
				sent += n
			}

			_, _ = io.Copy(m.next, bytes.NewReader(b[:nr]))

		}
	}
}

func getID(addr string) uint32 {
	i := strings.LastIndex(addr, ":")
	id := addr[i+1:]
	parsed, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err.Error())
	}

	return uint32(parsed)
}
