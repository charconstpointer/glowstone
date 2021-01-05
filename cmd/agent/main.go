package main

import (
	"flag"
	"log"
	"net"
	"time"

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
	log.Println("tunnel read upstream")
	b := make([]byte, 512)
	for {
		n, err := downstream.Read(b)
		log.Println("read from up", n)
		if n > 0 {
			tick := glowstone.Tick{
				Src:     "stringxddddd",
				Dest:    "string",
				Payload: b[:n],
			}
			msg, err := proto.Marshal(&tick)
			if err != nil {
				log.Println("marshall", err.Error())
				continue
			}
			nw, err := upstream.Write(msg)
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

func readUs(downstream net.Conn, upstream net.Conn) {
	log.Println("agent read upstream")
	b := make([]byte, 531)
	for {
		n, err := upstream.Read(b)
		if n > 0 {
			var tick glowstone.Tick
			err := proto.Unmarshal(b[:n], &tick)
			if err != nil {
				log.Println("marshall", err.Error())
				time.Sleep(250 * time.Millisecond)
			}
			nw, err := downstream.Write(tick.GetPayload())
			//nw, err := downstream.Write(b[:n])
			if err != nil {
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
