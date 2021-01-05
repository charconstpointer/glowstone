package main

import (
	"log"
	"net"
	"time"

	"github.com/charconstpointer/glowstone/pkg/glowstone"
	"google.golang.org/protobuf/proto"
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
		log.Println("new client")

		// go io.Copy(client, downstream)
		// go io.Copy(downstream, client)
		go readDs(downstream, client)
		go readUs(downstream, client)
	}(c)
	time.Sleep(123443 * time.Second)
}

func readUs(downstream net.Conn, upstream net.Conn) {
	log.Println("tunnel read upstream")
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
			log.Println("readDs", n)
			payload := b[:n]
			tick := glowstone.Tick{
				Src:     "string",
				Dest:    "string",
				Payload: payload,
			}
			msg, err := proto.Marshal(&tick)
			if err != nil {
				log.Println("marshall", err.Error())
				time.Sleep(250 * time.Millisecond)
			}
			nw, err := downstream.Write(msg)
			if nw < len(msg) {
				log.Println("partial write?", nw, len(msg))
				time.Sleep(250 * time.Millisecond)
			}
			if err != nil {
				log.Println(nw, err.Error())
				time.Sleep(250 * time.Millisecond)
			}
		}
		if err != nil {
			continue
			log.Println("tunnel read upstream", err.Error())
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func readDs(downstream net.Conn, upstream net.Conn) {
	for {
		b := make([]byte, 10000*1024)
		n, err := downstream.Read(b)
		if n == len(b) {
			time.Sleep(time.Second)
			log.Println("hol uppp")

		}
		if n > 0 {
			log.Println("tunnel read downstream", n)

			var tick glowstone.Tick
			err := proto.Unmarshal(b[:n], &tick)
			if err != nil {
				log.Println("unmarshal error", err.Error())
				time.Sleep(250 * time.Millisecond)
			}
			nw, err := upstream.Write(tick.GetPayload())
			log.Println("tunnel sent", nw, n)
			if err != nil {
				log.Println("tunnel read downstream error", nw, err.Error())
				time.Sleep(250 * time.Millisecond)
			}

		}
		if err != nil {
			continue
			log.Println(err.Error())
			time.Sleep(250 * time.Millisecond)
		}
	}
}
