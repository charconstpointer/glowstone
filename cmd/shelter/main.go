package main

import "net"

func main() {
	conn, _ := net.Listen("tcp", ":8889")
	conn.Accept()
}
