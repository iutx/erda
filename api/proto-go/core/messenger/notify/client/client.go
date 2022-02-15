// Code generated by protoc-gen-go-client. DO NOT EDIT.
// Sources: notify.proto

package client

import (
	context "context"

	grpc "github.com/erda-project/erda-infra/pkg/transport/grpc"
	pb "github.com/erda-project/erda-proto-go/core/messenger/notify/pb"
	grpc1 "google.golang.org/grpc"
)

// Client provide all service clients.
type Client interface {
	// NotifyService notify.proto
	NotifyService() pb.NotifyServiceClient
}

// New create client
func New(cc grpc.ClientConnInterface) Client {
	return &serviceClients{
		notifyService: pb.NewNotifyServiceClient(cc),
	}
}

type serviceClients struct {
	notifyService pb.NotifyServiceClient
}

func (c *serviceClients) NotifyService() pb.NotifyServiceClient {
	return c.notifyService
}

type notifyServiceWrapper struct {
	client pb.NotifyServiceClient
	opts   []grpc1.CallOption
}

func (s *notifyServiceWrapper) CreateNotifyHistory(ctx context.Context, req *pb.CreateNotifyHistoryRequest) (*pb.CreateNotifyHistoryResponse, error) {
	return s.client.CreateNotifyHistory(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}

func (s *notifyServiceWrapper) QueryNotifyHistories(ctx context.Context, req *pb.QueryNotifyHistoriesRequest) (*pb.QueryNotifyHistoriesResponse, error) {
	return s.client.QueryNotifyHistories(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}

func (s *notifyServiceWrapper) GetNotifyStatus(ctx context.Context, req *pb.GetNotifyStatusRequest) (*pb.GetNotifyStatusResponse, error) {
	return s.client.GetNotifyStatus(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}

func (s *notifyServiceWrapper) GetNotifyHistogram(ctx context.Context, req *pb.GetNotifyHistogramRequest) (*pb.GetNotifyHistogramResponse, error) {
	return s.client.GetNotifyHistogram(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}
