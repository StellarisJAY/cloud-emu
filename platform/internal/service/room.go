package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/platform/api/v1"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
)

type RoomService struct {
	v1.UnimplementedRoomServer
	ruc *biz.RoomUseCase
}

func NewRoomService(ruc *biz.RoomUseCase) v1.RoomServer {
	return &RoomService{ruc: ruc}
}

func (r *RoomService) ListMyRooms(ctx context.Context, request *v1.ListRoomRequest) (*v1.ListRoomResponse, error) {
	rooms, err := r.ruc.ListMyRooms(ctx, 1864580034254077952)
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
			AddTime:     room.AddTime.UnixMilli(),
		})
	}
	return &v1.ListRoomResponse{
		Rooms: result,
		Total: int32(len(result)),
	}, nil
}

func (r *RoomService) ListAllRooms(ctx context.Context, request *v1.ListRoomRequest) (*v1.ListRoomResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoomService) CreateRoom(ctx context.Context, request *v1.CreateRoomRequest) (*v1.CreateRoomResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoomService) GetRoom(ctx context.Context, request *v1.GetRoomRequest) (*v1.GetRoomResponse, error) {
	//TODO implement me
	panic("implement me")
}
