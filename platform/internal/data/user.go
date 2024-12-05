package data

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"time"
)

type UserRepo struct {
	data *Data
}

type UserModel struct {
	UserId   int64 `gorm:"primary_key"`
	UserName string
	NickName string
	Password string
	Phone    string
	Email    string
	Status   int32
	Role     int32
	AddTime  time.Time `gorm:"type:datetime"`
}

const UserTableName = "sys_user"

func NewUserRepo(data *Data) biz.UserRepo {
	return &UserRepo{data: data}
}

func (u *UserRepo) Create(ctx context.Context, user *biz.User) error {
	return u.data.db.Table(UserTableName).Create(user).Error
}

func (u *UserRepo) GetById(ctx context.Context, id int64) (*biz.User, error) {
	model := &UserModel{}
	err := u.data.db.Table(UserTableName).Where("user_id = ?", id).WithContext(ctx).First(model).Error
	if err != nil {
		return nil, err
	}
	return convertModelToBiz(model), nil
}

func (u *UserRepo) Update(ctx context.Context, user *biz.User) error {
	model := convertBizToModel(user)
	err := u.data.db.Table(UserTableName).Where("user_id = ?", user.UserId).Updates(model).Error
	return err
}

func (u *UserRepo) GetByUsername(ctx context.Context, userName string) (*biz.User, error) {
	model := &UserModel{}
	err := u.data.db.Table(UserTableName).Where("user_name = ?", userName).WithContext(ctx).First(model).Error
	if err != nil {
		return nil, err
	}
	return convertModelToBiz(model), nil
}

func convertModelToBiz(model *UserModel) *biz.User {
	return &biz.User{
		UserId:   model.UserId,
		UserName: model.UserName,
		NickName: model.NickName,
		Password: model.Password,
		Phone:    model.Phone,
		Email:    model.Email,
		Status:   model.Status,
		Role:     model.Role,
	}
}

func convertBizToModel(model *biz.User) *UserModel {
	return &UserModel{
		UserId:   model.UserId,
		UserName: model.UserName,
		NickName: model.NickName,
		Password: model.Password,
		Phone:    model.Phone,
		Email:    model.Email,
		Status:   model.Status,
		Role:     model.Role,
	}
}
