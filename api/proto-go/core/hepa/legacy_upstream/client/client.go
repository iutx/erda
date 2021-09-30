// Code generated by protoc-gen-go-client. DO NOT EDIT.
// Sources: legacy_upstream.proto

package client

import (
	context "context"
	grpc "github.com/erda-project/erda-infra/pkg/transport/grpc"
	pb "github.com/erda-project/erda-proto-go/core/hepa/legacy_upstream/pb"
	grpc1 "google.golang.org/grpc"
)

// Client provide all service clients.
type Client interface {
	// UpstreamService legacy_upstream.proto
	UpstreamService() pb.UpstreamServiceClient
}

// New create client
func New(cc grpc.ClientConnInterface) Client {
	return &serviceClients{
		upstreamService: pb.NewUpstreamServiceClient(cc),
	}
}

type serviceClients struct {
	upstreamService pb.UpstreamServiceClient
}

func (c *serviceClients) UpstreamService() pb.UpstreamServiceClient {
	return c.upstreamService
}

type upstreamServiceWrapper struct {
	client pb.UpstreamServiceClient
	opts   []grpc1.CallOption
}

func (s *upstreamServiceWrapper) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return s.client.Register(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}

func (s *upstreamServiceWrapper) AsyncRegister(ctx context.Context, req *pb.AsyncRegisterRequest) (*pb.AsyncRegisterResponse, error) {
	return s.client.AsyncRegister(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}
