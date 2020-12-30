package main

import (
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err.Error())
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}

		go func(c net.Conn) {
			log.Println("handle", c.RemoteAddr())
		}(conn)
	}
}
