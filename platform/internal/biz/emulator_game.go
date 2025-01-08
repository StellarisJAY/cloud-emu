package biz

import (
	"context"
	"fmt"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/conf"
	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"strconv"
	"strings"
	"time"
)

type EmulatorGame struct {
	GameId       int64
	EmulatorId   int64
	GameName     string
	Size         int32
	AddTime      time.Time
	CustomData   string
	AddUser      int64
	AddUserName  string
	Url          string
	EmulatorName string
	EmulatorType string
}

type EmulatorGameQuery struct {
	EmulatorId int64
	GameName   string
}

type EmulatorGameRepo interface {
	Create(ctx context.Context, game *EmulatorGame) error
	Upload(ctx context.Context, game *EmulatorGame, data []byte) error
	GetById(ctx context.Context, gameId int64) (*EmulatorGame, error)
	Delete(ctx context.Context, gameId int64) error
	DeleteFile(ctx context.Context, game *EmulatorGame) error
	ListGame(ctx context.Context, query EmulatorGameQuery, p *common.Pagination) ([]*EmulatorGame, error)
	Download(ctx context.Context, game *EmulatorGame) ([]byte, error)
	GetByEmulatorIdAndName(ctx context.Context, emulatorId int64, name string) (*EmulatorGame, error)
}

type EmulatorGameUseCase struct {
	emulatorGameRepo EmulatorGameRepo
	tm               Transaction
	snowflakeId      *snowflake.Node
	userRepo         UserRepo
	ac               *conf.Auth
	logger           *log.Helper
}

func NewEmulatorGameUseCase(emulatorGameRepo EmulatorGameRepo, tm Transaction, snowflakeId *snowflake.Node,
	userRepo UserRepo, ac *conf.Auth, logger log.Logger) *EmulatorGameUseCase {
	return &EmulatorGameUseCase{
		emulatorGameRepo: emulatorGameRepo,
		tm:               tm,
		snowflakeId:      snowflakeId,
		userRepo:         userRepo,
		ac:               ac,
		logger:           log.NewHelper(logger),
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

	if err := uc.upload(c, &EmulatorGame{GameName: gameName, EmulatorId: emulatorId}, data, claims.UserId); err != nil {
		er := errors.FromError(err)
		_ = c.JSON(200, struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Message: er.Message,
			Code:    int(er.Code),
		})
	} else {
		_ = c.JSON(200, struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Message: "上传成功",
			Code:    200,
		})
	}
	return nil
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
		game.Url = fmt.Sprintf("mongodb://cloud-emu/game_file/%d", game.GameId)
		err := uc.emulatorGameRepo.Create(ctx, game)
		if err != nil {
			uc.logger.Error("添加游戏错误:", err)
			return errors.New(500, "Database Error", "上传游戏失败")
		}
		if err := uc.emulatorGameRepo.Upload(ctx, game, data); err != nil {
			uc.logger.Error("上传mongodb错误:", err)
			return errors.New(500, "Upload Error", "上传游戏失败")
		}
		return nil
	})
}

func (uc *EmulatorGameUseCase) Delete(ctx context.Context, gameId int64) error {
	return uc.tm.Tx(ctx, func(ctx context.Context) error {
		game, err := uc.emulatorGameRepo.GetById(ctx, gameId)
		if err != nil {
			uc.logger.Error("删除游戏错误:", err)
			return v1.ErrorServiceError("删除失败")
		}
		if game == nil {
			return v1.ErrorNotFound("游戏不存在")
		}
		if err := uc.emulatorGameRepo.Delete(ctx, gameId); err != nil {
			uc.logger.Error("删除游戏错误:", err)
			return v1.ErrorServiceError("删除失败")
		}
		if err := uc.emulatorGameRepo.DeleteFile(ctx, game); err != nil {
			uc.logger.Error("删除游戏错误:", err)
			return v1.ErrorServiceError("删除失败")
		}
		return nil
	})
}
