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
	up   = flag.String("up", ":8888", "upstream port")
	down = flag.String("down", ":8889", "downstream port")
)

func main() {
	flag.Parse()
	dc, err := net.Listen("tcp", *down)
	if err != nil {
		log.Fatal(err.Error())
	}
	ds, err := dc.Accept()
	if err != nil {
		log.Fatal(err.Error())
	}

	uc, err := net.Listen("tcp", *up)
	us, err := uc.Accept()
	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		b := make([]byte, 1024*100)
		for {

			n, err := us.Read(b)
			if err != nil {
				log.Fatal(err.Error())
			}
			if n > 0 {
				tick := glowstone.Tick{
					Src:     us.RemoteAddr().String(),
					Dest:    "mc",
					Payload: b[:n],
				}
				// log.Printf("tick %v", tick)

				msg, err := proto.Marshal(&tick)
				if n != len(msg) {
					log.Println(n, len(msg))
				}
				if err != nil {
					log.Fatal(err.Error())
				}
				sent := 0
				for sent < len(msg) {
					bs, _ := ds.Write(msg[sent:])
					sent += bs
				}

				log.Println(n == len(b))
			}
		}
	}()

	go func(upstream net.Conn, downstream net.Conn) {
		b := make([]byte, 1024*100)
		for {
			n, err := downstream.Read(b)

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
					bs, err := upstream.Write(tick.Payload[sent:])
					if err != nil {
						log.Fatal(err.Error())
					}
					sent += bs
				}
				log.Printf("wrote %d bytes to client", sent)
			}
		}
	}(us, ds)
	time.Sleep(321323 * time.Second)
}
