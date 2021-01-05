package main

import (
	"github.com/charconstpointer/glowstone/pkg/glowstone"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"time"
)

func main() {
	c, _ := net.Listen("tcp", ":8889")

	go func(c net.Listener) {
		downstream, _ := c.Accept()
		log.Println("downstream connected")

		clients, err := net.Listen("tcp", ":9000")
		if err != nil {
			log.Println(err.Error())
		}
		client, _ := clients.Accept()

		go readDs(downstream, client)
		go readUs(downstream, client)
	}(c)
	time.Sleep(123443 * time.Second)
}

func readUs(downstream net.Conn, upstream net.Conn) {
	log.Println("tunnel read upstream")
	b := make([]byte, 512)
	for {
		n, err := upstream.Read(b)
		log.Println("read from up", n)
		if n > 0 {
			tick := glowstone.Tick{
				Src:     "string",
				Dest:    "stringxxxxxxxx",
				Payload: b[:n],
			}
			msg, err := proto.Marshal(&tick)
			if err != nil {
				log.Println("marshall", err.Error())
				continue
			}
			nw, err := downstream.Write(msg)
			//nw, err := downstream.Write(b[:n])

			if err != nil {
				log.Println(nw, err.Error())
				time.Sleep(250 * time.Millisecond)
			}
		}
		if err != nil {

			log.Println("tunnel read upstream", err.Error())
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func readDs(downstream net.Conn, upstream net.Conn) {
	log.Println("agent read upstream")
	b := make([]byte, 1000)
	for {
		n, err := downstream.Read(b)
		if n > 0 {
			var tick glowstone.Tick
			err := proto.Unmarshal(b[:n], &tick)
			if err != nil {
				log.Println(b[:n])
				log.Println("marshall", err.Error())
				time.Sleep(250 * time.Millisecond)
			}
			log.Println(tick.Src)
			log.Println(tick.Dest)
			log.Println(len(tick.Payload))
			nw, err := upstream.Write(tick.GetPayload())
			//nw, err := downstream.Write(b[:n])
			if err != nil {
				log.Println(b[:n])
				log.Println(nw, err.Error())
				time.Sleep(250 * time.Millisecond)
			}
		}
		if err != nil {
			log.Println(err.Error())
			time.Sleep(250 * time.Millisecond)
		}
	}
}
