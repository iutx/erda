// Code generated by protoc-gen-go-client. DO NOT EDIT.
// Sources: cmp_alert.proto

package client

import (
	context "context"

	grpc "github.com/erda-project/erda-infra/pkg/transport/grpc"
	pb "github.com/erda-project/erda-proto-go/cmp/alert/pb"
	grpc1 "google.golang.org/grpc"
)

// Client provide all service clients.
type Client interface {
	// AlertService cmp_alert.proto
	AlertService() pb.AlertServiceClient
}

// New create client
func New(cc grpc.ClientConnInterface) Client {
	return &serviceClients{
		alertService: pb.NewAlertServiceClient(cc),
	}
}

type serviceClients struct {
	alertService pb.AlertServiceClient
}

func (c *serviceClients) AlertService() pb.AlertServiceClient {
	return c.alertService
}

type alertServiceWrapper struct {
	client pb.AlertServiceClient
	opts   []grpc1.CallOption
}

func (s *alertServiceWrapper) GetAlertConditions(ctx context.Context, req *pb.GetAlertConditionsRequest) (*pb.GetAlertConditionsResponse, error) {
	return s.client.GetAlertConditions(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}

func (s *alertServiceWrapper) GetAlertConditionsValue(ctx context.Context, req *pb.GetAlertConditionsValueRequest) (*pb.GetAlertConditionsValueResponse, error) {
	return s.client.GetAlertConditionsValue(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}
