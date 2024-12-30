package gapi

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/token"
	u "backend-masterclass/util"
)

// Serves all gRPC requests for our banking service.

type Server struct {
	Store sqlc.Store
	// Responsible for creating the context for each route.
	// It will automatically send the context to the handler functions.
	TokenMaker token.Maker
	Config     u.Config
}

func NewServer(store sqlc.Store, config u.Config, maker token.Maker,
) (s *Server) {
	s = &Server{
		Store:      store,
		Config:     config,
		TokenMaker: maker,
	}

	return
}

func (server *Server) Start(address string) (err error) {
	return
}
