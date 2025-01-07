// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.0
// - protoc             v3.12.4
// source: room_instance.proto

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

const OperationRoomInstanceAddIceCandidate = "/v1.RoomInstance/AddIceCandidate"
const OperationRoomInstanceGetControllerPlayers = "/v1.RoomInstance/GetControllerPlayers"
const OperationRoomInstanceGetRoomInstance = "/v1.RoomInstance/GetRoomInstance"
const OperationRoomInstanceGetServerIceCandidate = "/v1.RoomInstance/GetServerIceCandidate"
const OperationRoomInstanceListGameHistory = "/v1.RoomInstance/ListGameHistory"
const OperationRoomInstanceOpenGameConnection = "/v1.RoomInstance/OpenGameConnection"
const OperationRoomInstanceRestartRoomInstance = "/v1.RoomInstance/RestartRoomInstance"
const OperationRoomInstanceSdpAnswer = "/v1.RoomInstance/SdpAnswer"
const OperationRoomInstanceSetControllerPlayer = "/v1.RoomInstance/SetControllerPlayer"

type RoomInstanceHTTPServer interface {
	AddIceCandidate(context.Context, *AddIceCandidateRequest) (*AddIceCandidateResponse, error)
	GetControllerPlayers(context.Context, *GetControllerPlayersRequest) (*GetControllerPlayersResponse, error)
	GetRoomInstance(context.Context, *GetRoomInstanceRequest) (*GetRoomInstanceResponse, error)
	GetServerIceCandidate(context.Context, *GetServerIceCandidateRequest) (*GetServerIceCandidateResponse, error)
	ListGameHistory(context.Context, *ListGameHistoryRequest) (*ListGameHistoryResponse, error)
	OpenGameConnection(context.Context, *OpenGameConnectionRequest) (*OpenGameConnectionResponse, error)
	RestartRoomInstance(context.Context, *RestartRoomInstanceRequest) (*RestartRoomInstanceResponse, error)
	SdpAnswer(context.Context, *SdpAnswerRequest) (*SdpAnswerResponse, error)
	SetControllerPlayer(context.Context, *SetControllerPlayerRequest) (*SetControllerPlayerResponse, error)
}

func RegisterRoomInstanceHTTPServer(s *http.Server, srv RoomInstanceHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/room-instance", _RoomInstance_GetRoomInstance0_HTTP_Handler(srv))
	r.GET("/api/v1/room-instance/history", _RoomInstance_ListGameHistory0_HTTP_Handler(srv))
	r.POST("/api/v1/room-instance/connect", _RoomInstance_OpenGameConnection0_HTTP_Handler(srv))
	r.POST("/api/v1/room-instance/sdp-answer", _RoomInstance_SdpAnswer0_HTTP_Handler(srv))
	r.POST("/api/v1/room-instance/ice-candidate", _RoomInstance_AddIceCandidate0_HTTP_Handler(srv))
	r.GET("/api/v1/room-instance/ice-candidate", _RoomInstance_GetServerIceCandidate0_HTTP_Handler(srv))
	r.POST("/api/v1/room-instance/restart", _RoomInstance_RestartRoomInstance0_HTTP_Handler(srv))
	r.GET("/api/v1/room-instance/controller-players", _RoomInstance_GetControllerPlayers0_HTTP_Handler(srv))
	r.POST("/api/v1/room-instance/controller-players", _RoomInstance_SetControllerPlayer0_HTTP_Handler(srv))
}

func _RoomInstance_GetRoomInstance0_HTTP_Handler(srv RoomInstanceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetRoomInstanceRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomInstanceGetRoomInstance)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetRoomInstance(ctx, req.(*GetRoomInstanceRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetRoomInstanceResponse)
		return ctx.Result(200, reply)
	}
}

func _RoomInstance_ListGameHistory0_HTTP_Handler(srv RoomInstanceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListGameHistoryRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomInstanceListGameHistory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListGameHistory(ctx, req.(*ListGameHistoryRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListGameHistoryResponse)
		return ctx.Result(200, reply)
	}
}

func _RoomInstance_OpenGameConnection0_HTTP_Handler(srv RoomInstanceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in OpenGameConnectionRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomInstanceOpenGameConnection)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.OpenGameConnection(ctx, req.(*OpenGameConnectionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*OpenGameConnectionResponse)
		return ctx.Result(200, reply)
	}
}

func _RoomInstance_SdpAnswer0_HTTP_Handler(srv RoomInstanceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in SdpAnswerRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomInstanceSdpAnswer)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SdpAnswer(ctx, req.(*SdpAnswerRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SdpAnswerResponse)
		return ctx.Result(200, reply)
	}
}

func _RoomInstance_AddIceCandidate0_HTTP_Handler(srv RoomInstanceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AddIceCandidateRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomInstanceAddIceCandidate)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AddIceCandidate(ctx, req.(*AddIceCandidateRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AddIceCandidateResponse)
		return ctx.Result(200, reply)
	}
}

func _RoomInstance_GetServerIceCandidate0_HTTP_Handler(srv RoomInstanceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetServerIceCandidateRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomInstanceGetServerIceCandidate)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetServerIceCandidate(ctx, req.(*GetServerIceCandidateRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetServerIceCandidateResponse)
		return ctx.Result(200, reply)
	}
}

func _RoomInstance_RestartRoomInstance0_HTTP_Handler(srv RoomInstanceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in RestartRoomInstanceRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomInstanceRestartRoomInstance)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.RestartRoomInstance(ctx, req.(*RestartRoomInstanceRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RestartRoomInstanceResponse)
		return ctx.Result(200, reply)
	}
}

