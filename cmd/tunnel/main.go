package main

import (
	"log"
	"net"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
	"google.golang.org/protobuf/proto"
)

func main() {
	c, _ := net.Listen("tcp", ":8889")
	downstream, _ := c.Accept()

	clients, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Println(err.Error())
	}
	client, _ := clients.Accept()

	go readDs(downstream, client)
	readUs(downstream, client)

}

func readUs(downstream net.Conn, upstream net.Conn) {
	b := make([]byte, 4096)
	for {
		n, err := upstream.Read(b)
		if n > 0 {
			tick := glowstone.Tick{
				Payload: b[:n],
			}
			msg, err := proto.Marshal(&tick)
			if err != nil {
				log.Println("marshall", err.Error())
				continue
			}
			nw, err := downstream.Write(msg)
			// nw, err := downstream.Write(b[:n])

			if err != nil {
				log.Println(nw, err.Error())
			}
		}
		if err != nil {
			log.Println("tunnel read upstream", err.Error())
		}
	}
}

func readDs(downstream net.Conn, upstream net.Conn) {
	b := make([]byte, 4096)
	for {
		n, err := downstream.Read(b)
		if n > 0 {
			var tick glowstone.Tick
			err := proto.Unmarshal(b[:n], &tick)
			if err != nil {
				log.Println("unmarshall", err.Error())
			}
			nw, err := upstream.Write(tick.Payload)
			// nw, err := upstream.Write(b[:n])
			if err != nil {
				log.Println(nw, err.Error())
			}
		}
		if err != nil {
			log.Println(err.Error())
		}
	}
}
