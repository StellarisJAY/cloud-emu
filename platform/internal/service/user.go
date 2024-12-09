package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
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
		return nil, err
	}
	return &v1.RegisterResponse{
		Code:    200,
		Message: "register success",
	}, nil
}

func (u *UserService) Login(ctx context.Context, request *v1.LoginRequest) (*v1.LoginResponse, error) {
	token, err := u.uuc.Login(ctx, request.UserName, request.Password)
	if err != nil {
		return nil, err
	}
	return &v1.LoginResponse{
		Code:    200,
		Message: "login success",
		Token:   token,
	}, nil
}
