// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: room_instance.proto

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
	RoomInstance_GetRoomInstance_FullMethodName       = "/v1.RoomInstance/GetRoomInstance"
	RoomInstance_ListGameHistory_FullMethodName       = "/v1.RoomInstance/ListGameHistory"
	RoomInstance_OpenGameConnection_FullMethodName    = "/v1.RoomInstance/OpenGameConnection"
	RoomInstance_SdpAnswer_FullMethodName             = "/v1.RoomInstance/SdpAnswer"
	RoomInstance_AddIceCandidate_FullMethodName       = "/v1.RoomInstance/AddIceCandidate"
	RoomInstance_GetServerIceCandidate_FullMethodName = "/v1.RoomInstance/GetServerIceCandidate"
	RoomInstance_RestartRoomInstance_FullMethodName   = "/v1.RoomInstance/RestartRoomInstance"
	RoomInstance_GetControllerPlayers_FullMethodName  = "/v1.RoomInstance/GetControllerPlayers"
	RoomInstance_SetControllerPlayer_FullMethodName   = "/v1.RoomInstance/SetControllerPlayer"
	RoomInstance_GetGraphicOptions_FullMethodName     = "/v1.RoomInstance/GetGraphicOptions"
	RoomInstance_SetGraphicOptions_FullMethodName     = "/v1.RoomInstance/SetGraphicOptions"
	RoomInstance_GetEmulatorSpeed_FullMethodName      = "/v1.RoomInstance/GetEmulatorSpeed"
	RoomInstance_SetEmulatorSpeed_FullMethodName      = "/v1.RoomInstance/SetEmulatorSpeed"
)

// RoomInstanceClient is the client API for RoomInstance service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoomInstanceClient interface {
	GetRoomInstance(ctx context.Context, in *GetRoomInstanceRequest, opts ...grpc.CallOption) (*GetRoomInstanceResponse, error)
	ListGameHistory(ctx context.Context, in *ListGameHistoryRequest, opts ...grpc.CallOption) (*ListGameHistoryResponse, error)
	OpenGameConnection(ctx context.Context, in *OpenGameConnectionRequest, opts ...grpc.CallOption) (*OpenGameConnectionResponse, error)
	SdpAnswer(ctx context.Context, in *SdpAnswerRequest, opts ...grpc.CallOption) (*SdpAnswerResponse, error)
	AddIceCandidate(ctx context.Context, in *AddIceCandidateRequest, opts ...grpc.CallOption) (*AddIceCandidateResponse, error)
	GetServerIceCandidate(ctx context.Context, in *GetServerIceCandidateRequest, opts ...grpc.CallOption) (*GetServerIceCandidateResponse, error)
	RestartRoomInstance(ctx context.Context, in *RestartRoomInstanceRequest, opts ...grpc.CallOption) (*RestartRoomInstanceResponse, error)
	GetControllerPlayers(ctx context.Context, in *GetControllerPlayersRequest, opts ...grpc.CallOption) (*GetControllerPlayersResponse, error)
	SetControllerPlayer(ctx context.Context, in *SetControllerPlayerRequest, opts ...grpc.CallOption) (*SetControllerPlayerResponse, error)
	GetGraphicOptions(ctx context.Context, in *GetGraphicOptionsRequest, opts ...grpc.CallOption) (*GetGraphicOptionsResponse, error)
	SetGraphicOptions(ctx context.Context, in *SetGraphicOptionsRequest, opts ...grpc.CallOption) (*SetGraphicOptionsResponse, error)
	GetEmulatorSpeed(ctx context.Context, in *GetEmulatorSpeedRequest, opts ...grpc.CallOption) (*GetEmulatorSpeedResponse, error)
	SetEmulatorSpeed(ctx context.Context, in *SetEmulatorSpeedRequest, opts ...grpc.CallOption) (*SetEmulatorSpeedResponse, error)
}

type roomInstanceClient struct {
	cc grpc.ClientConnInterface
}

func NewRoomInstanceClient(cc grpc.ClientConnInterface) RoomInstanceClient {
	return &roomInstanceClient{cc}
}

