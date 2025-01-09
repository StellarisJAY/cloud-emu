package biz

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/util"
	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type User struct {
	UserId   int64
	UserName string
	NickName string
	Password string
	Phone    string
	Email    string
	Status   int32
	Role     int32
	AddTime  time.Time
}

type UserEmailVerify struct {
	Id      int64
	UserId  int64
	Code    string
	AddTime time.Time
}

const (
	UserStatusAvailable int32 = iota + 1
	UserStatusBanned
	UserStatusNotActivated
)

const (
	UserRoleNormal int32 = iota + 1
	UserRoleAdmin
)

type UserUseCase struct {
	repo                UserRepo
	ar                  AuthRepo
	snowflakeId         *snowflake.Node
	userEmailVerifyRepo UserEmailVerifyRepo
	eh                  *util.EmailHelper
	tm                  Transaction
	logger              *log.Helper
}

type UserQuery struct {
	UserId   int64
	UserName string
	NickName string
	Status   int32
}

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	GetById(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, userName string) (*User, error)
	ListUser(ctx context.Context, query UserQuery, p *common.Pagination) ([]*User, error)
}

type UserEmailVerifyRepo interface {
	Create(ctx context.Context, user *UserEmailVerify) error
	GetByUserId(ctx context.Context, userId int64) (*UserEmailVerify, error)
	Update(ctx context.Context, user *UserEmailVerify) error
}

type AuthRepo interface {
	CreateToken(ctx context.Context, claims *LoginClaims) (string, error)
}

type LoginClaims struct {
	jwt.RegisteredClaims
	UserId  int64
	AppId   string
	LoginIp string
}

func NewUserUseCase(repo UserRepo, ar AuthRepo, snowflakeId *snowflake.Node, userEmailVerifyRepo UserEmailVerifyRepo,
	eh *util.EmailHelper, tm Transaction, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo:                repo,
		snowflakeId:         snowflakeId,
		ar:                  ar,
		userEmailVerifyRepo: userEmailVerifyRepo,
		tm:                  tm,
		logger:              log.NewHelper(logger),
		eh:                  eh,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, user *User) error {
	hash := sha256.Sum256([]byte(user.Password))
	user.Password = hex.EncodeToString(hash[:])
	user.AddTime = time.Now().Local()
	user.UserId = uc.snowflakeId.Generate().Int64()
	user.Role = UserRoleNormal
	user.Status = UserStatusNotActivated
	return uc.tm.Tx(ctx, func(ctx context.Context) error {
		if err := uc.repo.Create(ctx, user); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errors.New(400, "Register Failed", "用户名已被使用")
			}
			uc.logger.Errorf("注册用户错误: %v", err)
			return v1.ErrorServiceError("注册用户出错")
		}
		verify := UserEmailVerify{
			Id:      uc.snowflakeId.Generate().Int64(),
			UserId:  user.UserId,
			AddTime: time.Now().Local(),
			Code:    uc.newVerifyCode(),
		}
		// 发送邮箱验证码
		if err := uc.eh.Send(user.Email, "cloudemu", "您的注册验证码为:"+verify.Code); err != nil {
			uc.logger.Error("注册用户发送邮件错误:", err)
			return v1.ErrorServiceError("注册用户出错")
		}
		if err := uc.userEmailVerifyRepo.Create(ctx, &verify); err != nil {
			uc.logger.Error("注册用户错误:", err)
			return v1.ErrorServiceError("注册用户出错")
		}
		return nil
	})
}

func (uc *UserUseCase) GetById(ctx context.Context, id int64) (*User, error) {
	user, err := uc.repo.GetById(ctx, id)
	if err != nil {
		uc.logger.Error("获取用户详情错误:", err)
		return nil, errors.New(500, "Database Error", "无法获取用户详情")
	}
	if user == nil {
		return nil, v1.ErrorNotFound("用户不存在")
	}
	return user, nil
}

func (uc *UserUseCase) Login(ctx context.Context, userName, password string) (string, error) {
	user, _ := uc.repo.GetByUsername(ctx, userName)
	if user == nil {
		return "", v1.ErrorLoginFailed("登录失败，请检查用户名和密码")
	}
	hash := sha256.Sum256([]byte(password))
	password = hex.EncodeToString(hash[:])
	if password != user.Password {
		return "", v1.ErrorLoginFailed("登录失败，请检查用户名和密码")
	}
	token, err := uc.ar.CreateToken(ctx, &LoginClaims{UserId: user.UserId})
	if err != nil {
		return "", v1.ErrorLoginFailed("登录失败，请检查用户名和密码")
	}
	return token, nil
}

func (uc *UserUseCase) newVerifyCode() string {
	number := rand.Intn(9000) + 1000
	return fmt.Sprintf("%0d", number)
}

func (uc *UserUseCase) ActivateAccount(ctx context.Context, userId int64, code string) error {
	return uc.tm.Tx(ctx, func(ctx context.Context) error {
		user, _ := uc.repo.GetById(ctx, userId)
		if user == nil {
			return v1.ErrorNotFound("账号不存在")
		}
		if user.Status != UserStatusNotActivated {
			return nil
		}
		verify, _ := uc.userEmailVerifyRepo.GetByUserId(ctx, userId)
		if verify == nil {
			return v1.ErrorNotFound("账号不存在")
		}
		if verify.Code == code {
			user.Status = UserStatusAvailable
			return uc.repo.Update(ctx, user)
		}
		return v1.ErrorActivationFailed("激活失败，请检查激活码")
	})
}

func (uc *UserUseCase) ListUser(ctx context.Context, query UserQuery, p *common.Pagination) ([]*User, error) {
	return uc.repo.ListUser(ctx, query, p)
}
