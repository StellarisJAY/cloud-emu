package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
	"time"
)

type UserService struct {
	v1.UnimplementedUserServer
	uuc *biz.UserUseCase
}

func NewUserService(uuc *biz.UserUseCase) v1.UserServer {
	return &UserService{uuc: uuc}
}

func (u *UserService) Register(ctx context.Context, request *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	err := u.uuc.Register(ctx, &biz.User{
		UserName: request.UserName,
		Password: request.Password,
		Phone:    request.Phone,
		Email:    request.Email,
		NickName: request.NickName,
	})
	if err != nil {
		e := errors.FromError(err)
		return &v1.RegisterResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.RegisterResponse{
		Code:    200,
		Message: "注册成功",
	}, nil
}

func (u *UserService) Login(ctx context.Context, request *v1.LoginRequest) (*v1.LoginResponse, error) {
	token, err := u.uuc.Login(ctx, request.UserName, request.Password)
	if err != nil {
		e := errors.FromError(err)
		return &v1.LoginResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	if tr, ok := transport.FromServerContext(ctx); ok {
		tr.ReplyHeader().Set("Authorization", token)
	}

	return &v1.LoginResponse{
		Code:    200,
		Message: "登录成功",
	}, nil
}

func (u *UserService) ActivateAccount(ctx context.Context, request *v1.ActivateAccountRequest) (*v1.ActivateAccountResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := u.uuc.ActivateAccount(ctx, claims.UserId, request.Code)
	if err != nil {
		e := errors.FromError(err)
		return &v1.ActivateAccountResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.ActivateAccountResponse{
		Code:    200,
		Message: "激活成功",
	}, nil
}

func (u *UserService) ListUser(ctx context.Context, request *v1.ListUserRequest) (*v1.ListUserResponse, error) {
	p := &common.Pagination{Page: request.Page, PageSize: request.PageSize}
	users, err := u.uuc.ListUser(ctx, biz.UserQuery{
		UserName: request.UserName,
		NickName: request.NickName,
		Status:   request.Status,
	}, p)
	if err != nil {
		e := errors.FromError(err)
		return &v1.ListUserResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	userList := make([]*v1.UserDto, len(users))
	for i, u := range users {
		userList[i] = &v1.UserDto{
			UserId:   u.UserId,
			NickName: u.NickName,
			UserName: u.UserName,
			Status:   u.Status,
			AddTime:  u.AddTime.Format(time.DateTime),
		}
	}
	return &v1.ListUserResponse{
		Code:    200,
		Message: "查询成功",
		Data:    userList,
		Total:   p.Total,
	}, nil
}

func (u *UserService) GetUserDetail(ctx context.Context, request *v1.GetUserDetailRequest) (*v1.GetUserDetailResponse, error) {
	user, err := u.uuc.GetById(ctx, request.UserId)
	if err != nil {
		e := errors.FromError(err)
		return &v1.GetUserDetailResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.GetUserDetailResponse{
		Code:    200,
		Message: "查询成功",
		Data: &v1.UserDetailDto{
			UserId:   user.UserId,
			UserName: user.UserName,
			NickName: user.NickName,
			AddTime:  user.AddTime.Format(time.DateTime),
			Status:   user.Status,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
		},
	}, nil
}

func (u *UserService) GetLoginUserDetail(ctx context.Context, _ *v1.GetLoginUserDetailRequest) (*v1.GetLoginUserDetailResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	user, err := u.uuc.GetById(ctx, claims.UserId)
	if err != nil {
		e := errors.FromError(err)
		return &v1.GetLoginUserDetailResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.GetLoginUserDetailResponse{
		Code:    200,
		Message: "查询成功",
		Data: &v1.UserDetailDto{
			UserId:   user.UserId,
			UserName: user.UserName,
			NickName: user.NickName,
			AddTime:  user.AddTime.Format(time.DateTime),
			Status:   user.Status,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
		},
	}, nil
}
