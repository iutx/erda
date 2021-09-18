// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// Source: checker_v1.proto

package pb

import (
	context "context"

	transport "github.com/erda-project/erda-infra/pkg/transport"
	grpc1 "github.com/erda-project/erda-infra/pkg/transport/grpc"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion5

// CheckerV1ServiceClient is the client API for CheckerV1Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CheckerV1ServiceClient interface {
	CreateCheckerV1(ctx context.Context, in *CreateCheckerV1Request, opts ...grpc.CallOption) (*CreateCheckerV1Response, error)
	UpdateCheckerV1(ctx context.Context, in *UpdateCheckerV1Request, opts ...grpc.CallOption) (*UpdateCheckerV1Response, error)
	DeleteCheckerV1(ctx context.Context, in *DeleteCheckerV1Request, opts ...grpc.CallOption) (*DeleteCheckerV1Response, error)
	GetCheckerV1(ctx context.Context, in *GetCheckerV1Request, opts ...grpc.CallOption) (*GetCheckerV1Response, error)
	DescribeCheckersV1(ctx context.Context, in *DescribeCheckersV1Request, opts ...grpc.CallOption) (*DescribeCheckersV1Response, error)
	DescribeCheckerV1(ctx context.Context, in *DescribeCheckerV1Request, opts ...grpc.CallOption) (*DescribeCheckerV1Response, error)
	GetCheckerStatusV1(ctx context.Context, in *GetCheckerStatusV1Request, opts ...grpc.CallOption) (*GetCheckerStatusV1Response, error)
	// +depracated
	GetCheckerIssuesV1(ctx context.Context, in *GetCheckerIssuesV1Request, opts ...grpc.CallOption) (*GetCheckerIssuesV1Response, error)
}

type checkerV1ServiceClient struct {
	cc grpc1.ClientConnInterface
}

func NewCheckerV1ServiceClient(cc grpc1.ClientConnInterface) CheckerV1ServiceClient {
	return &checkerV1ServiceClient{cc}
}

