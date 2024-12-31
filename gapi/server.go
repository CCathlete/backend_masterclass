package gapi

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/rpc"
	"backend-masterclass/token"
	u "backend-masterclass/util"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Serves all gRPC requests for our banking service.

type Server struct {
	rpc.UnimplementedSimpleBankServer
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
	// Create a new gRPC server instance.
	grpcServer := grpc.NewServer()

	// Makes our server available for gRPC calls.
	rpc.RegisterSimpleBankServer(grpcServer, server)

	reflection.Register(grpcServer)

	// Listen on the specified TCP network address (e.g., "0.0.0.0:9090").
	// This creates a network listener which the gRPC server will accept connections from.
	listener, err := net.Listen("tcp", address)
	if err != nil {
		// If there's an error in creating the listener, return the error.
		return
	}

	log.Println("Starting gRPC server on", listener.Addr().String())

	// Start serving incoming connections on the listener.
	// This call blocks and serves connections until the server is stopped.
	err = grpcServer.Serve(listener)
	return
}
