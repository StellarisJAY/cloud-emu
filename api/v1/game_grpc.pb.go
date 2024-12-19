// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: game.proto

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
	Game_OpenGameInstance_FullMethodName           = "/v1.Game/OpenGameInstance"
	Game_GetRoomInstanceAccessToken_FullMethodName = "/v1.Game/GetRoomInstanceAccessToken"
	Game_ShutdownRoomInstance_FullMethodName       = "/v1.Game/ShutdownRoomInstance"
	Game_OpenGameConnection_FullMethodName         = "/v1.Game/OpenGameConnection"
	Game_SdpAnswer_FullMethodName                  = "/v1.Game/SdpAnswer"
	Game_AddIceCandidate_FullMethodName            = "/v1.Game/AddIceCandidate"
	Game_GetIceCandidate_FullMethodName            = "/v1.Game/GetIceCandidate"
)

// GameClient is the client API for Game service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GameClient interface {
	OpenGameInstance(ctx context.Context, in *OpenGameInstanceRequest, opts ...grpc.CallOption) (*OpenGameInstanceResponse, error)
	GetRoomInstanceAccessToken(ctx context.Context, in *GetRoomInstanceAccessTokenRequest, opts ...grpc.CallOption) (*GetRoomInstanceAccessTokenResponse, error)
	ShutdownRoomInstance(ctx context.Context, in *ShutdownRoomInstanceRequest, opts ...grpc.CallOption) (*ShutdownRoomInstanceResponse, error)
	OpenGameConnection(ctx context.Context, in *GameSrvOpenGameConnectionRequest, opts ...grpc.CallOption) (*GameSrvOpenGameConnectionResponse, error)
	SdpAnswer(ctx context.Context, in *GameSrvSdpAnswerRequest, opts ...grpc.CallOption) (*GameSrvSdpAnswerResponse, error)
	AddIceCandidate(ctx context.Context, in *GameSrvAddIceCandidateRequest, opts ...grpc.CallOption) (*GameSrvAddIceCandidateResponse, error)
	GetIceCandidate(ctx context.Context, in *GetIceCandidateRequest, opts ...grpc.CallOption) (*GetIceCandidateResponse, error)
}

type gameClient struct {
	cc grpc.ClientConnInterface
}

func NewGameClient(cc grpc.ClientConnInterface) GameClient {
	return &gameClient{cc}
}

