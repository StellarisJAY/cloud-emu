package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"time"
)

type RoomService struct {
	v1.UnimplementedRoomServer
	roomUC       *biz.RoomUseCase
	roomMemberUC *biz.RoomMemberUseCase
}

func NewRoomService(roomUC *biz.RoomUseCase, roomMemberUC *biz.RoomMemberUseCase) v1.RoomServer {
	return &RoomService{roomUC: roomUC, roomMemberUC: roomMemberUC}
}

func (r *RoomService) ListMyRooms(ctx context.Context, request *v1.ListRoomRequest) (*v1.ListRoomResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	query := biz.RoomQuery{
		RoomName:   request.RoomName,
		HostName:   request.HostName,
		JoinType:   request.JoinType,
		EmulatorId: request.EmulatorId,
		HostOnly:   request.HostOnly,
	}
	page := &common.Pagination{
		Page:     request.Page,
		PageSize: request.PageSize,
	}
	rooms, err := r.roomUC.ListMyRooms(ctx, claims.UserId, query, page)
	if err != nil {
		e := errors.FromError(err)
		return &v1.ListRoomResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	result := make([]*v1.RoomDto, 0, len(rooms))
	for _, room := range rooms {
		result = append(result, &v1.RoomDto{
			RoomId:      room.RoomId,
			RoomName:    room.RoomName,
			HostId:      room.HostId,
			HostName:    room.HostName,
			MemberCount: room.MemberCount,
			MemberLimit: room.MemberLimit,
			AddTime:     room.AddTime.Format(time.DateTime),
			EmulatorId:  room.EmulatorId,
			JoinType:    room.JoinType,
			IsHost:      room.IsHost,
		})
	}
	return &v1.ListRoomResponse{
		Code:    200,
		Message: "查询成功",
		Data:    result,
		Total:   page.Total,
	}, nil
}

func (r *RoomService) ListAllRooms(ctx context.Context, request *v1.ListRoomRequest) (*v1.ListRoomResponse, error) {
	query := biz.RoomQuery{
		RoomName:   request.RoomName,
		HostName:   request.HostName,
		JoinType:   request.JoinType,
		EmulatorId: request.EmulatorId,
	}
	page := &common.Pagination{
		Page:     request.Page,
		PageSize: request.PageSize,
	}
	rooms, err := r.roomUC.ListAllRooms(ctx, query, page)
	if err != nil {
		e := errors.FromError(err)
		return &v1.ListRoomResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	result := make([]*v1.RoomDto, 0, len(rooms))
	for _, room := range rooms {
		result = append(result, &v1.RoomDto{
			RoomId:      room.RoomId,
			RoomName:    room.RoomName,
			HostId:      room.HostId,
			HostName:    room.HostName,
			MemberCount: room.MemberCount,
			MemberLimit: room.MemberLimit,
			AddTime:     room.AddTime.Format(time.DateTime),
			EmulatorId:  room.EmulatorId,
			JoinType:    room.JoinType,
		})
	}
	return &v1.ListRoomResponse{
		Code:    200,
		Message: "查询成功",
		Data:    result,
		Total:   page.Total,
	}, nil
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
	err := r.roomUC.Create(ctx, room)
	if err != nil {
		e := errors.FromError(err)
		return &v1.CreateRoomResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.CreateRoomResponse{
		Code:    200,
		Message: "创建成功",
	}, nil
}

func (r *RoomService) GetRoom(ctx context.Context, request *v1.GetRoomRequest) (*v1.GetRoomResponse, error) {
	room, err := r.roomUC.GetById(ctx, request.Id)
	if err != nil {
		e := errors.FromError(err)
		return &v1.GetRoomResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.GetRoomResponse{
		Code:    200,
		Message: "查询成功",
		Data: &v1.RoomDto{
			RoomId:       room.RoomId,
			RoomName:     room.RoomName,
			HostId:       room.HostId,
			HostName:     room.HostName,
			JoinType:     room.JoinType,
			MemberCount:  room.MemberCount,
			MemberLimit:  room.MemberLimit,
			AddTime:      room.AddTime.Format(time.DateTime),
			EmulatorId:   room.EmulatorId,
			EmulatorName: room.EmulatorName,
			GameId:       room.GameId,
			GameName:     room.GameName,
			EmulatorType: room.EmulatorType,
		},
	}, nil
}

func (r *RoomService) JoinRoom(ctx context.Context, request *v1.JoinRoomRequest) (*v1.JoinRoomResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := r.roomUC.Join(ctx, request.RoomId, claims.UserId, request.Password)
	if err != nil {
		e := errors.FromError(err)
		return &v1.JoinRoomResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.JoinRoomResponse{
		Code:    200,
		Message: "加入成功",
	}, nil
}

func (r *RoomService) DeleteRoom(ctx context.Context, request *v1.DeleteRoomRequest) (*v1.DeleteRoomResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := r.roomUC.Delete(ctx, request.RoomId, claims.UserId)
	if err != nil {
		e := errors.FromError(err)
		return &v1.DeleteRoomResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.DeleteRoomResponse{
		Code:    200,
		Message: "删除成功",
	}, nil
}

func (r *RoomService) UpdateRoom(ctx context.Context, request *v1.UpdateRoomRequest) (*v1.UpdateRoomResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	room := &biz.Room{
		RoomId:   request.RoomId,
		RoomName: request.RoomName,
		Password: request.Password,
		JoinType: request.JoinType,
	}
	err := r.roomUC.UpdateRoom(ctx, room, claims.UserId)
	if err != nil {
		e := errors.FromError(err)
		return &v1.UpdateRoomResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.UpdateRoomResponse{
		Code:    200,
		Message: "修改成功",
	}, nil
}
