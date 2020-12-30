package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	//downstream
	ds, err := net.Listen("tcp", ":8889")
	if err != nil {
		log.Fatal(err.Error())
	}

	conn, err := ds.Accept()
	log.Println("tunnel connected")
	if err != nil {
		log.Fatal(err.Error())
	}

	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}

		go func(client net.Conn) {
			if err != nil {
				log.Fatal(err.Error())

			}
			go handleClientStream(conn, client)
			go handleServerStream(client, conn)
		}(client)
	}
}

func handleServerStream(c io.Writer, ds io.Reader) {
	if _, err := io.Copy(c, ds); err != nil {
		log.Println("stream closed")
		fmt.Println(err)
	}
}

func handleClientStream(ds io.Writer, c io.Reader) {
	if _, err := io.Copy(ds, c); err != nil {
		log.Println("stream closed")
		fmt.Println(err)
	}
}
