package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
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
	result, err := r.uc.OpenRoomInstance(ctx, request.RoomId, claims.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.GetRoomInstanceResponse{
		Code:    200,
		Message: "SUCCESS",
		Data: &v1.RoomInstanceDto{
			RoomInstanceId: result.RoomInstanceId,
			RoomId:         result.RoomId,
			ServerUrl:      result.ServerUrl,
		},
		AccessToken: result.AccessToken,
	}, nil
}

func (r *RoomInstanceService) ListGameHistory(ctx context.Context, request *v1.ListGameHistoryRequest) (*v1.ListGameHistoryResponse, error) {
	page := common.Pagination{
		Page:     request.Page,
		PageSize: request.PageSize,
	}
	instances, err := r.uc.ListRoomGameHistory(ctx, request.RoomId, &page)
	if err != nil {
		return nil, err
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
		Code:  200,
		Data:  result,
		Total: page.Total,
	}, nil
}
