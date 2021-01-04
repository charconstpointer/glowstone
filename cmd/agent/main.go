package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
	"github.com/golang/protobuf/proto"
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

	go func(upstream net.Conn) {
		conn, err := net.Dial("tcp", *mc)
		if err != nil {
			log.Println(err.Error())
		}
		log.Println("connected to mc")
		go readDs(conn, upstream)
		go readUs(conn, upstream)
	}(upstream)
	time.Sleep(100340 * time.Second)
}

func readDs(downstream net.Conn, upstream net.Conn) {
	log.Println("readDs", downstream.RemoteAddr().String())
	b := make([]byte, 1024*32)
	for {
		n, err := downstream.Read(b)
		if err != nil {
			log.Println(err.Error())
		}
		if n > 0 {
			log.Println("readDs", n)
			payload := b[:n]
			tick := glowstone.Tick{
				Src:     "string",
				Dest:    "string",
				Payload: payload,
			}
			msg, _ := proto.Marshal(&tick)
			nw, err := upstream.Write(msg)
			if err != nil {
				log.Println(nw, err.Error())
			}
			log.Println("readDs", nw)
		}
	}
}

func readUs(downstream net.Conn, upstream net.Conn) {
	log.Println("readUs", downstream.RemoteAddr().String())

	b := make([]byte, 1024*32)
	for {
		n, err := upstream.Read(b)
		if err != nil {
			log.Println(err.Error())
		}
		if n > 0 {
			log.Println("readUs", n)

			var tick glowstone.Tick
			proto.Unmarshal(b[:n], &tick)
			nw, err := downstream.Write(tick.Payload)
			log.Println("readUs", nw)
			if err != nil {
				log.Println(nw, err.Error())
			}
		}
	}
}
