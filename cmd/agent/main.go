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
	conn, _ := net.Dial("tcp", ":8889")
	downstream, _ := net.Dial("tpc", ":25565")
	go func(upstream net.Conn) {
		b := make([]byte, 32*1024)
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

				n, err := downstream.Write(tick.Payload)
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
	b := make([]byte, 32*1024)
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
			b, err := proto.Marshal(&tick)
			if err != nil {
				log.Fatal(err.Error())
			}

			n, err := upstream.Write(b)
			if err != nil {
				log.Fatal(err.Error())
			}
			log.Printf("wrote %d bytes upstream %s", n, upstream.RemoteAddr().String())
		}
	}
}
