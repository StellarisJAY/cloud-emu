package data

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"time"
)

type EmulatorGameEntity struct {
	GameId     int64
	EmulatorId int64
	GameName   string
	Size       int32
	AddTime    time.Time
	CustomData string
	AddUser    int64
	Url        string
}

type EmulatorGameRepo struct {
	data *Data
}

const EmulatorGameTableName = "emulator_game"
const EmulatorGameBucketName = "game_file"

func NewEmulatorGameRepo(data *Data) biz.EmulatorGameRepo {
	return &EmulatorGameRepo{data: data}
}

func getGameFileNameForGridFS(gameId int64) string {
	return fmt.Sprintf("%d", gameId)
}

func (e *EmulatorGameRepo) getGridFSBucket(dbName string, bucketName string) (*gridfs.Bucket, error) {
	bucket, err := gridfs.NewBucket(e.data.mongo.Database(dbName), options.GridFSBucket().SetName(bucketName))
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func (e *EmulatorGameRepo) Create(ctx context.Context, game *biz.EmulatorGame) error {
	return e.data.DB(ctx).Table(EmulatorGameTableName).WithContext(ctx).Create(convertEmulatorGameBizToEntity(game)).Error
}

func (e *EmulatorGameRepo) Upload(_ context.Context, game *biz.EmulatorGame, data []byte) error {
	bucket, err := e.getGridFSBucket(MongoDBName, EmulatorGameBucketName)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(data)
	return bucket.UploadFromStreamWithID(game.GameId, getGameFileNameForGridFS(game.GameId), reader)
}

func (e *EmulatorGameRepo) DeleteFile(ctx context.Context, game *biz.EmulatorGame) error {
	bucket, err := e.getGridFSBucket(MongoDBName, EmulatorGameBucketName)
	if err != nil {
		return err
	}
	return bucket.DeleteContext(ctx, game.GameId)
}

func (e *EmulatorGameRepo) ListGame(ctx context.Context, query biz.EmulatorGameQuery, p *common.Pagination) ([]*biz.EmulatorGame, error) {
	var result []*biz.EmulatorGame
	d := e.data.DB(ctx).Table(EmulatorGameTableName).WithContext(ctx)
	if query.EmulatorId != 0 {
		d = d.Where("emulator_id = ?", query.EmulatorId)
	}
	if query.GameName != "" {
		d = d.Where("game_name LIKE ?", "%"+query.GameName+"%")
	}
	if err := d.Scopes(common.WithPagination(p)).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (e *EmulatorGameRepo) GetById(ctx context.Context, gameId int64) (*biz.EmulatorGame, error) {
	var result *biz.EmulatorGame
	err := e.data.DB(ctx).Table(EmulatorGameTableName).WithContext(ctx).Where("game_id = ?", gameId).Scan(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else {
		return result, err
	}
}

func (e *EmulatorGameRepo) Delete(ctx context.Context, gameId int64) error {
	return e.data.DB(ctx).
		Table(EmulatorGameTableName).
		WithContext(ctx).
		Where("game_id = ?", gameId).
		Delete(&biz.EmulatorGame{}).Error
}

func (e *EmulatorGameRepo) Download(ctx context.Context, game *biz.EmulatorGame) ([]byte, error) {
	bucket, err := e.getGridFSBucket(MongoDBName, EmulatorGameBucketName)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer([]byte{})
	_, err = bucket.DownloadToStream(game.GameId, buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (e *EmulatorGameRepo) GetByEmulatorIdAndName(ctx context.Context, emulatorId int64, name string) (*biz.EmulatorGame, error) {
	var result *biz.EmulatorGame
	err := e.data.DB(ctx).Table(EmulatorGameTableName).WithContext(ctx).
		Where("emulator_id = ?", emulatorId).
		Where("game_name = ?", name).Scan(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else {
		return result, err
	}
}

func convertEmulatorGameBizToEntity(gameBiz *biz.EmulatorGame) *EmulatorGameEntity {
	return &EmulatorGameEntity{
		GameId:     gameBiz.GameId,
		EmulatorId: gameBiz.EmulatorId,
		GameName:   gameBiz.GameName,
		Size:       gameBiz.Size,
		AddTime:    gameBiz.AddTime,
		CustomData: gameBiz.CustomData,
		AddUser:    gameBiz.AddUser,
		Url:        gameBiz.Url,
	}
}
