//go:build refs

package ref

import (
	"backend-masterclass/controllers/grpc"
	"backend-masterclass/controllers/protoc"
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

type repo = grpc.Repo
type Gateway struct {
	BankService *grpc.Service
}

func NewGWService(r repo) (s *Gateway) {
	s = &Gateway{
		BankService: grpc.NewService(r),
	}

	return
}

func (gw *Gateway) Start(address string) (err error) {

	// -----------------Json configurations-------------------

	// Setting json options so that variable names would be identical to the proto file and request fields that are not in the protobuf would be discarded.
	jsonOptions := runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},

			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	)

	// ------Mapping http requests to gRPC calls--------------

	// Creating an object that maps urls to grpc calls.
	grpcMux := runtime.NewServeMux(jsonOptions)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Registering our bank service methods as handler functions in our grpcMux mapping.
	protoc.RegisterBankServiceHandlerServer(ctx, grpcMux,
		gw.BankService)
	if err != nil {
		return
	}

	// Creating an object that uses our grpcMux mapping in http requests.
	httpMux := http.NewServeMux()
	httpMux.Handle("/", grpcMux)

	// ------Listening on the given address (url:port)--------

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	log.Printf("Listening on %s\n", listener.Addr().String())

	http.Serve(listener, httpMux)
	return
}