func _RoomInstance_GetControllerPlayers0_HTTP_Handler(srv RoomInstanceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetControllerPlayersRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomInstanceGetControllerPlayers)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetControllerPlayers(ctx, req.(*GetControllerPlayersRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetControllerPlayersResponse)
		return ctx.Result(200, reply)
	}
}

func _RoomInstance_SetControllerPlayer0_HTTP_Handler(srv RoomInstanceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in SetControllerPlayerRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomInstanceSetControllerPlayer)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SetControllerPlayer(ctx, req.(*SetControllerPlayerRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SetControllerPlayerResponse)
		return ctx.Result(200, reply)
	}
}

type RoomInstanceHTTPClient interface {
	AddIceCandidate(ctx context.Context, req *AddIceCandidateRequest, opts ...http.CallOption) (rsp *AddIceCandidateResponse, err error)
	GetControllerPlayers(ctx context.Context, req *GetControllerPlayersRequest, opts ...http.CallOption) (rsp *GetControllerPlayersResponse, err error)
	GetRoomInstance(ctx context.Context, req *GetRoomInstanceRequest, opts ...http.CallOption) (rsp *GetRoomInstanceResponse, err error)
	GetServerIceCandidate(ctx context.Context, req *GetServerIceCandidateRequest, opts ...http.CallOption) (rsp *GetServerIceCandidateResponse, err error)
	ListGameHistory(ctx context.Context, req *ListGameHistoryRequest, opts ...http.CallOption) (rsp *ListGameHistoryResponse, err error)
	OpenGameConnection(ctx context.Context, req *OpenGameConnectionRequest, opts ...http.CallOption) (rsp *OpenGameConnectionResponse, err error)
	RestartRoomInstance(ctx context.Context, req *RestartRoomInstanceRequest, opts ...http.CallOption) (rsp *RestartRoomInstanceResponse, err error)
	SdpAnswer(ctx context.Context, req *SdpAnswerRequest, opts ...http.CallOption) (rsp *SdpAnswerResponse, err error)
	SetControllerPlayer(ctx context.Context, req *SetControllerPlayerRequest, opts ...http.CallOption) (rsp *SetControllerPlayerResponse, err error)
}

type RoomInstanceHTTPClientImpl struct {
	cc *http.Client
}

func NewRoomInstanceHTTPClient(client *http.Client) RoomInstanceHTTPClient {
	return &RoomInstanceHTTPClientImpl{client}
}

func (c *RoomInstanceHTTPClientImpl) AddIceCandidate(ctx context.Context, in *AddIceCandidateRequest, opts ...http.CallOption) (*AddIceCandidateResponse, error) {
	var out AddIceCandidateResponse
	pattern := "/api/v1/room-instance/ice-candidate"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomInstanceAddIceCandidate))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomInstanceHTTPClientImpl) GetControllerPlayers(ctx context.Context, in *GetControllerPlayersRequest, opts ...http.CallOption) (*GetControllerPlayersResponse, error) {
	var out GetControllerPlayersResponse
	pattern := "/api/v1/room-instance/controller-players"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoomInstanceGetControllerPlayers))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomInstanceHTTPClientImpl) GetRoomInstance(ctx context.Context, in *GetRoomInstanceRequest, opts ...http.CallOption) (*GetRoomInstanceResponse, error) {
	var out GetRoomInstanceResponse
	pattern := "/api/v1/room-instance"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoomInstanceGetRoomInstance))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomInstanceHTTPClientImpl) GetServerIceCandidate(ctx context.Context, in *GetServerIceCandidateRequest, opts ...http.CallOption) (*GetServerIceCandidateResponse, error) {
	var out GetServerIceCandidateResponse
	pattern := "/api/v1/room-instance/ice-candidate"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoomInstanceGetServerIceCandidate))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomInstanceHTTPClientImpl) ListGameHistory(ctx context.Context, in *ListGameHistoryRequest, opts ...http.CallOption) (*ListGameHistoryResponse, error) {
	var out ListGameHistoryResponse
	pattern := "/api/v1/room-instance/history"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoomInstanceListGameHistory))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomInstanceHTTPClientImpl) OpenGameConnection(ctx context.Context, in *OpenGameConnectionRequest, opts ...http.CallOption) (*OpenGameConnectionResponse, error) {
	var out OpenGameConnectionResponse
	pattern := "/api/v1/room-instance/connect"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomInstanceOpenGameConnection))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomInstanceHTTPClientImpl) RestartRoomInstance(ctx context.Context, in *RestartRoomInstanceRequest, opts ...http.CallOption) (*RestartRoomInstanceResponse, error) {
	var out RestartRoomInstanceResponse
	pattern := "/api/v1/room-instance/restart"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomInstanceRestartRoomInstance))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomInstanceHTTPClientImpl) SdpAnswer(ctx context.Context, in *SdpAnswerRequest, opts ...http.CallOption) (*SdpAnswerResponse, error) {
	var out SdpAnswerResponse
	pattern := "/api/v1/room-instance/sdp-answer"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomInstanceSdpAnswer))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomInstanceHTTPClientImpl) SetControllerPlayer(ctx context.Context, in *SetControllerPlayerRequest, opts ...http.CallOption) (*SetControllerPlayerResponse, error) {
	var out SetControllerPlayerResponse
	pattern := "/api/v1/room-instance/controller-players"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomInstanceSetControllerPlayer))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
