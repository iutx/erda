// Code generated by protoc-gen-go-client. DO NOT EDIT.
// Sources: settings.proto

package client

import (
	context "context"

	grpc "github.com/erda-project/erda-infra/pkg/transport/grpc"
	pb "github.com/erda-project/erda-proto-go/core/monitor/settings/pb"
	grpc1 "google.golang.org/grpc"
)

// Client provide all service clients.
type Client interface {
	// SettingsService settings.proto
	SettingsService() pb.SettingsServiceClient
}

// New create client
func New(cc grpc.ClientConnInterface) Client {
	return &serviceClients{
		settingsService: pb.NewSettingsServiceClient(cc),
	}
}

type serviceClients struct {
	settingsService pb.SettingsServiceClient
}

func (c *serviceClients) SettingsService() pb.SettingsServiceClient {
	return c.settingsService
}

type settingsServiceWrapper struct {
	client pb.SettingsServiceClient
	opts   []grpc1.CallOption
}

func (s *settingsServiceWrapper) GetSettings(ctx context.Context, req *pb.GetSettingsRequest) (*pb.GetSettingsResponse, error) {
	return s.client.GetSettings(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}

func (s *settingsServiceWrapper) PutSettings(ctx context.Context, req *pb.PutSettingsRequest) (*pb.PutSettingsResponse, error) {
	return s.client.PutSettings(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}

func (s *settingsServiceWrapper) RegisterMonitorConfig(ctx context.Context, req *pb.RegisterMonitorConfigRequest) (*pb.RegisterMonitorConfigResponse, error) {
	return s.client.RegisterMonitorConfig(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}
