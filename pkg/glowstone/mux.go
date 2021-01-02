package glowstone

import (
	"log"
	"net"
)

type RpcServer struct {
	upstream net.Conn
}

func (s *RpcServer) Listen(stream Glow_ListenServer) error {
	go func(stream Glow_ListenServer) {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Println("cannot recv msg", err.Error())
				continue
			}
			n, err := s.upstream.Write(msg.Payload)
			if err != nil {
				log.Println("could not push to upstream", err.Error())
			}
			log.Printf("wrote %d bytes to upstream", n)

		}
	}(stream)
	b := make([]byte, 32*1024)
	for {
		n, err := s.upstream.Read(b)
		if err != nil {
			log.Println(err.Error())
		}
		msg := Tick{
			Payload: b[:n],
			Src:     "Src",
			Dest:    "Dest",
		}
		err = stream.Send(&msg)
		if err != nil {
			log.Println("could not send msg", err.Error())
		}
	}

}

func (s *RpcServer) mustEmbedUnimplementedGlowServer() error {
	return nil
}
