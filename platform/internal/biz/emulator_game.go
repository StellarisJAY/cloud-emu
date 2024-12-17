package biz

import (
	"context"
	"fmt"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/conf"
	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"strconv"
	"strings"
	"time"
)

type EmulatorGame struct {
	GameId      int64
	EmulatorId  int64
	GameName    string
	Size        int32
	AddTime     time.Time
	CustomData  string
	AddUser     int64
	AddUserName string
	Url         string
}

type EmulatorGameQuery struct {
	EmulatorId int64
	GameName   string
}

type EmulatorGameRepo interface {
	Create(ctx context.Context, game *EmulatorGame) error
	Upload(ctx context.Context, game *EmulatorGame, data []byte) error
	Delete(ctx context.Context, gameId int64) error
	ListGame(ctx context.Context, query EmulatorGameQuery, p *common.Pagination) ([]*EmulatorGame, error)
}

type EmulatorGameUseCase struct {
	emulatorGameRepo EmulatorGameRepo
	tm               Transaction
	snowflakeId      *snowflake.Node
	userRepo         UserRepo
	ac               *conf.Auth
}

func NewEmulatorGameUseCase(emulatorGameRepo EmulatorGameRepo, tm Transaction, snowflakeId *snowflake.Node,
	userRepo UserRepo, ac *conf.Auth) *EmulatorGameUseCase {
	return &EmulatorGameUseCase{
		emulatorGameRepo: emulatorGameRepo,
		tm:               tm,
		snowflakeId:      snowflakeId,
		userRepo:         userRepo,
		ac:               ac,
	}
}

func (uc *EmulatorGameUseCase) ListGame(ctx context.Context, query EmulatorGameQuery, p *common.Pagination) ([]*EmulatorGame, error) {
	return uc.emulatorGameRepo.ListGame(ctx, query, p)
}

func (uc *EmulatorGameUseCase) Upload(c http.Context) error {
	request := c.Request()
	file, _, err := request.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	gameName := request.FormValue("gameName")
	emulatorId, _ := strconv.ParseInt(request.FormValue("emulatorId"), 10, 64)
	token := request.Header.Get("Authorization")
	if token == "" {
		return errors.New(401, "Unauthorized", "没有登录")
	}
	token = strings.TrimPrefix(token, "Bearer ")
	claims := &LoginClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(uc.ac.JwtSecret), nil
	})
	if err != nil {
		return errors.New(401, "Unauthorized", "token无效")
	}

	return uc.upload(c, &EmulatorGame{GameName: gameName, EmulatorId: emulatorId}, data, claims.UserId)
}

func (uc *EmulatorGameUseCase) upload(ctx context.Context, game *EmulatorGame, data []byte, userId int64) error {
	user, _ := uc.userRepo.GetById(ctx, userId)
	if user == nil || user.Role != UserRoleAdmin || user.Status != UserStatusAvailable {
		return v1.ErrorAccessDenied("没有上传游戏权限")
	}
	return uc.tm.Tx(ctx, func(ctx context.Context) error {
		game.GameId = uc.snowflakeId.Generate().Int64()
		game.Size = int32(len(data))
		game.AddTime = time.Now()
		game.AddUser = userId
		game.Url = fmt.Sprintf("mongodb://game-file/%d/%d", game.EmulatorId, game.GameId)
		err := uc.emulatorGameRepo.Create(ctx, game)
		if err != nil {
			return errors.New(500, "Database Error", "上传游戏失败")
		}
		if err := uc.emulatorGameRepo.Upload(ctx, game, data); err != nil {
			return errors.New(500, "Upload Error", "上传游戏失败")
		}
		return nil
	})
}
