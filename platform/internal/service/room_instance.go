package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
)

type RoomInstanceService struct {
	v1.UnimplementedRoomInstanceServer
	uc *biz.RoomInstanceUseCase
}

func NewRoomInstanceService(uc *biz.RoomInstanceUseCase) v1.RoomInstanceServer {
	return &RoomInstanceService{uc: uc}
}

func (r *RoomInstanceService) GetRoomInstance(ctx context.Context, request *v1.GetRoomInstanceRequest) (*v1.GetRoomInstanceResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	tr, _ := transport.FromServerContext(ctx)
	ua := tr.RequestHeader().Get("User-Agent")
	ip := tr.RequestHeader().Get("X-Remote-Ip")
	result, err := r.uc.OpenRoomInstance(ctx, request.RoomId, biz.RoomMemberAuth{UserId: claims.UserId, Ip: ip, AppId: ua})
	if err != nil {
		e := errors.FromError(err)
		return &v1.GetRoomInstanceResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.GetRoomInstanceResponse{
		Code:    200,
		Message: "操作成功",
		Data: &v1.GetRoomInstanceResult{
			RoomInstance: &v1.RoomInstanceDto{
				RoomInstanceId: result.RoomInstanceId,
				RoomId:         result.RoomId,
				ServerIp:       result.ServerIp,
			},
			AccessToken: result.AccessToken,
		},
	}, nil
}

func (r *RoomInstanceService) ListGameHistory(ctx context.Context, request *v1.ListGameHistoryRequest) (*v1.ListGameHistoryResponse, error) {
	page := common.Pagination{
		Page:     request.Page,
		PageSize: request.PageSize,
	}
	instances, err := r.uc.ListRoomGameHistory(ctx, request.RoomId, &page)
	if err != nil {
		e := errors.FromError(err)
		return &v1.ListGameHistoryResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	result := make([]*v1.RoomInstanceDto, 0, len(instances))
	for _, instance := range instances {
		result = append(result, &v1.RoomInstanceDto{
			RoomInstanceId: instance.RoomInstanceId,
			RoomId:         instance.RoomId,
			EmulatorId:     instance.EmulatorId,
			EmulatorName:   instance.EmulatorName,
		})
	}
	return &v1.ListGameHistoryResponse{
		Code:    200,
		Message: "查询成功",
		Data:    result,
		Total:   page.Total,
	}, nil
}

func (r *RoomInstanceService) OpenGameConnection(ctx context.Context, request *v1.OpenGameConnectionRequest) (*v1.OpenGameConnectionResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	tr, _ := transport.FromServerContext(ctx)
	ua := tr.RequestHeader().Get("User-Agent")
	ip := tr.RequestHeader().Get("X-Remote-Ip")
	sdpOffer, err := r.uc.OpenGameConnection(ctx, request.RoomId, request.Token, biz.RoomMemberAuth{
		UserId: claims.UserId,
		Ip:     ip,
		AppId:  ua,
	})
	if err != nil {
		e := errors.FromError(err)
		return &v1.OpenGameConnectionResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.OpenGameConnectionResponse{
		Code:    200,
		Message: "创建成功",
		Data:    &v1.OpenGameConnectionResult{SdpOffer: sdpOffer},
	}, nil
}

func (r *RoomInstanceService) SdpAnswer(ctx context.Context, request *v1.SdpAnswerRequest) (*v1.SdpAnswerResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	tr, _ := transport.FromServerContext(ctx)
	ua := tr.RequestHeader().Get("User-Agent")
	ip := tr.RequestHeader().Get("X-Remote-Ip")
	err := r.uc.SdpAnswer(ctx, request.RoomId, request.Token, biz.RoomMemberAuth{
		UserId: claims.UserId,
		Ip:     ip,
		AppId:  ua,
	}, request.SdpAnswer)
	if err != nil {
		e := errors.FromError(err)
		return &v1.SdpAnswerResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.SdpAnswerResponse{
		Code:    200,
		Message: "创建成功",
	}, nil
}

func (r *RoomInstanceService) AddIceCandidate(ctx context.Context, request *v1.AddIceCandidateRequest) (*v1.AddIceCandidateResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	tr, _ := transport.FromServerContext(ctx)
	ua := tr.RequestHeader().Get("User-Agent")
	ip := tr.RequestHeader().Get("X-Remote-Ip")
	err := r.uc.AddICECandidate(ctx, request.RoomId, request.Token, biz.RoomMemberAuth{
		UserId: claims.UserId,
		Ip:     ip,
		AppId:  ua,
	}, request.IceCandidate)
	if err != nil {
		e := errors.FromError(err)
		return &v1.AddIceCandidateResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.AddIceCandidateResponse{
		Code:    200,
		Message: "创建成功",
	}, nil
}

func (r *RoomInstanceService) GetServerIceCandidate(ctx context.Context, request *v1.GetServerIceCandidateRequest) (*v1.GetServerIceCandidateResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	tr, _ := transport.FromServerContext(ctx)
	ua := tr.RequestHeader().Get("User-Agent")
	ip := tr.RequestHeader().Get("X-Remote-Ip")
	candidates, err := r.uc.GetServerICECandidates(ctx, request.RoomId, request.Token, biz.RoomMemberAuth{
		UserId: claims.UserId,
		Ip:     ip,
		AppId:  ua,
	})
	if err != nil {
		e := errors.FromError(err)
		return &v1.GetServerIceCandidateResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.GetServerIceCandidateResponse{
		Code:    200,
		Message: "创建成功",
		Data:    candidates,
	}, nil
}
