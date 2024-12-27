package data

import (
	"bytes"
	"context"
	"errors"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/biz"
	"net/url"
	"strconv"
	"strings"
)

type GameFileRepo struct {
	data *Data
}

func NewGameFileRepo(data *Data) biz.GameFileRepo {
	return &GameFileRepo{
		data: data,
	}
}

func (g *GameFileRepo) GetGameData(_ context.Context, game string) ([]byte, error) {
	u, _ := url.Parse(game)
	if u.Scheme != "mongodb" {
		return nil, errors.New("file system not implemented")
	}
	database := u.Host
	bucketName := strings.Split(u.Path, "/")[1]
	gameId := strings.Split(u.Path, "/")[2]
	bucket, err := g.data.getGridFSBucket(database, bucketName)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer([]byte{})
	id, _ := strconv.ParseInt(gameId, 10, 64)
	_, err = bucket.DownloadToStream(id, buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