func (c *gameClient) OpenGameInstance(ctx context.Context, in *OpenGameInstanceRequest, opts ...grpc.CallOption) (*OpenGameInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpenGameInstanceResponse)
	err := c.cc.Invoke(ctx, Game_OpenGameInstance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) GetRoomInstanceAccessToken(ctx context.Context, in *GetRoomInstanceAccessTokenRequest, opts ...grpc.CallOption) (*GetRoomInstanceAccessTokenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoomInstanceAccessTokenResponse)
	err := c.cc.Invoke(ctx, Game_GetRoomInstanceAccessToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) ShutdownRoomInstance(ctx context.Context, in *ShutdownRoomInstanceRequest, opts ...grpc.CallOption) (*ShutdownRoomInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShutdownRoomInstanceResponse)
	err := c.cc.Invoke(ctx, Game_ShutdownRoomInstance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) OpenGameConnection(ctx context.Context, in *GameSrvOpenGameConnectionRequest, opts ...grpc.CallOption) (*GameSrvOpenGameConnectionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GameSrvOpenGameConnectionResponse)
	err := c.cc.Invoke(ctx, Game_OpenGameConnection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) SdpAnswer(ctx context.Context, in *GameSrvSdpAnswerRequest, opts ...grpc.CallOption) (*GameSrvSdpAnswerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GameSrvSdpAnswerResponse)
	err := c.cc.Invoke(ctx, Game_SdpAnswer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) AddIceCandidate(ctx context.Context, in *GameSrvAddIceCandidateRequest, opts ...grpc.CallOption) (*GameSrvAddIceCandidateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GameSrvAddIceCandidateResponse)
	err := c.cc.Invoke(ctx, Game_AddIceCandidate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) GetIceCandidate(ctx context.Context, in *GetIceCandidateRequest, opts ...grpc.CallOption) (*GetIceCandidateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetIceCandidateResponse)
	err := c.cc.Invoke(ctx, Game_GetIceCandidate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GameServer is the server API for Game service.
// All implementations must embed UnimplementedGameServer
// for forward compatibility.
type GameServer interface {
	OpenGameInstance(context.Context, *OpenGameInstanceRequest) (*OpenGameInstanceResponse, error)
	GetRoomInstanceAccessToken(context.Context, *GetRoomInstanceAccessTokenRequest) (*GetRoomInstanceAccessTokenResponse, error)
	ShutdownRoomInstance(context.Context, *ShutdownRoomInstanceRequest) (*ShutdownRoomInstanceResponse, error)
	OpenGameConnection(context.Context, *GameSrvOpenGameConnectionRequest) (*GameSrvOpenGameConnectionResponse, error)
	SdpAnswer(context.Context, *GameSrvSdpAnswerRequest) (*GameSrvSdpAnswerResponse, error)
	AddIceCandidate(context.Context, *GameSrvAddIceCandidateRequest) (*GameSrvAddIceCandidateResponse, error)
	GetIceCandidate(context.Context, *GetIceCandidateRequest) (*GetIceCandidateResponse, error)
	mustEmbedUnimplementedGameServer()
}

// UnimplementedGameServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGameServer struct{}

func (UnimplementedGameServer) OpenGameInstance(context.Context, *OpenGameInstanceRequest) (*OpenGameInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpenGameInstance not implemented")
}
func (UnimplementedGameServer) GetRoomInstanceAccessToken(context.Context, *GetRoomInstanceAccessTokenRequest) (*GetRoomInstanceAccessTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomInstanceAccessToken not implemented")
}
func (UnimplementedGameServer) ShutdownRoomInstance(context.Context, *ShutdownRoomInstanceRequest) (*ShutdownRoomInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShutdownRoomInstance not implemented")
}
func (UnimplementedGameServer) OpenGameConnection(context.Context, *GameSrvOpenGameConnectionRequest) (*GameSrvOpenGameConnectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpenGameConnection not implemented")
}
func (UnimplementedGameServer) SdpAnswer(context.Context, *GameSrvSdpAnswerRequest) (*GameSrvSdpAnswerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SdpAnswer not implemented")
}
func (UnimplementedGameServer) AddIceCandidate(context.Context, *GameSrvAddIceCandidateRequest) (*GameSrvAddIceCandidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddIceCandidate not implemented")
}
func (UnimplementedGameServer) GetIceCandidate(context.Context, *GetIceCandidateRequest) (*GetIceCandidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIceCandidate not implemented")
}
func (UnimplementedGameServer) mustEmbedUnimplementedGameServer() {}
func (UnimplementedGameServer) testEmbeddedByValue()              {}

// UnsafeGameServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GameServer will
// result in compilation errors.
type UnsafeGameServer interface {
	mustEmbedUnimplementedGameServer()
}

func RegisterGameServer(s grpc.ServiceRegistrar, srv GameServer) {
	// If the following call pancis, it indicates UnimplementedGameServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Game_ServiceDesc, srv)
}

func _Game_OpenGameInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenGameInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).OpenGameInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Game_OpenGameInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).OpenGameInstance(ctx, req.(*OpenGameInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_GetRoomInstanceAccessToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomInstanceAccessTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).GetRoomInstanceAccessToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Game_GetRoomInstanceAccessToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).GetRoomInstanceAccessToken(ctx, req.(*GetRoomInstanceAccessTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_ShutdownRoomInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShutdownRoomInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).ShutdownRoomInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Game_ShutdownRoomInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).ShutdownRoomInstance(ctx, req.(*ShutdownRoomInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_OpenGameConnection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GameSrvOpenGameConnectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).OpenGameConnection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Game_OpenGameConnection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).OpenGameConnection(ctx, req.(*GameSrvOpenGameConnectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_SdpAnswer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GameSrvSdpAnswerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).SdpAnswer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Game_SdpAnswer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).SdpAnswer(ctx, req.(*GameSrvSdpAnswerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_AddIceCandidate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GameSrvAddIceCandidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).AddIceCandidate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Game_AddIceCandidate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).AddIceCandidate(ctx, req.(*GameSrvAddIceCandidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_GetIceCandidate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIceCandidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).GetIceCandidate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Game_GetIceCandidate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).GetIceCandidate(ctx, req.(*GetIceCandidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Game_ServiceDesc is the grpc.ServiceDesc for Game service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Game_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.Game",
	HandlerType: (*GameServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OpenGameInstance",
			Handler:    _Game_OpenGameInstance_Handler,
		},
		{
			MethodName: "GetRoomInstanceAccessToken",
			Handler:    _Game_GetRoomInstanceAccessToken_Handler,
		},
		{
			MethodName: "ShutdownRoomInstance",
			Handler:    _Game_ShutdownRoomInstance_Handler,
		},
		{
			MethodName: "OpenGameConnection",
			Handler:    _Game_OpenGameConnection_Handler,
		},
		{
			MethodName: "SdpAnswer",
			Handler:    _Game_SdpAnswer_Handler,
		},
		{
			MethodName: "AddIceCandidate",
			Handler:    _Game_AddIceCandidate_Handler,
		},
		{
			MethodName: "GetIceCandidate",
			Handler:    _Game_GetIceCandidate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "game.proto",
}
