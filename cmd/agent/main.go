package main

import (
	"flag"
	"log"
	"net"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
	"github.com/golang/protobuf/proto"
)

var (
	tunnel = flag.String("tunnel", ":8889", "tunnel address")
	mc     = flag.String("mc", ":25565", "minecraft server downstream address")
)

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", ":8889")
	if err != nil {
		log.Fatal(err.Error())
	}
	downstream, err := net.Dial("tcp", ":25565")
	if err != nil {
		log.Fatal(err.Error())
	}
	go func(upstream net.Conn) {
		b := make([]byte, 1024*100)
		for {
			n, err := upstream.Read(b)

			if err != nil {
				log.Fatal(err.Error())
				continue
			}
			if n > 0 {
				var tick glowstone.Tick
				err := proto.Unmarshal(b[:n], &tick)
				if err != nil {
					log.Fatal(err.Error())
				}
				sent := 0
				for sent < len(tick.Payload) {
					bs, _ := downstream.Write(tick.Payload[sent:])
					sent += bs
				}
				log.Println(n == len(b))
				if err != nil {
					log.Fatal(err.Error())
				}
				log.Printf("wrote %d bytes to minecraft", n)
			}
		}
	}(conn)

	start(conn.RemoteAddr().String(), downstream, conn)

}
func start(dest string, downstream net.Conn, upstream net.Conn) {
	b := make([]byte, 1024*100)
	for {
		n, err := downstream.Read(b)

		if err != nil {
			log.Fatal(err.Error())
			continue
		}
		if n > 0 {
			tick := glowstone.Tick{
				Src:     "mc",
				Dest:    dest,
				Payload: b[:n],
			}
			// log.Printf("tick %v", tick)
			msg, err := proto.Marshal(&tick)
			if err != nil {
				log.Fatal(err.Error())
			}
			sent := 0
			for sent < len(msg) {
				bs, _ := upstream.Write(msg[sent:])
				sent += bs
			}
			if err != nil {
				log.Fatal(err.Error())
			}
			log.Println(n == len(b))
			if n != len(msg) {
				log.Println(sent, len(msg))
			}
			log.Printf("wrote %d bytes upstream %s", n, upstream.RemoteAddr().String())
		}
	}
}
