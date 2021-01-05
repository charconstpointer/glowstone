package main

import (
	"flag"
	"log"
	"net"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
	"google.golang.org/protobuf/proto"
)

var (
	tunnel = flag.String("tunnel", ":8889", "tunnel address")
	mc     = flag.String("mc", ":25565", "minecraft server downstream address")
)

func main() {
	flag.Parse()
	upstream, err := net.Dial("tcp", *tunnel)

	if err != nil {
		log.Println(err.Error())
	}
	conn, err := net.Dial("tcp", *mc)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("connected to mc")
	go readDs(conn, upstream)
	readUs(conn, upstream)
}

func readDs(downstream net.Conn, upstream net.Conn) {
	b := make([]byte, 4096)
	for {
		n, err := downstream.Read(b)
		if n > 0 {
			tick := glowstone.Tick{
				Payload: b[:n],
			}
			msg, err := proto.Marshal(&tick)
			if err != nil {
				log.Println("marshall", err.Error())

			}
			nw, err := upstream.Write(msg)
			// nw, err := upstream.Write(b[:n])

			if err != nil {
				log.Println(nw, err.Error())

			}
		}
		if err != nil {
			log.Println("tunnel read upstream", err.Error())
		}
	}
}

func readUs(downstream net.Conn, upstream net.Conn) {
	b := make([]byte, 4096)
	for {
		n, err := upstream.Read(b)
		if n > 0 {
			var tick glowstone.Tick
			err := proto.Unmarshal(b[:n], &tick)
			if err != nil {
				log.Println("marshall", err.Error())
			}
			nw, err := downstream.Write(tick.GetPayload())
			// nw, err := downstream.Write(b[:n])
			if err != nil {
				log.Println(nw, err.Error())
			}
		}
		if err != nil {
			log.Println(err.Error())
		}
	}
}
