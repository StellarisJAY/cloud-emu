package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/biz"
	"github.com/redis/go-redis/v9"
	"time"
)

type MemberAuthRepo struct {
	data *Data
}

const memberTokenExpiration = time.Minute * 1

func roomMemberAuthCacheKey(roomId int64, token string) string {
	return fmt.Sprintf("/cloud-emu/room-member/auth/%d/%s", roomId, token)
}

func (m *MemberAuthRepo) StoreAuthInfo(token string, roomId int64, info *biz.MemberAuthInfo) error {
	data, _ := json.Marshal(info)
	cmd := m.data.redis.Set(context.Background(), roomMemberAuthCacheKey(roomId, token), data, memberTokenExpiration)
	return cmd.Err()
}

func (m *MemberAuthRepo) GetAuthInfo(token string, roomId int64) (*biz.MemberAuthInfo, error) {
	cmd := m.data.redis.Get(context.Background(), roomMemberAuthCacheKey(roomId, token))
	if err := cmd.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	data := cmd.Val()
	result := &biz.MemberAuthInfo{}
	_ = json.Unmarshal([]byte(data), result)
	return result, nil
}

func NewMemberAuthRepo(data *Data) biz.MemberAuthRepo {
	return &MemberAuthRepo{data: data}
}
