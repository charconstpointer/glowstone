package main

import (
	"flag"
	"log"
	"net"

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
	b := make([]byte, 32*1024)
	go func() {
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
				msg, err := proto.Marshal(&tick)
				if err != nil {
					log.Fatal(err.Error())
				}
				n, err = ds.Write(msg)
				log.Println(n)
			}
		}
	}()

	for {

		n, err := us.Read(b)

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

			n, err := ds.Write(tick.Payload)
			if err != nil {
				log.Fatal(err.Error())
			}
			log.Printf("wrote %d bytes to minecraft", n)
		}
	}
}