func (c *roomInstanceClient) GetRoomInstance(ctx context.Context, in *GetRoomInstanceRequest, opts ...grpc.CallOption) (*GetRoomInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoomInstanceResponse)
	err := c.cc.Invoke(ctx, RoomInstance_GetRoomInstance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomInstanceClient) ListGameHistory(ctx context.Context, in *ListGameHistoryRequest, opts ...grpc.CallOption) (*ListGameHistoryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListGameHistoryResponse)
	err := c.cc.Invoke(ctx, RoomInstance_ListGameHistory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomInstanceClient) OpenGameConnection(ctx context.Context, in *OpenGameConnectionRequest, opts ...grpc.CallOption) (*OpenGameConnectionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpenGameConnectionResponse)
	err := c.cc.Invoke(ctx, RoomInstance_OpenGameConnection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomInstanceClient) SdpAnswer(ctx context.Context, in *SdpAnswerRequest, opts ...grpc.CallOption) (*SdpAnswerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SdpAnswerResponse)
	err := c.cc.Invoke(ctx, RoomInstance_SdpAnswer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomInstanceClient) AddIceCandidate(ctx context.Context, in *AddIceCandidateRequest, opts ...grpc.CallOption) (*AddIceCandidateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddIceCandidateResponse)
	err := c.cc.Invoke(ctx, RoomInstance_AddIceCandidate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomInstanceClient) GetServerIceCandidate(ctx context.Context, in *GetServerIceCandidateRequest, opts ...grpc.CallOption) (*GetServerIceCandidateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetServerIceCandidateResponse)
	err := c.cc.Invoke(ctx, RoomInstance_GetServerIceCandidate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomInstanceClient) RestartRoomInstance(ctx context.Context, in *RestartRoomInstanceRequest, opts ...grpc.CallOption) (*RestartRoomInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RestartRoomInstanceResponse)
	err := c.cc.Invoke(ctx, RoomInstance_RestartRoomInstance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomInstanceClient) GetControllerPlayers(ctx context.Context, in *GetControllerPlayersRequest, opts ...grpc.CallOption) (*GetControllerPlayersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetControllerPlayersResponse)
	err := c.cc.Invoke(ctx, RoomInstance_GetControllerPlayers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomInstanceClient) SetControllerPlayer(ctx context.Context, in *SetControllerPlayerRequest, opts ...grpc.CallOption) (*SetControllerPlayerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetControllerPlayerResponse)
	err := c.cc.Invoke(ctx, RoomInstance_SetControllerPlayer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomInstanceClient) GetGraphicOptions(ctx context.Context, in *GetGraphicOptionsRequest, opts ...grpc.CallOption) (*GetGraphicOptionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetGraphicOptionsResponse)
	err := c.cc.Invoke(ctx, RoomInstance_GetGraphicOptions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomInstanceClient) SetGraphicOptions(ctx context.Context, in *SetGraphicOptionsRequest, opts ...grpc.CallOption) (*SetGraphicOptionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetGraphicOptionsResponse)
	err := c.cc.Invoke(ctx, RoomInstance_SetGraphicOptions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomInstanceClient) GetEmulatorSpeed(ctx context.Context, in *GetEmulatorSpeedRequest, opts ...grpc.CallOption) (*GetEmulatorSpeedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetEmulatorSpeedResponse)
	err := c.cc.Invoke(ctx, RoomInstance_GetEmulatorSpeed_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomInstanceClient) SetEmulatorSpeed(ctx context.Context, in *SetEmulatorSpeedRequest, opts ...grpc.CallOption) (*SetEmulatorSpeedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetEmulatorSpeedResponse)
	err := c.cc.Invoke(ctx, RoomInstance_SetEmulatorSpeed_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoomInstanceServer is the server API for RoomInstance service.
// All implementations must embed UnimplementedRoomInstanceServer
// for forward compatibility.
type RoomInstanceServer interface {
	GetRoomInstance(context.Context, *GetRoomInstanceRequest) (*GetRoomInstanceResponse, error)
	ListGameHistory(context.Context, *ListGameHistoryRequest) (*ListGameHistoryResponse, error)
	OpenGameConnection(context.Context, *OpenGameConnectionRequest) (*OpenGameConnectionResponse, error)
	SdpAnswer(context.Context, *SdpAnswerRequest) (*SdpAnswerResponse, error)
	AddIceCandidate(context.Context, *AddIceCandidateRequest) (*AddIceCandidateResponse, error)
	GetServerIceCandidate(context.Context, *GetServerIceCandidateRequest) (*GetServerIceCandidateResponse, error)
	RestartRoomInstance(context.Context, *RestartRoomInstanceRequest) (*RestartRoomInstanceResponse, error)
	GetControllerPlayers(context.Context, *GetControllerPlayersRequest) (*GetControllerPlayersResponse, error)
	SetControllerPlayer(context.Context, *SetControllerPlayerRequest) (*SetControllerPlayerResponse, error)
	GetGraphicOptions(context.Context, *GetGraphicOptionsRequest) (*GetGraphicOptionsResponse, error)
	SetGraphicOptions(context.Context, *SetGraphicOptionsRequest) (*SetGraphicOptionsResponse, error)
	GetEmulatorSpeed(context.Context, *GetEmulatorSpeedRequest) (*GetEmulatorSpeedResponse, error)
	SetEmulatorSpeed(context.Context, *SetEmulatorSpeedRequest) (*SetEmulatorSpeedResponse, error)
	mustEmbedUnimplementedRoomInstanceServer()
}

// UnimplementedRoomInstanceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRoomInstanceServer struct{}

func (UnimplementedRoomInstanceServer) GetRoomInstance(context.Context, *GetRoomInstanceRequest) (*GetRoomInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomInstance not implemented")
}
func (UnimplementedRoomInstanceServer) ListGameHistory(context.Context, *ListGameHistoryRequest) (*ListGameHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListGameHistory not implemented")
}
func (UnimplementedRoomInstanceServer) OpenGameConnection(context.Context, *OpenGameConnectionRequest) (*OpenGameConnectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpenGameConnection not implemented")
}
func (UnimplementedRoomInstanceServer) SdpAnswer(context.Context, *SdpAnswerRequest) (*SdpAnswerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SdpAnswer not implemented")
}
func (UnimplementedRoomInstanceServer) AddIceCandidate(context.Context, *AddIceCandidateRequest) (*AddIceCandidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddIceCandidate not implemented")
}
func (UnimplementedRoomInstanceServer) GetServerIceCandidate(context.Context, *GetServerIceCandidateRequest) (*GetServerIceCandidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetServerIceCandidate not implemented")
}
func (UnimplementedRoomInstanceServer) RestartRoomInstance(context.Context, *RestartRoomInstanceRequest) (*RestartRoomInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RestartRoomInstance not implemented")
}
func (UnimplementedRoomInstanceServer) GetControllerPlayers(context.Context, *GetControllerPlayersRequest) (*GetControllerPlayersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetControllerPlayers not implemented")
}
func (UnimplementedRoomInstanceServer) SetControllerPlayer(context.Context, *SetControllerPlayerRequest) (*SetControllerPlayerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetControllerPlayer not implemented")
}
func (UnimplementedRoomInstanceServer) GetGraphicOptions(context.Context, *GetGraphicOptionsRequest) (*GetGraphicOptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGraphicOptions not implemented")
}
func (UnimplementedRoomInstanceServer) SetGraphicOptions(context.Context, *SetGraphicOptionsRequest) (*SetGraphicOptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetGraphicOptions not implemented")
}
func (UnimplementedRoomInstanceServer) GetEmulatorSpeed(context.Context, *GetEmulatorSpeedRequest) (*GetEmulatorSpeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEmulatorSpeed not implemented")
}
func (UnimplementedRoomInstanceServer) SetEmulatorSpeed(context.Context, *SetEmulatorSpeedRequest) (*SetEmulatorSpeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetEmulatorSpeed not implemented")
}
func (UnimplementedRoomInstanceServer) mustEmbedUnimplementedRoomInstanceServer() {}
func (UnimplementedRoomInstanceServer) testEmbeddedByValue()                      {}

// UnsafeRoomInstanceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoomInstanceServer will
// result in compilation errors.
type UnsafeRoomInstanceServer interface {
	mustEmbedUnimplementedRoomInstanceServer()
}

func RegisterRoomInstanceServer(s grpc.ServiceRegistrar, srv RoomInstanceServer) {
	// If the following call pancis, it indicates UnimplementedRoomInstanceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RoomInstance_ServiceDesc, srv)
}

func _RoomInstance_GetRoomInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomInstanceServer).GetRoomInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomInstance_GetRoomInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomInstanceServer).GetRoomInstance(ctx, req.(*GetRoomInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomInstance_ListGameHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListGameHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomInstanceServer).ListGameHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomInstance_ListGameHistory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomInstanceServer).ListGameHistory(ctx, req.(*ListGameHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomInstance_OpenGameConnection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenGameConnectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomInstanceServer).OpenGameConnection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomInstance_OpenGameConnection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomInstanceServer).OpenGameConnection(ctx, req.(*OpenGameConnectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomInstance_SdpAnswer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SdpAnswerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomInstanceServer).SdpAnswer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomInstance_SdpAnswer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomInstanceServer).SdpAnswer(ctx, req.(*SdpAnswerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomInstance_AddIceCandidate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddIceCandidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomInstanceServer).AddIceCandidate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomInstance_AddIceCandidate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomInstanceServer).AddIceCandidate(ctx, req.(*AddIceCandidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomInstance_GetServerIceCandidate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetServerIceCandidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomInstanceServer).GetServerIceCandidate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomInstance_GetServerIceCandidate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomInstanceServer).GetServerIceCandidate(ctx, req.(*GetServerIceCandidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomInstance_RestartRoomInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestartRoomInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomInstanceServer).RestartRoomInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomInstance_RestartRoomInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomInstanceServer).RestartRoomInstance(ctx, req.(*RestartRoomInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomInstance_GetControllerPlayers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetControllerPlayersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomInstanceServer).GetControllerPlayers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomInstance_GetControllerPlayers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomInstanceServer).GetControllerPlayers(ctx, req.(*GetControllerPlayersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomInstance_SetControllerPlayer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetControllerPlayerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomInstanceServer).SetControllerPlayer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomInstance_SetControllerPlayer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomInstanceServer).SetControllerPlayer(ctx, req.(*SetControllerPlayerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomInstance_GetGraphicOptions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGraphicOptionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomInstanceServer).GetGraphicOptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomInstance_GetGraphicOptions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomInstanceServer).GetGraphicOptions(ctx, req.(*GetGraphicOptionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomInstance_SetGraphicOptions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetGraphicOptionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomInstanceServer).SetGraphicOptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomInstance_SetGraphicOptions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomInstanceServer).SetGraphicOptions(ctx, req.(*SetGraphicOptionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomInstance_GetEmulatorSpeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEmulatorSpeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomInstanceServer).GetEmulatorSpeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomInstance_GetEmulatorSpeed_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomInstanceServer).GetEmulatorSpeed(ctx, req.(*GetEmulatorSpeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomInstance_SetEmulatorSpeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetEmulatorSpeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomInstanceServer).SetEmulatorSpeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomInstance_SetEmulatorSpeed_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomInstanceServer).SetEmulatorSpeed(ctx, req.(*SetEmulatorSpeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RoomInstance_ServiceDesc is the grpc.ServiceDesc for RoomInstance service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RoomInstance_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.RoomInstance",
	HandlerType: (*RoomInstanceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRoomInstance",
			Handler:    _RoomInstance_GetRoomInstance_Handler,
		},
		{
			MethodName: "ListGameHistory",
			Handler:    _RoomInstance_ListGameHistory_Handler,
		},
		{
			MethodName: "OpenGameConnection",
			Handler:    _RoomInstance_OpenGameConnection_Handler,
		},
		{
			MethodName: "SdpAnswer",
			Handler:    _RoomInstance_SdpAnswer_Handler,
		},
		{
			MethodName: "AddIceCandidate",
			Handler:    _RoomInstance_AddIceCandidate_Handler,
		},
		{
			MethodName: "GetServerIceCandidate",
			Handler:    _RoomInstance_GetServerIceCandidate_Handler,
		},
		{
			MethodName: "RestartRoomInstance",
			Handler:    _RoomInstance_RestartRoomInstance_Handler,
		},
		{
			MethodName: "GetControllerPlayers",
			Handler:    _RoomInstance_GetControllerPlayers_Handler,
		},
		{
			MethodName: "SetControllerPlayer",
			Handler:    _RoomInstance_SetControllerPlayer_Handler,
		},
		{
			MethodName: "GetGraphicOptions",
			Handler:    _RoomInstance_GetGraphicOptions_Handler,
		},
		{
			MethodName: "SetGraphicOptions",
			Handler:    _RoomInstance_SetGraphicOptions_Handler,
		},
		{
			MethodName: "GetEmulatorSpeed",
			Handler:    _RoomInstance_GetEmulatorSpeed_Handler,
		},
		{
			MethodName: "SetEmulatorSpeed",
			Handler:    _RoomInstance_SetEmulatorSpeed_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "room_instance.proto",
}
