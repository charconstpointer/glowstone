package main

import (
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
			go copyClient(conn, client)
			go copyServer(client, conn)
		}(client)
	}
}

func copyClient(c io.Writer, ds net.Conn) {
	size := 32 * 1024
	buf := make([]byte, size)
	for {
		nr, er := ds.Read(buf)
		if er != nil {
			log.Println("er != nil ")
			ds.Close()
		}
		if nr > 0 {
			log.Println("nr", nr)
			log.Println(string(buf))
			// id := make([]byte, 64)
			// copy(id, "marysia")
			// copy(id, ds.RemoteAddr().String())
			// toWrite := append(buf[0:nr], id...)
			// log.Println("TOWRITE", string(toWrite))
			// time.Sleep(time.Second)
			nw, ew := c.Write(buf[:nr])
			// nw, ew := c.Write(buf[0:nr])
			if nw > 0 {
				log.Println(nw)
			}
			if ew != nil {
				log.Println("ew != nil")
				break
			}
			if nr != nw {
				log.Println("nr != nw ", nr, nw)
				break
			}
		}
		if er != nil {
			log.Println(er.Error())
			ds.Close()
			break
		}
	}
}

func copyServer(c io.Writer, ds net.Conn) {
	size := 32 * 1024
	buf := make([]byte, size)
	for {
		nr, er := ds.Read(buf)
		if er != nil {
			log.Println("er != nil ")
			ds.Close()
		}
		if nr > 0 {
			// log.Println("B", nr)
			// id := string(buf[nr-64:])
			// log.Println("ID", id)
			// toWrite := buf[0 : nr-64]
			// toWrite := buf
			nw, ew := c.Write(buf[:nr])
			// nw, ew := c.Write(buf[0:nr])
			if nw > 0 {
				log.Println("sent up", nw)
			}
			if ew != nil {
				log.Println("ew != nil")
				break
			}
			if nr != nw {
				// if nr != nw {
				log.Println("nr != nw ")
				break
			}
		}
		if er != nil {
			log.Println(er.Error())
			ds.Close()
			break
		}
	}
}