func (c *checkerV1ServiceClient) CreateCheckerV1(ctx context.Context, in *CreateCheckerV1Request, opts ...grpc.CallOption) (*CreateCheckerV1Response, error) {
	out := new(CreateCheckerV1Response)
	err := c.cc.Invoke(ctx, "/erda.msp.apm.checker.CheckerV1Service/CreateCheckerV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkerV1ServiceClient) UpdateCheckerV1(ctx context.Context, in *UpdateCheckerV1Request, opts ...grpc.CallOption) (*UpdateCheckerV1Response, error) {
	out := new(UpdateCheckerV1Response)
	err := c.cc.Invoke(ctx, "/erda.msp.apm.checker.CheckerV1Service/UpdateCheckerV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkerV1ServiceClient) DeleteCheckerV1(ctx context.Context, in *DeleteCheckerV1Request, opts ...grpc.CallOption) (*DeleteCheckerV1Response, error) {
	out := new(DeleteCheckerV1Response)
	err := c.cc.Invoke(ctx, "/erda.msp.apm.checker.CheckerV1Service/DeleteCheckerV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkerV1ServiceClient) GetCheckerV1(ctx context.Context, in *GetCheckerV1Request, opts ...grpc.CallOption) (*GetCheckerV1Response, error) {
	out := new(GetCheckerV1Response)
	err := c.cc.Invoke(ctx, "/erda.msp.apm.checker.CheckerV1Service/GetCheckerV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkerV1ServiceClient) DescribeCheckersV1(ctx context.Context, in *DescribeCheckersV1Request, opts ...grpc.CallOption) (*DescribeCheckersV1Response, error) {
	out := new(DescribeCheckersV1Response)
	err := c.cc.Invoke(ctx, "/erda.msp.apm.checker.CheckerV1Service/DescribeCheckersV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkerV1ServiceClient) DescribeCheckerV1(ctx context.Context, in *DescribeCheckerV1Request, opts ...grpc.CallOption) (*DescribeCheckerV1Response, error) {
	out := new(DescribeCheckerV1Response)
	err := c.cc.Invoke(ctx, "/erda.msp.apm.checker.CheckerV1Service/DescribeCheckerV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkerV1ServiceClient) GetCheckerStatusV1(ctx context.Context, in *GetCheckerStatusV1Request, opts ...grpc.CallOption) (*GetCheckerStatusV1Response, error) {
	out := new(GetCheckerStatusV1Response)
	err := c.cc.Invoke(ctx, "/erda.msp.apm.checker.CheckerV1Service/GetCheckerStatusV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkerV1ServiceClient) GetCheckerIssuesV1(ctx context.Context, in *GetCheckerIssuesV1Request, opts ...grpc.CallOption) (*GetCheckerIssuesV1Response, error) {
	out := new(GetCheckerIssuesV1Response)
	err := c.cc.Invoke(ctx, "/erda.msp.apm.checker.CheckerV1Service/GetCheckerIssuesV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CheckerV1ServiceServer is the server API for CheckerV1Service service.
// All implementations should embed UnimplementedCheckerV1ServiceServer
// for forward compatibility
type CheckerV1ServiceServer interface {
	CreateCheckerV1(context.Context, *CreateCheckerV1Request) (*CreateCheckerV1Response, error)
	UpdateCheckerV1(context.Context, *UpdateCheckerV1Request) (*UpdateCheckerV1Response, error)
	DeleteCheckerV1(context.Context, *DeleteCheckerV1Request) (*DeleteCheckerV1Response, error)
	GetCheckerV1(context.Context, *GetCheckerV1Request) (*GetCheckerV1Response, error)
	DescribeCheckersV1(context.Context, *DescribeCheckersV1Request) (*DescribeCheckersV1Response, error)
	DescribeCheckerV1(context.Context, *DescribeCheckerV1Request) (*DescribeCheckerV1Response, error)
	GetCheckerStatusV1(context.Context, *GetCheckerStatusV1Request) (*GetCheckerStatusV1Response, error)
	// +depracated
	GetCheckerIssuesV1(context.Context, *GetCheckerIssuesV1Request) (*GetCheckerIssuesV1Response, error)
}

// UnimplementedCheckerV1ServiceServer should be embedded to have forward compatible implementations.
type UnimplementedCheckerV1ServiceServer struct {
}

func (*UnimplementedCheckerV1ServiceServer) CreateCheckerV1(context.Context, *CreateCheckerV1Request) (*CreateCheckerV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCheckerV1 not implemented")
}
func (*UnimplementedCheckerV1ServiceServer) UpdateCheckerV1(context.Context, *UpdateCheckerV1Request) (*UpdateCheckerV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCheckerV1 not implemented")
}
func (*UnimplementedCheckerV1ServiceServer) DeleteCheckerV1(context.Context, *DeleteCheckerV1Request) (*DeleteCheckerV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCheckerV1 not implemented")
}
func (*UnimplementedCheckerV1ServiceServer) GetCheckerV1(context.Context, *GetCheckerV1Request) (*GetCheckerV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCheckerV1 not implemented")
}
func (*UnimplementedCheckerV1ServiceServer) DescribeCheckersV1(context.Context, *DescribeCheckersV1Request) (*DescribeCheckersV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeCheckersV1 not implemented")
}
func (*UnimplementedCheckerV1ServiceServer) DescribeCheckerV1(context.Context, *DescribeCheckerV1Request) (*DescribeCheckerV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeCheckerV1 not implemented")
}
func (*UnimplementedCheckerV1ServiceServer) GetCheckerStatusV1(context.Context, *GetCheckerStatusV1Request) (*GetCheckerStatusV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCheckerStatusV1 not implemented")
}
func (*UnimplementedCheckerV1ServiceServer) GetCheckerIssuesV1(context.Context, *GetCheckerIssuesV1Request) (*GetCheckerIssuesV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCheckerIssuesV1 not implemented")
}

func RegisterCheckerV1ServiceServer(s grpc1.ServiceRegistrar, srv CheckerV1ServiceServer, opts ...grpc1.HandleOption) {
	s.RegisterService(_get_CheckerV1Service_serviceDesc(srv, opts...), srv)
}

var _CheckerV1Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "erda.msp.apm.checker.CheckerV1Service",
	HandlerType: (*CheckerV1ServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "checker_v1.proto",
}

func _get_CheckerV1Service_serviceDesc(srv CheckerV1ServiceServer, opts ...grpc1.HandleOption) *grpc.ServiceDesc {
	h := grpc1.DefaultHandleOptions()
	for _, op := range opts {
		op(h)
	}

	_CheckerV1Service_CreateCheckerV1_Handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.CreateCheckerV1(ctx, req.(*CreateCheckerV1Request))
	}
	var _CheckerV1Service_CreateCheckerV1_info transport.ServiceInfo
	if h.Interceptor != nil {
		_CheckerV1Service_CreateCheckerV1_info = transport.NewServiceInfo("erda.msp.apm.checker.CheckerV1Service", "CreateCheckerV1", srv)
		_CheckerV1Service_CreateCheckerV1_Handler = h.Interceptor(_CheckerV1Service_CreateCheckerV1_Handler)
	}

	_CheckerV1Service_UpdateCheckerV1_Handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.UpdateCheckerV1(ctx, req.(*UpdateCheckerV1Request))
	}
	var _CheckerV1Service_UpdateCheckerV1_info transport.ServiceInfo
	if h.Interceptor != nil {
		_CheckerV1Service_UpdateCheckerV1_info = transport.NewServiceInfo("erda.msp.apm.checker.CheckerV1Service", "UpdateCheckerV1", srv)
		_CheckerV1Service_UpdateCheckerV1_Handler = h.Interceptor(_CheckerV1Service_UpdateCheckerV1_Handler)
	}

	_CheckerV1Service_DeleteCheckerV1_Handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.DeleteCheckerV1(ctx, req.(*DeleteCheckerV1Request))
	}
	var _CheckerV1Service_DeleteCheckerV1_info transport.ServiceInfo
	if h.Interceptor != nil {
		_CheckerV1Service_DeleteCheckerV1_info = transport.NewServiceInfo("erda.msp.apm.checker.CheckerV1Service", "DeleteCheckerV1", srv)
		_CheckerV1Service_DeleteCheckerV1_Handler = h.Interceptor(_CheckerV1Service_DeleteCheckerV1_Handler)
	}

	_CheckerV1Service_GetCheckerV1_Handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.GetCheckerV1(ctx, req.(*GetCheckerV1Request))
	}
	var _CheckerV1Service_GetCheckerV1_info transport.ServiceInfo
	if h.Interceptor != nil {
		_CheckerV1Service_GetCheckerV1_info = transport.NewServiceInfo("erda.msp.apm.checker.CheckerV1Service", "GetCheckerV1", srv)
		_CheckerV1Service_GetCheckerV1_Handler = h.Interceptor(_CheckerV1Service_GetCheckerV1_Handler)
	}

	_CheckerV1Service_DescribeCheckersV1_Handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.DescribeCheckersV1(ctx, req.(*DescribeCheckersV1Request))
	}
	var _CheckerV1Service_DescribeCheckersV1_info transport.ServiceInfo
	if h.Interceptor != nil {
		_CheckerV1Service_DescribeCheckersV1_info = transport.NewServiceInfo("erda.msp.apm.checker.CheckerV1Service", "DescribeCheckersV1", srv)
		_CheckerV1Service_DescribeCheckersV1_Handler = h.Interceptor(_CheckerV1Service_DescribeCheckersV1_Handler)
	}

	_CheckerV1Service_DescribeCheckerV1_Handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.DescribeCheckerV1(ctx, req.(*DescribeCheckerV1Request))
	}
	var _CheckerV1Service_DescribeCheckerV1_info transport.ServiceInfo
	if h.Interceptor != nil {
		_CheckerV1Service_DescribeCheckerV1_info = transport.NewServiceInfo("erda.msp.apm.checker.CheckerV1Service", "DescribeCheckerV1", srv)
		_CheckerV1Service_DescribeCheckerV1_Handler = h.Interceptor(_CheckerV1Service_DescribeCheckerV1_Handler)
	}

	_CheckerV1Service_GetCheckerStatusV1_Handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.GetCheckerStatusV1(ctx, req.(*GetCheckerStatusV1Request))
	}
	var _CheckerV1Service_GetCheckerStatusV1_info transport.ServiceInfo
	if h.Interceptor != nil {
		_CheckerV1Service_GetCheckerStatusV1_info = transport.NewServiceInfo("erda.msp.apm.checker.CheckerV1Service", "GetCheckerStatusV1", srv)
		_CheckerV1Service_GetCheckerStatusV1_Handler = h.Interceptor(_CheckerV1Service_GetCheckerStatusV1_Handler)
	}

	_CheckerV1Service_GetCheckerIssuesV1_Handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.GetCheckerIssuesV1(ctx, req.(*GetCheckerIssuesV1Request))
	}
	var _CheckerV1Service_GetCheckerIssuesV1_info transport.ServiceInfo
	if h.Interceptor != nil {
		_CheckerV1Service_GetCheckerIssuesV1_info = transport.NewServiceInfo("erda.msp.apm.checker.CheckerV1Service", "GetCheckerIssuesV1", srv)
		_CheckerV1Service_GetCheckerIssuesV1_Handler = h.Interceptor(_CheckerV1Service_GetCheckerIssuesV1_Handler)
	}

	var serviceDesc = _CheckerV1Service_serviceDesc
	serviceDesc.Methods = []grpc.MethodDesc{
		{
			MethodName: "CreateCheckerV1",
			Handler: func(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(CreateCheckerV1Request)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil && h.Interceptor == nil {
					return srv.(CheckerV1ServiceServer).CreateCheckerV1(ctx, in)
				}
				if h.Interceptor != nil {
					ctx = context.WithValue(ctx, transport.ServiceInfoContextKey, _CheckerV1Service_CreateCheckerV1_info)
				}
				if interceptor == nil {
					return _CheckerV1Service_CreateCheckerV1_Handler(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/erda.msp.apm.checker.CheckerV1Service/CreateCheckerV1",
				}
				return interceptor(ctx, in, info, _CheckerV1Service_CreateCheckerV1_Handler)
			},
		},
		{
			MethodName: "UpdateCheckerV1",
			Handler: func(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(UpdateCheckerV1Request)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil && h.Interceptor == nil {
					return srv.(CheckerV1ServiceServer).UpdateCheckerV1(ctx, in)
				}
				if h.Interceptor != nil {
					ctx = context.WithValue(ctx, transport.ServiceInfoContextKey, _CheckerV1Service_UpdateCheckerV1_info)
				}
				if interceptor == nil {
					return _CheckerV1Service_UpdateCheckerV1_Handler(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/erda.msp.apm.checker.CheckerV1Service/UpdateCheckerV1",
				}
				return interceptor(ctx, in, info, _CheckerV1Service_UpdateCheckerV1_Handler)
			},
		},
		{
			MethodName: "DeleteCheckerV1",
			Handler: func(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(DeleteCheckerV1Request)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil && h.Interceptor == nil {
					return srv.(CheckerV1ServiceServer).DeleteCheckerV1(ctx, in)
				}
				if h.Interceptor != nil {
					ctx = context.WithValue(ctx, transport.ServiceInfoContextKey, _CheckerV1Service_DeleteCheckerV1_info)
				}
				if interceptor == nil {
					return _CheckerV1Service_DeleteCheckerV1_Handler(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/erda.msp.apm.checker.CheckerV1Service/DeleteCheckerV1",
				}
				return interceptor(ctx, in, info, _CheckerV1Service_DeleteCheckerV1_Handler)
			},
		},
		{
			MethodName: "GetCheckerV1",
			Handler: func(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(GetCheckerV1Request)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil && h.Interceptor == nil {
					return srv.(CheckerV1ServiceServer).GetCheckerV1(ctx, in)
				}
				if h.Interceptor != nil {
					ctx = context.WithValue(ctx, transport.ServiceInfoContextKey, _CheckerV1Service_GetCheckerV1_info)
				}
				if interceptor == nil {
					return _CheckerV1Service_GetCheckerV1_Handler(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/erda.msp.apm.checker.CheckerV1Service/GetCheckerV1",
				}
				return interceptor(ctx, in, info, _CheckerV1Service_GetCheckerV1_Handler)
			},
		},
		{
			MethodName: "DescribeCheckersV1",
			Handler: func(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(DescribeCheckersV1Request)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil && h.Interceptor == nil {
					return srv.(CheckerV1ServiceServer).DescribeCheckersV1(ctx, in)
				}
				if h.Interceptor != nil {
					ctx = context.WithValue(ctx, transport.ServiceInfoContextKey, _CheckerV1Service_DescribeCheckersV1_info)
				}
				if interceptor == nil {
					return _CheckerV1Service_DescribeCheckersV1_Handler(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/erda.msp.apm.checker.CheckerV1Service/DescribeCheckersV1",
				}
				return interceptor(ctx, in, info, _CheckerV1Service_DescribeCheckersV1_Handler)
			},
		},
		{
			MethodName: "DescribeCheckerV1",
			Handler: func(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(DescribeCheckerV1Request)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil && h.Interceptor == nil {
					return srv.(CheckerV1ServiceServer).DescribeCheckerV1(ctx, in)
				}
				if h.Interceptor != nil {
					ctx = context.WithValue(ctx, transport.ServiceInfoContextKey, _CheckerV1Service_DescribeCheckerV1_info)
				}
				if interceptor == nil {
					return _CheckerV1Service_DescribeCheckerV1_Handler(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/erda.msp.apm.checker.CheckerV1Service/DescribeCheckerV1",
				}
				return interceptor(ctx, in, info, _CheckerV1Service_DescribeCheckerV1_Handler)
			},
		},
		{
			MethodName: "GetCheckerStatusV1",
			Handler: func(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(GetCheckerStatusV1Request)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil && h.Interceptor == nil {
					return srv.(CheckerV1ServiceServer).GetCheckerStatusV1(ctx, in)
				}
				if h.Interceptor != nil {
					ctx = context.WithValue(ctx, transport.ServiceInfoContextKey, _CheckerV1Service_GetCheckerStatusV1_info)
				}
				if interceptor == nil {
					return _CheckerV1Service_GetCheckerStatusV1_Handler(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/erda.msp.apm.checker.CheckerV1Service/GetCheckerStatusV1",
				}
				return interceptor(ctx, in, info, _CheckerV1Service_GetCheckerStatusV1_Handler)
			},
		},
		{
			MethodName: "GetCheckerIssuesV1",
			Handler: func(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(GetCheckerIssuesV1Request)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil && h.Interceptor == nil {
					return srv.(CheckerV1ServiceServer).GetCheckerIssuesV1(ctx, in)
				}
				if h.Interceptor != nil {
					ctx = context.WithValue(ctx, transport.ServiceInfoContextKey, _CheckerV1Service_GetCheckerIssuesV1_info)
				}
				if interceptor == nil {
					return _CheckerV1Service_GetCheckerIssuesV1_Handler(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/erda.msp.apm.checker.CheckerV1Service/GetCheckerIssuesV1",
				}
				return interceptor(ctx, in, info, _CheckerV1Service_GetCheckerIssuesV1_Handler)
			},
		},
	}
	return &serviceDesc
}
