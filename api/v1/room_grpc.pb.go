// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: room.proto

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
	Room_ListMyRooms_FullMethodName    = "/v1.Room/ListMyRooms"
	Room_ListAllRooms_FullMethodName   = "/v1.Room/ListAllRooms"
	Room_CreateRoom_FullMethodName     = "/v1.Room/CreateRoom"
	Room_GetRoom_FullMethodName        = "/v1.Room/GetRoom"
	Room_ListRoomMember_FullMethodName = "/v1.Room/ListRoomMember"
)

// RoomClient is the client API for Room service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoomClient interface {
	ListMyRooms(ctx context.Context, in *ListRoomRequest, opts ...grpc.CallOption) (*ListRoomResponse, error)
	ListAllRooms(ctx context.Context, in *ListRoomRequest, opts ...grpc.CallOption) (*ListRoomResponse, error)
	CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*CreateRoomResponse, error)
	GetRoom(ctx context.Context, in *GetRoomRequest, opts ...grpc.CallOption) (*GetRoomResponse, error)
	ListRoomMember(ctx context.Context, in *ListRoomMemberRequest, opts ...grpc.CallOption) (*ListRoomMemberResponse, error)
}

type roomClient struct {
	cc grpc.ClientConnInterface
}

func NewRoomClient(cc grpc.ClientConnInterface) RoomClient {
	return &roomClient{cc}
}

func (c *roomClient) ListMyRooms(ctx context.Context, in *ListRoomRequest, opts ...grpc.CallOption) (*ListRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListRoomResponse)
	err := c.cc.Invoke(ctx, Room_ListMyRooms_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) ListAllRooms(ctx context.Context, in *ListRoomRequest, opts ...grpc.CallOption) (*ListRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListRoomResponse)
	err := c.cc.Invoke(ctx, Room_ListAllRooms_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*CreateRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateRoomResponse)
	err := c.cc.Invoke(ctx, Room_CreateRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) GetRoom(ctx context.Context, in *GetRoomRequest, opts ...grpc.CallOption) (*GetRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoomResponse)
	err := c.cc.Invoke(ctx, Room_GetRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) ListRoomMember(ctx context.Context, in *ListRoomMemberRequest, opts ...grpc.CallOption) (*ListRoomMemberResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListRoomMemberResponse)
	err := c.cc.Invoke(ctx, Room_ListRoomMember_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoomServer is the server API for Room service.
// All implementations must embed UnimplementedRoomServer
// for forward compatibility.
type RoomServer interface {
	ListMyRooms(context.Context, *ListRoomRequest) (*ListRoomResponse, error)
	ListAllRooms(context.Context, *ListRoomRequest) (*ListRoomResponse, error)
	CreateRoom(context.Context, *CreateRoomRequest) (*CreateRoomResponse, error)
	GetRoom(context.Context, *GetRoomRequest) (*GetRoomResponse, error)
	ListRoomMember(context.Context, *ListRoomMemberRequest) (*ListRoomMemberResponse, error)
	mustEmbedUnimplementedRoomServer()
}

// UnimplementedRoomServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRoomServer struct{}

func (UnimplementedRoomServer) ListMyRooms(context.Context, *ListRoomRequest) (*ListRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMyRooms not implemented")
}
func (UnimplementedRoomServer) ListAllRooms(context.Context, *ListRoomRequest) (*ListRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAllRooms not implemented")
}
func (UnimplementedRoomServer) CreateRoom(context.Context, *CreateRoomRequest) (*CreateRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRoom not implemented")
}
func (UnimplementedRoomServer) GetRoom(context.Context, *GetRoomRequest) (*GetRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoom not implemented")
}
func (UnimplementedRoomServer) ListRoomMember(context.Context, *ListRoomMemberRequest) (*ListRoomMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRoomMember not implemented")
}
func (UnimplementedRoomServer) mustEmbedUnimplementedRoomServer() {}
func (UnimplementedRoomServer) testEmbeddedByValue()              {}

// UnsafeRoomServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoomServer will
// result in compilation errors.
type UnsafeRoomServer interface {
	mustEmbedUnimplementedRoomServer()
}

func RegisterRoomServer(s grpc.ServiceRegistrar, srv RoomServer) {
	// If the following call pancis, it indicates UnimplementedRoomServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Room_ServiceDesc, srv)
}

func _Room_ListMyRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).ListMyRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Room_ListMyRooms_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).ListMyRooms(ctx, req.(*ListRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_ListAllRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).ListAllRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Room_ListAllRooms_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).ListAllRooms(ctx, req.(*ListRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_CreateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).CreateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Room_CreateRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).CreateRoom(ctx, req.(*CreateRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_GetRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).GetRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Room_GetRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).GetRoom(ctx, req.(*GetRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_ListRoomMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRoomMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).ListRoomMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Room_ListRoomMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).ListRoomMember(ctx, req.(*ListRoomMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Room_ServiceDesc is the grpc.ServiceDesc for Room service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Room_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.Room",
	HandlerType: (*RoomServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListMyRooms",
			Handler:    _Room_ListMyRooms_Handler,
		},
		{
			MethodName: "ListAllRooms",
			Handler:    _Room_ListAllRooms_Handler,
		},
		{
			MethodName: "CreateRoom",
			Handler:    _Room_CreateRoom_Handler,
		},
		{
			MethodName: "GetRoom",
			Handler:    _Room_GetRoom_Handler,
		},
		{
			MethodName: "ListRoomMember",
			Handler:    _Room_ListRoomMember_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "room.proto",
}
