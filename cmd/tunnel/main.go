package main

import (
	"context"
	"flag"
	"github.com/charconstpointer/glowstone/pkg/glowstone"
	"golang.org/x/sync/errgroup"
	"log"
)

var (
	up   = flag.String("up", ":8888", "upstream port")
	down = flag.String("down", ":8889", "downstream port")
)

func main() {
	flag.Parse()
	server := glowstone.ListenRpc(*down,*up)
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return server.ListenUp()
	})

	err := g.Wait()
	if err != nil {
		log.Fatal(err.Error())
	}
}
