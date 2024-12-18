package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"time"
)

type RoomMemberService struct {
	v1.UnimplementedRoomMemberServer
	roomMemberUC *biz.RoomMemberUseCase
}

func NewRoomMemberService(roomMemberUC *biz.RoomMemberUseCase) v1.RoomMemberServer {
	return &RoomMemberService{roomMemberUC: roomMemberUC}
}

func (r *RoomMemberService) ListRoomMember(ctx context.Context, request *v1.ListRoomMemberRequest) (*v1.ListRoomMemberResponse, error) {
	members, err := r.roomMemberUC.ListRoomMembers(ctx, request.RoomId)
	if err != nil {
		e := errors.FromError(err)
		return &v1.ListRoomMemberResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	result := make([]*v1.RoomMemberDto, 0, len(members))
	for _, member := range members {
		result = append(result, &v1.RoomMemberDto{
			RoomId:       member.RoomId,
			UserId:       member.UserId,
			RoomMemberId: member.RoomMemberId,
			UserName:     member.UserName,
			NickName:     member.NickName,
			Role:         member.Role,
			AddTime:      member.AddTime.Format(time.DateTime),
		})
	}
	return &v1.ListRoomMemberResponse{
		Code:    200,
		Message: "查询成功",
		Data:    result,
		Total:   int32(len(result)),
	}, nil
}

func (r *RoomMemberService) InviteRoomMember(ctx context.Context, request *v1.InviteRoomMemberRequest) (*v1.InviteRoomMemberResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := r.roomMemberUC.InviteRoomMember(ctx, claims.UserId, request.UserId, request.RoomId)
	if err != nil {
		e := errors.FromError(err)
		return &v1.InviteRoomMemberResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.InviteRoomMemberResponse{
		Code:    200,
		Message: "操作成功",
	}, nil
}

func (r *RoomMemberService) GetUserRoomMember(ctx context.Context, request *v1.GetUserRoomMemberRequest) (*v1.GetUserRoomMemberResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	rm, err := r.roomMemberUC.GetByRoomAndUser(ctx, request.RoomId, claims.UserId)
	if err != nil {
		e := errors.FromError(err)
		return &v1.GetUserRoomMemberResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.GetUserRoomMemberResponse{
		Code:    200,
		Message: "查询成功",
		Data: &v1.UserRoomMember{
			RoomMemberId: rm.RoomMemberId,
			RoomId:       rm.RoomId,
			UserId:       rm.UserId,
			Role:         rm.Role,
			AddTime:      rm.AddTime.Format(time.DateTime),
		},
	}, nil
}
