// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: keyboard_binding.proto

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
	KeyboardBinding_ListKeyboardBinding_FullMethodName   = "/v1.KeyboardBinding/ListKeyboardBinding"
	KeyboardBinding_CreateKeyboardBinding_FullMethodName = "/v1.KeyboardBinding/CreateKeyboardBinding"
	KeyboardBinding_UpdateKeyboardBinding_FullMethodName = "/v1.KeyboardBinding/UpdateKeyboardBinding"
	KeyboardBinding_DeleteKeyboardBinding_FullMethodName = "/v1.KeyboardBinding/DeleteKeyboardBinding"
)

// KeyboardBindingClient is the client API for KeyboardBinding service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeyboardBindingClient interface {
	ListKeyboardBinding(ctx context.Context, in *ListKeyboardBindingRequest, opts ...grpc.CallOption) (*ListKeyboardBindingResponse, error)
	CreateKeyboardBinding(ctx context.Context, in *CreateKeyboardBindingRequest, opts ...grpc.CallOption) (*CreateKeyboardBindingResponse, error)
	UpdateKeyboardBinding(ctx context.Context, in *UpdateKeyboardBindingRequest, opts ...grpc.CallOption) (*UpdateKeyboardBindingResponse, error)
	DeleteKeyboardBinding(ctx context.Context, in *DeleteKeyboardBindingRequest, opts ...grpc.CallOption) (*DeleteKeyboardBindingResponse, error)
}

type keyboardBindingClient struct {
	cc grpc.ClientConnInterface
}

func NewKeyboardBindingClient(cc grpc.ClientConnInterface) KeyboardBindingClient {
	return &keyboardBindingClient{cc}
}

func (c *keyboardBindingClient) ListKeyboardBinding(ctx context.Context, in *ListKeyboardBindingRequest, opts ...grpc.CallOption) (*ListKeyboardBindingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListKeyboardBindingResponse)
	err := c.cc.Invoke(ctx, KeyboardBinding_ListKeyboardBinding_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyboardBindingClient) CreateKeyboardBinding(ctx context.Context, in *CreateKeyboardBindingRequest, opts ...grpc.CallOption) (*CreateKeyboardBindingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateKeyboardBindingResponse)
	err := c.cc.Invoke(ctx, KeyboardBinding_CreateKeyboardBinding_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyboardBindingClient) UpdateKeyboardBinding(ctx context.Context, in *UpdateKeyboardBindingRequest, opts ...grpc.CallOption) (*UpdateKeyboardBindingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateKeyboardBindingResponse)
	err := c.cc.Invoke(ctx, KeyboardBinding_UpdateKeyboardBinding_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyboardBindingClient) DeleteKeyboardBinding(ctx context.Context, in *DeleteKeyboardBindingRequest, opts ...grpc.CallOption) (*DeleteKeyboardBindingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteKeyboardBindingResponse)
	err := c.cc.Invoke(ctx, KeyboardBinding_DeleteKeyboardBinding_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeyboardBindingServer is the server API for KeyboardBinding service.
// All implementations must embed UnimplementedKeyboardBindingServer
// for forward compatibility.
type KeyboardBindingServer interface {
	ListKeyboardBinding(context.Context, *ListKeyboardBindingRequest) (*ListKeyboardBindingResponse, error)
	CreateKeyboardBinding(context.Context, *CreateKeyboardBindingRequest) (*CreateKeyboardBindingResponse, error)
	UpdateKeyboardBinding(context.Context, *UpdateKeyboardBindingRequest) (*UpdateKeyboardBindingResponse, error)
	DeleteKeyboardBinding(context.Context, *DeleteKeyboardBindingRequest) (*DeleteKeyboardBindingResponse, error)
	mustEmbedUnimplementedKeyboardBindingServer()
}

// UnimplementedKeyboardBindingServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedKeyboardBindingServer struct{}

func (UnimplementedKeyboardBindingServer) ListKeyboardBinding(context.Context, *ListKeyboardBindingRequest) (*ListKeyboardBindingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListKeyboardBinding not implemented")
}
func (UnimplementedKeyboardBindingServer) CreateKeyboardBinding(context.Context, *CreateKeyboardBindingRequest) (*CreateKeyboardBindingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateKeyboardBinding not implemented")
}
func (UnimplementedKeyboardBindingServer) UpdateKeyboardBinding(context.Context, *UpdateKeyboardBindingRequest) (*UpdateKeyboardBindingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateKeyboardBinding not implemented")
}
func (UnimplementedKeyboardBindingServer) DeleteKeyboardBinding(context.Context, *DeleteKeyboardBindingRequest) (*DeleteKeyboardBindingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteKeyboardBinding not implemented")
}
func (UnimplementedKeyboardBindingServer) mustEmbedUnimplementedKeyboardBindingServer() {}
func (UnimplementedKeyboardBindingServer) testEmbeddedByValue()                         {}

// UnsafeKeyboardBindingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeyboardBindingServer will
// result in compilation errors.
type UnsafeKeyboardBindingServer interface {
	mustEmbedUnimplementedKeyboardBindingServer()
}

func RegisterKeyboardBindingServer(s grpc.ServiceRegistrar, srv KeyboardBindingServer) {
	// If the following call pancis, it indicates UnimplementedKeyboardBindingServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&KeyboardBinding_ServiceDesc, srv)
}

func _KeyboardBinding_ListKeyboardBinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListKeyboardBindingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyboardBindingServer).ListKeyboardBinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeyboardBinding_ListKeyboardBinding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyboardBindingServer).ListKeyboardBinding(ctx, req.(*ListKeyboardBindingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyboardBinding_CreateKeyboardBinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateKeyboardBindingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyboardBindingServer).CreateKeyboardBinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeyboardBinding_CreateKeyboardBinding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyboardBindingServer).CreateKeyboardBinding(ctx, req.(*CreateKeyboardBindingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyboardBinding_UpdateKeyboardBinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateKeyboardBindingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyboardBindingServer).UpdateKeyboardBinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeyboardBinding_UpdateKeyboardBinding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyboardBindingServer).UpdateKeyboardBinding(ctx, req.(*UpdateKeyboardBindingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyboardBinding_DeleteKeyboardBinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteKeyboardBindingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyboardBindingServer).DeleteKeyboardBinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeyboardBinding_DeleteKeyboardBinding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyboardBindingServer).DeleteKeyboardBinding(ctx, req.(*DeleteKeyboardBindingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KeyboardBinding_ServiceDesc is the grpc.ServiceDesc for KeyboardBinding service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KeyboardBinding_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.KeyboardBinding",
	HandlerType: (*KeyboardBindingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListKeyboardBinding",
			Handler:    _KeyboardBinding_ListKeyboardBinding_Handler,
		},
		{
			MethodName: "CreateKeyboardBinding",
			Handler:    _KeyboardBinding_CreateKeyboardBinding_Handler,
		},
		{
			MethodName: "UpdateKeyboardBinding",
			Handler:    _KeyboardBinding_UpdateKeyboardBinding_Handler,
		},
		{
			MethodName: "DeleteKeyboardBinding",
			Handler:    _KeyboardBinding_DeleteKeyboardBinding_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "keyboard_binding.proto",
}
