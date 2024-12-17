package data

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
)

type UserEmailVerifyRepo struct {
	data *Data
}

const UserEmailVerifyCollectionName = "user_email_verify"

func NewUserEmailVerifyRepo(data *Data) biz.UserEmailVerifyRepo {
	return &UserEmailVerifyRepo{
		data: data,
	}
}

func (u *UserEmailVerifyRepo) Create(ctx context.Context, user *biz.UserEmailVerify) error {
	return u.data.DB(ctx).Table(UserEmailVerifyCollectionName).Create(user).WithContext(ctx).Error
}

func (u *UserEmailVerifyRepo) GetByUserId(ctx context.Context, userId int64) (*biz.UserEmailVerify, error) {
	var result *biz.UserEmailVerify
	err := u.data.DB(ctx).Table(UserEmailVerifyCollectionName).Where("user_id = ?", userId).First(&result).WithContext(ctx).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *UserEmailVerifyRepo) Update(ctx context.Context, user *biz.UserEmailVerify) error {
	return u.data.DB(ctx).Table(UserEmailVerifyCollectionName).Where("user_id = ?", user.UserId).Updates(user).WithContext(ctx).Error
}
