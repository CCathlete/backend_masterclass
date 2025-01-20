package grpc

import (
	"backend-masterclass/controllers/protoc"
	"backend-masterclass/models/entities"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Service struct {
	protoc.UnimplementedBankServiceServer
	Repo entities.Repo
}

func NewService(repo entities.Repo) (service *Service) {
	service = &Service{
		Repo: repo,
	}

	return
}

func (service *Service) Start(address string) (err error) {

	server := grpc.NewServer()

	protoc.RegisterBankServiceServer(server, service)
	reflection.Register(server)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return
	}

	log.Printf("gRPC server started at %s", listener.Addr().String())

	err = server.Serve(listener)
	return
}
