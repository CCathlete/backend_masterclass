//go:build refs

package ref

import (
	"backend-masterclass/controllers/protoc"
	"backend-masterclass/models/entities"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Repo = entities.Repo

type Service struct {
	protoc.UnimplementedBankServiceServer
	Repo entities.Repo
}

func NewService(repo Repo) (s *Service) {
	s = &Service{
		Repo: repo,
	}

	return
}

func (service *Service) Start(address string) (err error) {

	// Creating a gRPC server and registering our bank service.
	server := grpc.NewServer()
	protoc.RegisterBankServiceServer(server, service)

	// Registering a reflection service on our server.
	reflection.Register(server)

	// Listening on the given address (url:port).
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	log.Printf("Starting server on %s\n", listener.Addr().
		String())

	// Starting the server.
	err = server.Serve(listener)

	return
}
