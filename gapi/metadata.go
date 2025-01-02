package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGateWayUserAgentKey = "grpcgateway-user-agent"
	// The IP of the client, the key x-forwarded-host is the gateway's IP.
	grpcGateWayClientIPKey = "x-forwarded-for"

	grpcUserAgentKey = "user-agent"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (md *Metadata) GetFromRaw(ctx context.Context, rawMD metadata.MD) {
	// -----------------In case of a gateway.-----------------------------
	if userAgents :=
		rawMD.Get(grpcGateWayUserAgentKey); len(userAgents) > 0 {
		md.UserAgent = userAgents[0]
	}

	if clientIPs :=
		rawMD.Get(grpcGateWayClientIPKey); len(clientIPs) > 0 {
		md.ClientIP = clientIPs[0]
	}

	// -----------------In case of pure gRPC.-----------------------------
	if userAgents :=
		rawMD.Get(grpcUserAgentKey); len(userAgents) > 0 {
		md.UserAgent = userAgents[0]
	}

	if peerInfo, exists := peer.FromContext(ctx); exists {
		md.ClientIP = peerInfo.Addr.String()
	}
}

func (server *Server) extractMetadata(ctx context.Context,
) (md *Metadata) {
	md = &Metadata{}

	rawMD, exists := metadata.FromIncomingContext(ctx)
	if !exists {
		return
	}

	md.GetFromRaw(ctx, rawMD)

	return
}
