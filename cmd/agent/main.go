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
	log.Println("agent read downstream", downstream.RemoteAddr().String())
	for {
		b := make([]byte, 10000*1024)
		n, err := downstream.Read(b)
		if n == len(b) {
			time.Sleep(time.Second)
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
		}
		if n > 0 {
			log.Println("readDs", n)
			payload := b[:n]
			tick := &glowstone.Tick{
				Src:     "string",
				Dest:    "string",
				Payload: payload,
			}
			msg, err := proto.Marshal(tick)
			if err != nil {
				log.Println("marshall", err.Error())
				time.Sleep(250 * time.Millisecond)
			}
			nw, err := upstream.Write(msg)
			if nw < len(msg) {
				log.Println("partial write?")
				time.Sleep(250 * time.Millisecond)
			}
			log.Println("agent sent", nw, n)
			if err != nil {
				log.Println(nw, err.Error())
				time.Sleep(250 * time.Millisecond)
			}
			log.Println("readDs", nw)
		}

		if err != nil {
			log.Println(err.Error())
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func readUs(downstream net.Conn, upstream net.Conn) {
	log.Println("agent read upstream")

	for {
		b := make([]byte, 10000*1024)
		n, err := upstream.Read(b)
		if n == len(b) {
			time.Sleep(time.Second)
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
			log.Println("hol uppp")
		}
		if n > 0 {
			log.Println("readUs", n)

			var tick glowstone.Tick
			err := proto.Unmarshal(b[:n], &tick)
			if err != nil {
				log.Println("marshall", err.Error())
				time.Sleep(250 * time.Millisecond)
			}
			if len(tick.GetPayload()) > 2097152 {
				log.Fatal("2097152")
				time.Sleep(250 * time.Millisecond)
			}
			nw, err := downstream.Write(tick.GetPayload())
			log.Println("readUs", nw)
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
