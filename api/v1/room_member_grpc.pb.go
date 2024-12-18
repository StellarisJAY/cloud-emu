// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: room_member.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RoomMember_ListRoomMember_FullMethodName    = "/v1.RoomMember/ListRoomMember"
	RoomMember_InviteRoomMember_FullMethodName  = "/v1.RoomMember/InviteRoomMember"
	RoomMember_GetUserRoomMember_FullMethodName = "/v1.RoomMember/GetUserRoomMember"
)

// RoomMemberClient is the client API for RoomMember service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoomMemberClient interface {
	ListRoomMember(ctx context.Context, in *ListRoomMemberRequest, opts ...grpc.CallOption) (*ListRoomMemberResponse, error)
	InviteRoomMember(ctx context.Context, in *InviteRoomMemberRequest, opts ...grpc.CallOption) (*InviteRoomMemberResponse, error)
	GetUserRoomMember(ctx context.Context, in *GetUserRoomMemberRequest, opts ...grpc.CallOption) (*GetUserRoomMemberResponse, error)
}

type roomMemberClient struct {
	cc grpc.ClientConnInterface
}

func NewRoomMemberClient(cc grpc.ClientConnInterface) RoomMemberClient {
	return &roomMemberClient{cc}
}

func (c *roomMemberClient) ListRoomMember(ctx context.Context, in *ListRoomMemberRequest, opts ...grpc.CallOption) (*ListRoomMemberResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListRoomMemberResponse)
	err := c.cc.Invoke(ctx, RoomMember_ListRoomMember_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomMemberClient) InviteRoomMember(ctx context.Context, in *InviteRoomMemberRequest, opts ...grpc.CallOption) (*InviteRoomMemberResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(InviteRoomMemberResponse)
	err := c.cc.Invoke(ctx, RoomMember_InviteRoomMember_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomMemberClient) GetUserRoomMember(ctx context.Context, in *GetUserRoomMemberRequest, opts ...grpc.CallOption) (*GetUserRoomMemberResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserRoomMemberResponse)
	err := c.cc.Invoke(ctx, RoomMember_GetUserRoomMember_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoomMemberServer is the server API for RoomMember service.
// All implementations must embed UnimplementedRoomMemberServer
// for forward compatibility.
type RoomMemberServer interface {
	ListRoomMember(context.Context, *ListRoomMemberRequest) (*ListRoomMemberResponse, error)
	InviteRoomMember(context.Context, *InviteRoomMemberRequest) (*InviteRoomMemberResponse, error)
	GetUserRoomMember(context.Context, *GetUserRoomMemberRequest) (*GetUserRoomMemberResponse, error)
	mustEmbedUnimplementedRoomMemberServer()
}

// UnimplementedRoomMemberServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRoomMemberServer struct{}

func (UnimplementedRoomMemberServer) ListRoomMember(context.Context, *ListRoomMemberRequest) (*ListRoomMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRoomMember not implemented")
}
func (UnimplementedRoomMemberServer) InviteRoomMember(context.Context, *InviteRoomMemberRequest) (*InviteRoomMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InviteRoomMember not implemented")
}
func (UnimplementedRoomMemberServer) GetUserRoomMember(context.Context, *GetUserRoomMemberRequest) (*GetUserRoomMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserRoomMember not implemented")
}
func (UnimplementedRoomMemberServer) mustEmbedUnimplementedRoomMemberServer() {}
func (UnimplementedRoomMemberServer) testEmbeddedByValue()                    {}

// UnsafeRoomMemberServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoomMemberServer will
// result in compilation errors.
type UnsafeRoomMemberServer interface {
	mustEmbedUnimplementedRoomMemberServer()
}

func RegisterRoomMemberServer(s grpc.ServiceRegistrar, srv RoomMemberServer) {
	// If the following call pancis, it indicates UnimplementedRoomMemberServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RoomMember_ServiceDesc, srv)
}

func _RoomMember_ListRoomMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRoomMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomMemberServer).ListRoomMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomMember_ListRoomMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomMemberServer).ListRoomMember(ctx, req.(*ListRoomMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomMember_InviteRoomMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InviteRoomMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomMemberServer).InviteRoomMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomMember_InviteRoomMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomMemberServer).InviteRoomMember(ctx, req.(*InviteRoomMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomMember_GetUserRoomMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRoomMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomMemberServer).GetUserRoomMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomMember_GetUserRoomMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomMemberServer).GetUserRoomMember(ctx, req.(*GetUserRoomMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RoomMember_ServiceDesc is the grpc.ServiceDesc for RoomMember service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RoomMember_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.RoomMember",
	HandlerType: (*RoomMemberServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListRoomMember",
			Handler:    _RoomMember_ListRoomMember_Handler,
		},
		{
			MethodName: "InviteRoomMember",
			Handler:    _RoomMember_InviteRoomMember_Handler,
		},
		{
			MethodName: "GetUserRoomMember",
			Handler:    _RoomMember_GetUserRoomMember_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "room_member.proto",
}
