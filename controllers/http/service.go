package http

import (
	"backend-masterclass/controllers/grpc"
	"backend-masterclass/models/entities"
)

type Service struct {
	*grpc.Service
}

func NewService(repo entities.Repo) (s *Service) {
	s = &Service{
		Service: grpc.NewService(repo),
	}

	return
}

func (service *Service) Start(address string) {

}
