package glowstone

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type GlowstoneServer interface {
	Listen(addr string) error
	CreateServer() error
}

type HTTPServer struct {
	servers []*Server
}

type Server struct {
	ID        string `json:"id"`
	CreatedAt int    `json:"createdAt"`
}

func NewServer() *HTTPServer {
	return &HTTPServer{
		servers: make([]*Server, 0),
	}
}

func (s *HTTPServer) addServers(serv *Server) error {
	for _, sv := range s.servers {
		if sv.ID == serv.ID {
			return fmt.Errorf("server %s already exists on glowstone server", serv.ID)
		}
	}
	s.servers = append(s.servers, serv)
	return nil
}

func (s *HTTPServer) CreateServer() error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "itzg/minecraft-server",
		Env:   []string{"EULA=TRUE"},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"25565/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "25565",
				},
			},
		},
	}, nil, nil, "")
	if err != nil {
		return err
	}
	serv := &Server{
		ID:        resp.ID,
		CreatedAt: int(time.Now().Unix()),
	}
	s.addServers(serv)
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return err
}

func (s *HTTPServer) Listen(addr string) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/servers", s.handleCreateServer)
	r.Get("/servers", s.handleGetServers)
	err := http.ListenAndServe(addr, r)
	return err
}

func (s *HTTPServer) handleCreateServer(rw http.ResponseWriter, r *http.Request) {
	go s.CreateServer()
	rw.WriteHeader(http.StatusAccepted)
}

func (s *HTTPServer) handleGetServers(rw http.ResponseWriter, r *http.Request) {
	servers := s.servers
	if len(servers) == 0 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusOK)
	b, err := json.Marshal(servers)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(b)
}
