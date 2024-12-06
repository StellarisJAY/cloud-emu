package biz

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/golang-jwt/jwt/v5"
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

type UserUseCase struct {
	repo        UserRepo
	ar          AuthRepo
	snowflakeId *snowflake.Node
}

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	GetById(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, userName string) (*User, error)
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

func NewUserUseCase(repo UserRepo, ar AuthRepo, snowflakeId *snowflake.Node) *UserUseCase {
	return &UserUseCase{repo: repo, snowflakeId: snowflakeId, ar: ar}
}

func (uc *UserUseCase) Register(ctx context.Context, user *User) error {
	hash := sha256.Sum256([]byte(user.Password))
	user.Password = hex.EncodeToString(hash[:])
	user.AddTime = time.Now().Local()
	user.UserId = uc.snowflakeId.Generate().Int64()
	err := uc.repo.Create(ctx, user)
	return err
}

func (uc *UserUseCase) GetById(ctx context.Context, id int64) (*User, error) {
	return uc.repo.GetById(ctx, id)
}

func (uc *UserUseCase) Login(ctx context.Context, userName, password string) (string, error) {
	user, err := uc.repo.GetByUsername(ctx, userName)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256([]byte(password))
	password = hex.EncodeToString(hash[:])
	if password != user.Password {
		return "", errors.New(401, "wrong password", "password error")
	}
	token, err := uc.ar.CreateToken(ctx, &LoginClaims{UserId: user.UserId})
	if err != nil {
		return "", err
	}
	return token, nil
}
