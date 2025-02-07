// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.0
// - protoc             v3.12.4
// source: macro.proto

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

const OperationMacroApplyMacro = "/v1.Macro/ApplyMacro"
const OperationMacroCreateMacro = "/v1.Macro/CreateMacro"
const OperationMacroDeleteMacro = "/v1.Macro/DeleteMacro"
const OperationMacroListMacros = "/v1.Macro/ListMacros"

type MacroHTTPServer interface {
	ApplyMacro(context.Context, *ApplyMacroRequest) (*ApplyMacroResponse, error)
	CreateMacro(context.Context, *CreateMacroRequest) (*CreateMacroResponse, error)
	DeleteMacro(context.Context, *DeleteMacroRequest) (*DeleteMacroResponse, error)
	ListMacros(context.Context, *ListMacrosRequest) (*ListMacrosResponse, error)
}

func RegisterMacroHTTPServer(s *http.Server, srv MacroHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/macros", _Macro_ListMacros0_HTTP_Handler(srv))
	r.POST("/api/v1/macros", _Macro_CreateMacro0_HTTP_Handler(srv))
	r.DELETE("/api/v1/macros/{macroId}", _Macro_DeleteMacro0_HTTP_Handler(srv))
	r.POST("/api/v1/macros/apply", _Macro_ApplyMacro0_HTTP_Handler(srv))
}

func _Macro_ListMacros0_HTTP_Handler(srv MacroHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListMacrosRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationMacroListMacros)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListMacros(ctx, req.(*ListMacrosRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListMacrosResponse)
		return ctx.Result(200, reply)
	}
}

func _Macro_CreateMacro0_HTTP_Handler(srv MacroHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateMacroRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationMacroCreateMacro)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateMacro(ctx, req.(*CreateMacroRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateMacroResponse)
		return ctx.Result(200, reply)
	}
}

func _Macro_DeleteMacro0_HTTP_Handler(srv MacroHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteMacroRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationMacroDeleteMacro)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteMacro(ctx, req.(*DeleteMacroRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteMacroResponse)
		return ctx.Result(200, reply)
	}
}

func _Macro_ApplyMacro0_HTTP_Handler(srv MacroHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ApplyMacroRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationMacroApplyMacro)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ApplyMacro(ctx, req.(*ApplyMacroRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ApplyMacroResponse)
		return ctx.Result(200, reply)
	}
}

type MacroHTTPClient interface {
	ApplyMacro(ctx context.Context, req *ApplyMacroRequest, opts ...http.CallOption) (rsp *ApplyMacroResponse, err error)
	CreateMacro(ctx context.Context, req *CreateMacroRequest, opts ...http.CallOption) (rsp *CreateMacroResponse, err error)
	DeleteMacro(ctx context.Context, req *DeleteMacroRequest, opts ...http.CallOption) (rsp *DeleteMacroResponse, err error)
	ListMacros(ctx context.Context, req *ListMacrosRequest, opts ...http.CallOption) (rsp *ListMacrosResponse, err error)
}

type MacroHTTPClientImpl struct {
	cc *http.Client
}

func NewMacroHTTPClient(client *http.Client) MacroHTTPClient {
	return &MacroHTTPClientImpl{client}
}

func (c *MacroHTTPClientImpl) ApplyMacro(ctx context.Context, in *ApplyMacroRequest, opts ...http.CallOption) (*ApplyMacroResponse, error) {
	var out ApplyMacroResponse
	pattern := "/api/v1/macros/apply"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationMacroApplyMacro))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *MacroHTTPClientImpl) CreateMacro(ctx context.Context, in *CreateMacroRequest, opts ...http.CallOption) (*CreateMacroResponse, error) {
	var out CreateMacroResponse
	pattern := "/api/v1/macros"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationMacroCreateMacro))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *MacroHTTPClientImpl) DeleteMacro(ctx context.Context, in *DeleteMacroRequest, opts ...http.CallOption) (*DeleteMacroResponse, error) {
	var out DeleteMacroResponse
	pattern := "/api/v1/macros/{macroId}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationMacroDeleteMacro))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *MacroHTTPClientImpl) ListMacros(ctx context.Context, in *ListMacrosRequest, opts ...http.CallOption) (*ListMacrosResponse, error) {
	var out ListMacrosResponse
	pattern := "/api/v1/macros"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationMacroListMacros))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
