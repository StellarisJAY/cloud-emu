package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
)

type RoomService struct {
	v1.UnimplementedRoomServer
	ruc *biz.RoomUseCase
}

func NewRoomService(ruc *biz.RoomUseCase) v1.RoomServer {
	return &RoomService{ruc: ruc}
}

func (r *RoomService) ListMyRooms(ctx context.Context, request *v1.ListRoomRequest) (*v1.ListRoomResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	query := biz.RoomQuery{}
	page := &common.Pagination{
		Page:     request.Page,
		PageSize: request.PageSize,
	}
	rooms, err := r.ruc.ListMyRooms(ctx, claims.UserId, query, page)
	if err != nil {
		return nil, err
	}
	result := make([]*v1.RoomDto, 0, len(rooms))
	for _, room := range rooms {
		result = append(result, &v1.RoomDto{
			RoomId:      room.RoomId,
			RoomName:    room.RoomName,
			HostId:      room.HostId,
			HostName:    room.HostName,
			MemberLimit: room.MemberLimit,
			AddTime:     room.AddTime.Format("2006-01-02 15:04:05"),
			EmulatorId:  room.EmulatorId,
		})
	}
	return &v1.ListRoomResponse{
		Rooms: result,
		Total: page.Total,
	}, nil
}

func (r *RoomService) ListAllRooms(ctx context.Context, request *v1.ListRoomRequest) (*v1.ListRoomResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoomService) CreateRoom(ctx context.Context, request *v1.CreateRoomRequest) (*v1.CreateRoomResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	room := &biz.Room{
		RoomName:    request.Name,
		Description: request.Description,
		HostId:      claims.UserId,
		MemberLimit: request.MemberLimit,
		Password:    request.Password,
		JoinType:    request.JoinType,
	}
	err := r.ruc.Create(ctx, room)
	if err != nil {
		return nil, err
	}
	return &v1.CreateRoomResponse{Id: room.RoomId}, nil
}

func (r *RoomService) GetRoom(ctx context.Context, request *v1.GetRoomRequest) (*v1.GetRoomResponse, error) {
	//TODO implement me
	panic("implement me")
}
