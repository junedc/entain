// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package sport

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SportClient is the client API for Sport service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SportClient interface {
	// ListEvents will return a collection of all races.
	ListEvents(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*ListEventsResponse, error)
	// GetEvent will return a race by id
	GetEvent(ctx context.Context, in *GetEventRequest, opts ...grpc.CallOption) (*Event, error)
}

type sportClient struct {
	cc grpc.ClientConnInterface
}

func NewSportClient(cc grpc.ClientConnInterface) SportClient {
	return &sportClient{cc}
}

func (c *sportClient) ListEvents(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*ListEventsResponse, error) {
	out := new(ListEventsResponse)
	err := c.cc.Invoke(ctx, "/sport.Sport/ListEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sportClient) GetEvent(ctx context.Context, in *GetEventRequest, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/sport.Sport/GetEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SportServer is the server API for Sport service.
// All implementations should embed UnimplementedSportServer
// for forward compatibility
type SportServer interface {
	// ListEvents will return a collection of all races.
	ListEvents(context.Context, *ListEventsRequest) (*ListEventsResponse, error)
	// GetEvent will return a race by id
	GetEvent(context.Context, *GetEventRequest) (*Event, error)
}

// UnimplementedSportServer should be embedded to have forward compatible implementations.
type UnimplementedSportServer struct {
}

func (UnimplementedSportServer) ListEvents(context.Context, *ListEventsRequest) (*ListEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEvents not implemented")
}
func (UnimplementedSportServer) GetEvent(context.Context, *GetEventRequest) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvent not implemented")
}

// UnsafeSportServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SportServer will
// result in compilation errors.
type UnsafeSportServer interface {
	mustEmbedUnimplementedSportServer()
}

func RegisterSportServer(s grpc.ServiceRegistrar, srv SportServer) {
	s.RegisterService(&Sport_ServiceDesc, srv)
}

func _Sport_ListEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SportServer).ListEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sport.Sport/ListEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SportServer).ListEvents(ctx, req.(*ListEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sport_GetEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SportServer).GetEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sport.Sport/GetEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SportServer).GetEvent(ctx, req.(*GetEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Sport_ServiceDesc is the grpc.ServiceDesc for Sport service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sport_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sport.Sport",
	HandlerType: (*SportServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListEvents",
			Handler:    _Sport_ListEvents_Handler,
		},
		{
			MethodName: "GetEvent",
			Handler:    _Sport_GetEvent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sport/sport.proto",
}
