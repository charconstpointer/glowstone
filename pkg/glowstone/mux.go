package glowstone

import (
	context "context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"golang.org/x/sync/errgroup"
)

type RpcServer struct {
	upstreams []net.Conn
	upAddr    string
}

func NewRpcServer(up string) *RpcServer {
	return &RpcServer{
		upAddr: up,
	}
}

func (s *RpcServer) ListenUp() error {
	l, err := net.Listen("tcp", s.upAddr)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		log.Println("new client connected", conn.RemoteAddr())
		if err != nil {
			return err
		}
		s.upstreams = append(s.upstreams, conn)
	}
}

func (s *RpcServer) Listen(stream Glow_ListenServer) error {
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	for _, upstream := range s.upstreams {
		g.Go(func() error {
			return listenUpstream(upstream, stream)
		})
	}
	g.Go(func() error {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Println("cannot receive", err.Error())
				time.Sleep(time.Second)
				continue
			}
			log.Println(s.upstreams)
			n, err := s.upstreams[0].Write(msg.Payload)
			if err != nil {
				log.Println("cant write")
			}
			log.Println("wrote", n, "bytes up")
		}
		return nil
	})
	err := g.Wait()
	return err
}

func listenUpstream(upstream net.Conn, stream Glow_ListenServer) error {
	b := make([]byte, 32*1024)
	for {
		n, err := upstream.Read(b)
		if err != nil {
			return err
		}
		msg := Tick{
			Payload: b[:n],
			Src:     "Src",
			Dest:    "Dest",
		}
		err = stream.Send(&msg)
		if err != nil {
			log.Println("could not send msg", err.Error())
			return err
		}
	}
}

func (s *RpcServer) mustEmbedUnimplementedGlowServer() {
}

func ListenRpc(down string, up string) *RpcServer {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost%s", down))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	server := NewRpcServer(up)
	RegisterGlowServer(grpcServer, server)
	go grpcServer.Serve(lis)
	return server

}
