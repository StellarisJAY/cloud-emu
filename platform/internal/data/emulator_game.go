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
	GameId       int64
	EmulatorType string
	GameName     string
	Size         int32
	AddTime      time.Time
	CustomData   string
	AddUser      int64
	Url          string
	Md5          string
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
	d := e.data.DB(ctx).Table(EmulatorGameTableName).
		WithContext(ctx)
	if query.EmulatorType != "" {
		d = d.Where("emulator_type = ?", query.EmulatorType)
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

func (e *EmulatorGameRepo) GetByEmulatorTypeAndName(ctx context.Context, emulatorType string, name string) (*biz.EmulatorGame, error) {
	var result *biz.EmulatorGame
	err := e.data.DB(ctx).Table(EmulatorGameTableName).WithContext(ctx).
		Where("emulator_type = ?", emulatorType).
		Where("game_name = ?", name).Scan(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else {
		return result, err
	}
}

func (e *EmulatorGameRepo) CountSame(ctx context.Context, game *biz.EmulatorGame) (int, error) {
	var result int64
	err := e.data.DB(ctx).Table(EmulatorGameTableName).
		Where("emulator_type = ?", game.EmulatorType).
		Where("(game_name =? OR md5 =?)", game.GameName, game.Md5).
		WithContext(ctx).
		Count(&result).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return int(result), nil
}

func convertEmulatorGameBizToEntity(gameBiz *biz.EmulatorGame) *EmulatorGameEntity {
	return &EmulatorGameEntity{
		GameId:       gameBiz.GameId,
		EmulatorType: gameBiz.EmulatorType,
		GameName:     gameBiz.GameName,
		Size:         gameBiz.Size,
		AddTime:      gameBiz.AddTime,
		CustomData:   gameBiz.CustomData,
		AddUser:      gameBiz.AddUser,
		Url:          gameBiz.Url,
		Md5:          gameBiz.Md5,
	}
}
