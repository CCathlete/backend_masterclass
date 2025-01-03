package gapi

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/rpc"
	"backend-masterclass/token"
	u "backend-masterclass/util"
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "backend-masterclass/doc/statik" // Binary static files.
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
	// Creating a new gRPC server instance.
	grpcServer := grpc.NewServer()

	// Makes our server available for gRPC calls.
	rpc.RegisterSimpleBankServer(grpcServer, server)

	reflection.Register(grpcServer)

	// Listening on the specified TCP network address (e.g., "0.0.0.0:9090").
	// This creates a network listener which the gRPC server will accept connections from.
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return
	}

	log.Println("Starting gRPC server on", listener.Addr().String())

	// Start serving incoming connections on the listener.
	// This call blocks and serves connections until the server is stopped.
	err = grpcServer.Serve(listener)
	return
}

func (server *Server) StartGatewayServer(address string,
) (err error) {

	jsonOption :=
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		})

	// Gets HTTP requests and forwards them to the gRPC server.
	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Registering our gRPC server to get calls from grpcMux.
	err = rpc.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		return
	}

	// Wrapping our grpcMux with an HTTP mux so all HTTP routes are handled by the grpc-gateway and converted to gRPC calls.
	HTTPMux := http.NewServeMux()
	HTTPMux.Handle("/", grpcMux)

	// Creating a handler that serves static files.
	staticFs, err := fs.New()
	if err != nil {
		return
	}

	// Strips the /swagger/ from paths it gets and uses the rest as a filesystem path.
	staticHandler :=
		http.StripPrefix("/swagger/", http.FileServer(staticFs))

	HTTPMux.Handle("/swagger/", staticHandler)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return
	}

	log.Println(
		"Starting HTTP gateway server on", listener.Addr().String())

	err = http.Serve(listener, HTTPMux)
	return
}
