package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/platform/api/v1"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/StellrisJAY/cloud-emu/util"
)

type RoomInstanceService struct {
	v1.UnimplementedRoomInstanceServer
	uc *biz.RoomInstanceUseCase
}

func NewRoomInstanceService(uc *biz.RoomInstanceUseCase) v1.RoomInstanceServer {
	return &RoomInstanceService{uc: uc}
}

func (r *RoomInstanceService) GetRoomInstance(ctx context.Context, request *v1.GetRoomInstanceRequest) (*v1.GetRoomInstanceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoomInstanceService) ListGameHistory(ctx context.Context, request *v1.ListGameHistoryRequest) (*v1.ListGameHistoryResponse, error) {
	page := util.Pagination{
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
