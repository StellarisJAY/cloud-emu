// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.0
// - protoc             v3.12.4
// source: room_member.proto

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

const OperationRoomMemberGetUserRoomMember = "/v1.RoomMember/GetUserRoomMember"
const OperationRoomMemberInviteRoomMember = "/v1.RoomMember/InviteRoomMember"
const OperationRoomMemberListRoomMember = "/v1.RoomMember/ListRoomMember"

type RoomMemberHTTPServer interface {
	GetUserRoomMember(context.Context, *GetUserRoomMemberRequest) (*GetUserRoomMemberResponse, error)
	InviteRoomMember(context.Context, *InviteRoomMemberRequest) (*InviteRoomMemberResponse, error)
	ListRoomMember(context.Context, *ListRoomMemberRequest) (*ListRoomMemberResponse, error)
}

func RegisterRoomMemberHTTPServer(s *http.Server, srv RoomMemberHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/room-member", _RoomMember_ListRoomMember0_HTTP_Handler(srv))
	r.POST("/api/v1/room-member/invite", _RoomMember_InviteRoomMember0_HTTP_Handler(srv))
	r.GET("/api/v1/room-member/user", _RoomMember_GetUserRoomMember0_HTTP_Handler(srv))
}

func _RoomMember_ListRoomMember0_HTTP_Handler(srv RoomMemberHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListRoomMemberRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomMemberListRoomMember)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListRoomMember(ctx, req.(*ListRoomMemberRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListRoomMemberResponse)
		return ctx.Result(200, reply)
	}
}

func _RoomMember_InviteRoomMember0_HTTP_Handler(srv RoomMemberHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in InviteRoomMemberRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomMemberInviteRoomMember)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.InviteRoomMember(ctx, req.(*InviteRoomMemberRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*InviteRoomMemberResponse)
		return ctx.Result(200, reply)
	}
}

func _RoomMember_GetUserRoomMember0_HTTP_Handler(srv RoomMemberHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetUserRoomMemberRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomMemberGetUserRoomMember)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetUserRoomMember(ctx, req.(*GetUserRoomMemberRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetUserRoomMemberResponse)
		return ctx.Result(200, reply)
	}
}

type RoomMemberHTTPClient interface {
	GetUserRoomMember(ctx context.Context, req *GetUserRoomMemberRequest, opts ...http.CallOption) (rsp *GetUserRoomMemberResponse, err error)
	InviteRoomMember(ctx context.Context, req *InviteRoomMemberRequest, opts ...http.CallOption) (rsp *InviteRoomMemberResponse, err error)
	ListRoomMember(ctx context.Context, req *ListRoomMemberRequest, opts ...http.CallOption) (rsp *ListRoomMemberResponse, err error)
}

type RoomMemberHTTPClientImpl struct {
	cc *http.Client
}

func NewRoomMemberHTTPClient(client *http.Client) RoomMemberHTTPClient {
	return &RoomMemberHTTPClientImpl{client}
}

func (c *RoomMemberHTTPClientImpl) GetUserRoomMember(ctx context.Context, in *GetUserRoomMemberRequest, opts ...http.CallOption) (*GetUserRoomMemberResponse, error) {
	var out GetUserRoomMemberResponse
	pattern := "/api/v1/room-member/user"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoomMemberGetUserRoomMember))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomMemberHTTPClientImpl) InviteRoomMember(ctx context.Context, in *InviteRoomMemberRequest, opts ...http.CallOption) (*InviteRoomMemberResponse, error) {
	var out InviteRoomMemberResponse
	pattern := "/api/v1/room-member/invite"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomMemberInviteRoomMember))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomMemberHTTPClientImpl) ListRoomMember(ctx context.Context, in *ListRoomMemberRequest, opts ...http.CallOption) (*ListRoomMemberResponse, error) {
	var out ListRoomMemberResponse
	pattern := "/api/v1/room-member"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoomMemberListRoomMember))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
