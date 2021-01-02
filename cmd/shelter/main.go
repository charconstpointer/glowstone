package main

import (
	"context"
	"github.com/charconstpointer/glowstone/pkg/glowstone"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	conn, err := grpc.Dial(":8889", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	client := glowstone.NewGlowClient(conn)
	stream, err := client.Listen(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	downstream, err := net.Dial("tcp", ":25565")
	if err != nil {
		log.Fatal(err.Error())
	}

	//propagate up
	go func(stream glowstone.Glow_ListenClient) {
		b := make([]byte, 32*1024)
		for {
			n, err := downstream.Read(b)
			if err != nil {
				log.Fatal(err.Error())
			}
			err = stream.Send(&glowstone.Tick{
				Src:     "downstream",
				Dest:    "tunnel",
				Payload: b[:n],
			})
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}(stream)
	//propagate down
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Println(err.Error())
			time.Sleep(100 * time.Millisecond)
			continue
		}
		n, err := downstream.Write(msg.Payload)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println("wrote", n, "bytes to downstream")
	}

}
