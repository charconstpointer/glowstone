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
	log.Println("readDs", downstream.RemoteAddr().String())
	b := make([]byte, 1024*32)
	for {
		n, err := upstream.Read(b)
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
			nw, err := downstream.Write(msg)
			if err != nil {
				log.Println(nw, err.Error())
			}
			log.Println("readDs", nw)
		}
	}
}

func readDs(downstream net.Conn, upstream net.Conn) {
	log.Println("readUs", downstream.RemoteAddr().String())
	b := make([]byte, 1024*32)
	for {
		n, err := downstream.Read(b)
		if err != nil {
			log.Println(err.Error())
		}
		if n > 0 {
			log.Println("readDs", n)

			var tick glowstone.Tick
			proto.Unmarshal(b[:n], &tick)
			nw, err := upstream.Write(tick.Payload)
			if err != nil {
				log.Println(nw, err.Error())
			}
			log.Println("readDs", nw)
		}
	}
}
