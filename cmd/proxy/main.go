package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/pkg/errors"
)

var (
	downstream = flag.String("downstream", "http://google.com", "downstream path")
	laddr = flag.String("laddr", "localhost:4013", "address for proxy to listen on")
)

func main() {
	flag.Parse()
	log.Println(*downstream)
	l, err := net.Listen("tcp", *laddr)
	if err != nil {
		errors.Wrap(err, "could not start proxy")
	}
	for {
		log.Println("waiting for clients")
		s, err := l.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		go func(s net.Conn) {
			log.Println(s.RemoteAddr())
			downStream, err := net.Dial("tpc", "localhost:25565")
			
			log.Println(downStream)
			if err != nil {
				errors.Wrap(err, "could not created downstream TCP client")
				return
			}
			go func() {
				log.Println("cr reader")
				cr := bufio.NewReader(s)
				p := make([]byte, 4)
				for {
					n, err := cr.Read(p)
					if err == io.EOF {
						break
					}
					fmt.Println(string(p[:n]))
				}
			}()

			go func() {
				log.Println("sr reader")
				sr := bufio.NewReader(downStream)
				p := make([]byte, 4)
				for {
					n, err := sr.Read(p)
					if err == io.EOF {
						break
					}
					fmt.Println(string(p[:n]))
				}
			}()

		}(s)
	}
	time.Sleep(1000 * time.Second)
}
