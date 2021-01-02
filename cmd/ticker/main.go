package main

import (
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		n, _ := conn.Write([]byte("streaming?"))
		log.Println("wrote ", n, "bytes")
		time.Sleep(time.Second)
	}
}
