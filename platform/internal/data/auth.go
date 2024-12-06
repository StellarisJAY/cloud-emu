package data

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/StellrisJAY/cloud-emu/platform/internal/conf"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AuthRepo struct {
	cfg *conf.Auth
}

func NewAuthRepo(cfg *conf.Auth) biz.AuthRepo {
	return &AuthRepo{cfg: cfg}
}

func (a *AuthRepo) CreateToken(ctx context.Context, claims *biz.LoginClaims) (string, error) {
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	claims.Issuer = "cloudenmu-platform"
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	c := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := c.SignedString([]byte(a.cfg.JwtSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}
