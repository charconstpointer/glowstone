package glowstone

import (
	"context"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server interface {
	Listen(addr string) error
	CreateServer() error
}

type GlowServer struct {
}

func NewServer() *GlowServer {
	return &GlowServer{}
}

func (s *GlowServer) CreateServer() error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "itzg/minecraft-server",
		Env:   []string{"EULA=TRUE"},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	return err
}

func (s *GlowServer) Listen(addr string) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/servers", func(rw http.ResponseWriter, r *http.Request) {
		err := s.CreateServer()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("could not create new server"))
			return
		}
		rw.WriteHeader(http.StatusAccepted)
	})
	err := http.ListenAndServe(addr, r)
	return err
}
