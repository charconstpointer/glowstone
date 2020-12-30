package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":8889")
	if err != nil {
		log.Fatal(err.Error())
	}
	server := ":25565"
	ds, err := net.Dial("tcp", server)
	log.Println(ds.LocalAddr())
	if err != nil {
		log.Fatal(err.Error())

	}
	log.Println("connected to downstream server", ds.RemoteAddr())
	go handleClientStream(ds, conn)
	go handleServerStream(conn, ds)
	time.Sleep(time.Second * 999999)
}
func ignore() {
	// l, err := net.Listen("tcp", ":8888")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// for {
	// 	conn, err := l.Accept()
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}

	// 	go func(c net.Conn) {
	// 		server := ":25565"
	// 		log.Println("handle new client", c.RemoteAddr(), "connecting to server", server)
	// 		ds, err := net.Dial("tcp", server)
	// 		log.Println(ds.LocalAddr())
	// 		if err != nil {
	// 			log.Fatal(err.Error())

	// 		}
	// 		log.Println("connected to downstream server", ds.RemoteAddr())
	// 		go handleClientStream(ds, c)
	// 		go handleServerStream(c, ds)

	// 	}(conn)
	// }
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
