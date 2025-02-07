// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: notification.proto

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
	Notification_ListInboxNotifications_FullMethodName   = "/v1.Notification/ListInboxNotifications"
	Notification_DeleteInboxNotifications_FullMethodName = "/v1.Notification/DeleteInboxNotifications"
	Notification_ClearInbox_FullMethodName               = "/v1.Notification/ClearInbox"
)

// NotificationClient is the client API for Notification service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificationClient interface {
	ListInboxNotifications(ctx context.Context, in *ListInboxNotificationRequest, opts ...grpc.CallOption) (*ListInboxNotificationResponse, error)
	DeleteInboxNotifications(ctx context.Context, in *DeleteInboxNotificationRequest, opts ...grpc.CallOption) (*DeleteInboxNotificationResponse, error)
	ClearInbox(ctx context.Context, in *ClearInboxRequest, opts ...grpc.CallOption) (*ClearInboxResponse, error)
}

type notificationClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationClient(cc grpc.ClientConnInterface) NotificationClient {
	return &notificationClient{cc}
}

func (c *notificationClient) ListInboxNotifications(ctx context.Context, in *ListInboxNotificationRequest, opts ...grpc.CallOption) (*ListInboxNotificationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListInboxNotificationResponse)
	err := c.cc.Invoke(ctx, Notification_ListInboxNotifications_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationClient) DeleteInboxNotifications(ctx context.Context, in *DeleteInboxNotificationRequest, opts ...grpc.CallOption) (*DeleteInboxNotificationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteInboxNotificationResponse)
	err := c.cc.Invoke(ctx, Notification_DeleteInboxNotifications_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationClient) ClearInbox(ctx context.Context, in *ClearInboxRequest, opts ...grpc.CallOption) (*ClearInboxResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClearInboxResponse)
	err := c.cc.Invoke(ctx, Notification_ClearInbox_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationServer is the server API for Notification service.
// All implementations must embed UnimplementedNotificationServer
// for forward compatibility.
type NotificationServer interface {
	ListInboxNotifications(context.Context, *ListInboxNotificationRequest) (*ListInboxNotificationResponse, error)
	DeleteInboxNotifications(context.Context, *DeleteInboxNotificationRequest) (*DeleteInboxNotificationResponse, error)
	ClearInbox(context.Context, *ClearInboxRequest) (*ClearInboxResponse, error)
	mustEmbedUnimplementedNotificationServer()
}

// UnimplementedNotificationServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNotificationServer struct{}

func (UnimplementedNotificationServer) ListInboxNotifications(context.Context, *ListInboxNotificationRequest) (*ListInboxNotificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListInboxNotifications not implemented")
}
func (UnimplementedNotificationServer) DeleteInboxNotifications(context.Context, *DeleteInboxNotificationRequest) (*DeleteInboxNotificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteInboxNotifications not implemented")
}
func (UnimplementedNotificationServer) ClearInbox(context.Context, *ClearInboxRequest) (*ClearInboxResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearInbox not implemented")
}
func (UnimplementedNotificationServer) mustEmbedUnimplementedNotificationServer() {}
func (UnimplementedNotificationServer) testEmbeddedByValue()                      {}

// UnsafeNotificationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServer will
// result in compilation errors.
type UnsafeNotificationServer interface {
	mustEmbedUnimplementedNotificationServer()
}

func RegisterNotificationServer(s grpc.ServiceRegistrar, srv NotificationServer) {
	// If the following call pancis, it indicates UnimplementedNotificationServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Notification_ServiceDesc, srv)
}

func _Notification_ListInboxNotifications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListInboxNotificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServer).ListInboxNotifications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Notification_ListInboxNotifications_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServer).ListInboxNotifications(ctx, req.(*ListInboxNotificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notification_DeleteInboxNotifications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteInboxNotificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServer).DeleteInboxNotifications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Notification_DeleteInboxNotifications_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServer).DeleteInboxNotifications(ctx, req.(*DeleteInboxNotificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notification_ClearInbox_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearInboxRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServer).ClearInbox(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Notification_ClearInbox_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServer).ClearInbox(ctx, req.(*ClearInboxRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Notification_ServiceDesc is the grpc.ServiceDesc for Notification service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Notification_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.Notification",
	HandlerType: (*NotificationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListInboxNotifications",
			Handler:    _Notification_ListInboxNotifications_Handler,
		},
		{
			MethodName: "DeleteInboxNotifications",
			Handler:    _Notification_DeleteInboxNotifications_Handler,
		},
		{
			MethodName: "ClearInbox",
			Handler:    _Notification_ClearInbox_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "notification.proto",
}
