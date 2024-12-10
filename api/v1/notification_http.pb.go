// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.0
// - protoc             v3.12.4
// source: notification.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationNotificationListInboxNotifications = "/v1.Notification/ListInboxNotifications"

type NotificationHTTPServer interface {
	ListInboxNotifications(context.Context, *ListInboxNotificationRequest) (*ListInboxNotificationResponse, error)
}

func RegisterNotificationHTTPServer(s *http.Server, srv NotificationHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/inbox/notifications", _Notification_ListInboxNotifications0_HTTP_Handler(srv))
}

func _Notification_ListInboxNotifications0_HTTP_Handler(srv NotificationHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListInboxNotificationRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationNotificationListInboxNotifications)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListInboxNotifications(ctx, req.(*ListInboxNotificationRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListInboxNotificationResponse)
		return ctx.Result(200, reply)
	}
}

type NotificationHTTPClient interface {
	ListInboxNotifications(ctx context.Context, req *ListInboxNotificationRequest, opts ...http.CallOption) (rsp *ListInboxNotificationResponse, err error)
}

type NotificationHTTPClientImpl struct {
	cc *http.Client
}

func NewNotificationHTTPClient(client *http.Client) NotificationHTTPClient {
	return &NotificationHTTPClientImpl{client}
}

func (c *NotificationHTTPClientImpl) ListInboxNotifications(ctx context.Context, in *ListInboxNotificationRequest, opts ...http.CallOption) (*ListInboxNotificationResponse, error) {
	var out ListInboxNotificationResponse
	pattern := "/api/v1/inbox/notifications"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationNotificationListInboxNotifications))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
