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

type GameSaveRepo struct {
	d *Data
}

type GameSaveEntity struct {
	SaveId     int64
	RoomId     int64
	EmulatorId int64
	GameId     int64
	AddTime    time.Time
	FileUrl    string
	SaveName   string
	Md5        string
}

const GameSaveTableName = "game_save"
const GameSaveBucketName = "game_save"

func NewGameSaveRepo(d *Data) biz.GameSaveRepo {
	return &GameSaveRepo{d: d}
}

func (g *GameSaveRepo) getBucket(dbName, bucketName string) (*gridfs.Bucket, error) {
	return gridfs.NewBucket(g.d.mongo.Database(dbName), options.GridFSBucket().SetName(bucketName))
}

func getGameSaveNameForGridFS(saveId int64) string {
	return fmt.Sprintf("%d", saveId)
}

func (g *GameSaveRepo) Create(ctx context.Context, save *biz.GameSave) error {
	return g.d.DB(ctx).Table(GameSaveTableName).Create(gameSaveBizToEntity(save)).WithContext(ctx).Error
}

func (g *GameSaveRepo) Upload(_ context.Context, save *biz.GameSave, data []byte) error {
	bucket, err := g.d.getGridFSBucket(MongoDBName, GameSaveBucketName)
	if err != nil {
		return err
	}
	return bucket.UploadFromStreamWithID(save.SaveId, getGameSaveNameForGridFS(save.SaveId), bytes.NewReader(data))
}

func (g *GameSaveRepo) Download(_ context.Context, save *biz.GameSave) ([]byte, error) {
	bucket, err := g.d.getGridFSBucket(MongoDBName, GameSaveBucketName)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer([]byte{})
	_, err = bucket.DownloadToStream(save.SaveId, buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil

}

func (g *GameSaveRepo) List(ctx context.Context, query biz.GameSaveQuery, p *common.Pagination) ([]*biz.GameSave, error) {
	var result []*biz.GameSave
	d := g.d.DB(ctx).Table(GameSaveTableName + " gs").Select("gs.*, emulator_name, game_name, room_name, e.emulator_type, e.emulator_code").
		Joins("INNER JOIN emulator e ON gs.emulator_id = e.emulator_id").
		Joins("INNER JOIN sys_room sr ON sr.room_id = gs.room_id").
		Joins("INNER JOIN emulator_game eg ON gs.game_id = eg.game_id").
		Order("gs.save_id DESC")
	if query.RoomId != 0 {
		d = d.Where("gs.room_id = ?", query.RoomId)
	}
	if query.EmulatorId != 0 {
		d = d.Where("gs.emulator_id =?", query.EmulatorId)
	}
	if query.GameId != 0 {
		d = d.Where("gs.game_id = ?", query.GameId)
	}
	if query.HostId != 0 {
		d = d.Where("sr.host_id = ?", query.HostId)
	}
	err := d.Scopes(common.WithPagination(p)).WithContext(ctx).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return result, nil
}

func (g *GameSaveRepo) Delete(ctx context.Context, saveId int64) error {
	return g.d.DB(ctx).Table(GameSaveTableName).
		Where("save_id = ?", saveId).
		Delete(&GameSaveEntity{}).
		WithContext(ctx).
		Error
}

func (g *GameSaveRepo) DeleteFile(_ context.Context, save *biz.GameSave) error {
	bucket, err := g.d.getGridFSBucket(MongoDBName, GameSaveBucketName)
	if err != nil {
		return err
	}
	return bucket.Delete(save.SaveId)
}

func (g *GameSaveRepo) Get(ctx context.Context, saveId int64) (*biz.GameSave, error) {
	var result *biz.GameSave
	err := g.d.DB(ctx).Table(GameSaveTableName).WithContext(ctx).Where("save_id =?", saveId).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return result, nil
}

func (g *GameSaveRepo) GetDetail(ctx context.Context, saveId int64) (*biz.GameSave, error) {
	var result *biz.GameSave
	err := g.d.DB(ctx).Table(GameSaveTableName+" gs").Select("gs.*, emulator_name, game_name, room_name, e.emulator_type, e.emulator_code").
		Joins("INNER JOIN emulator e ON gs.emulator_id = e.emulator_id").
		Joins("INNER JOIN sys_room sr ON sr.room_id = gs.room_id").
		Joins("INNER JOIN emulator_game eg ON gs.game_id = eg.game_id").
		Where("gs.save_id = ?", saveId).
		WithContext(ctx).
		Scan(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else {
		return result, err
	}
}

func (g *GameSaveRepo) Rename(ctx context.Context, saveId int64, saveName string) error {
	return g.d.DB(ctx).Table(GameSaveTableName).
		Where("save_id =?", saveId).
		Update("save_name", saveName).
		WithContext(ctx).
		Error
}

func (g *GameSaveRepo) Exist(ctx context.Context, roomId int64, md5 string) (bool, error) {
	var result int64
	err := g.d.DB(ctx).Table(GameSaveTableName).
		Where("room_id =?", roomId).
		Where("md5 =?", md5).
		WithContext(ctx).
		Count(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func gameSaveBizToEntity(save *biz.GameSave) *GameSaveEntity {
	return &GameSaveEntity{
		SaveId:     save.SaveId,
		RoomId:     save.RoomId,
		EmulatorId: save.EmulatorId,
		GameId:     save.GameId,
		AddTime:    save.AddTime,
		FileUrl:    save.FileUrl,
		SaveName:   save.SaveName,
		Md5:        save.Md5,
	}
}
