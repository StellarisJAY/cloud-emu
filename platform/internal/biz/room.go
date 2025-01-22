package biz

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redsync/redsync/v4"
	"time"
)

type Room struct {
	RoomId       int64     `json:"roomId"`
	RoomName     string    `json:"roomName"`
	Description  string    `json:"description"`
	HostId       int64     `json:"hostId"`
	HostName     string    `json:"hostName"`
	MemberCount  int32     `json:"memberCount"`
	MemberLimit  int32     `json:"memberLimit"`
	AddTime      time.Time `json:"addTime"`
	Password     string    `json:"password"`
	JoinType     int32     `json:"joinRule"`
	EmulatorId   int64     `json:"emulatorId"`
	EmulatorName string    `json:"emulatorName"`
	GameId       int64     `json:"gameId"`
	GameName     string    `json:"gameName"`
	EmulatorType string    `json:"emulatorType"`
	EmulatorCode string    `json:"emulatorCode"`
	IsHost       bool      `json:"isHost"`
}

type RoomUseCase struct {
	repo             RoomRepo
	snowflakeId      *snowflake.Node
	userRepo         UserRepo
	tm               Transaction
	roomMemberRepo   RoomMemberRepo
	redSync          *redsync.Redsync
	logger           *log.Helper
	roomInstanceRepo RoomInstanceRepo
	gameSaveRepo     GameSaveRepo
	gameServerRepo   GameServerRepo
}

type RoomQuery struct {
	HostId     int64
	RoomName   string
	HostName   string
	JoinType   int32
	MemberId   int64
	EmulatorId int64
	HostOnly   bool
}

const (
	RoomJoinTypePublic int32 = iota + 1
	RoomJoinTypePassword
)

type RoomRepo interface {
	Create(ctx context.Context, room *Room) error
	GetById(ctx context.Context, id int64) (*Room, error)
	Update(ctx context.Context, room *Room) error
	ListJoinedRooms(ctx context.Context, query RoomQuery, page *common.Pagination) ([]*Room, error)
	ListRooms(ctx context.Context, query RoomQuery, page *common.Pagination) ([]*Room, error)
	Delete(ctx context.Context, id int64) error
}

func NewRoomUseCase(repo RoomRepo, snowflakeId *snowflake.Node, userRepo UserRepo, tm Transaction,
	roomMemberRepo RoomMemberRepo, redsync *redsync.Redsync, logger log.Logger, roomInstanceRepo RoomInstanceRepo,
	gameSaveRepo GameSaveRepo, gameServerRepo GameServerRepo) *RoomUseCase {
	return &RoomUseCase{repo: repo, snowflakeId: snowflakeId, userRepo: userRepo, tm: tm, roomMemberRepo: roomMemberRepo,
		redSync: redsync, logger: log.NewHelper(logger), roomInstanceRepo: roomInstanceRepo, gameSaveRepo: gameSaveRepo,
		gameServerRepo: gameServerRepo}
}

func (r *RoomUseCase) Create(ctx context.Context, room *Room) error {
	return r.tm.Tx(ctx, func(ctx context.Context) error {
		user, _ := r.userRepo.GetById(ctx, room.HostId)
		if user == nil || user.Status != UserStatusAvailable {
			return v1.ErrorAccessDenied("当前用户无法创建房间")
		}
		room.RoomId = r.snowflakeId.Generate().Int64()
		if room.Password != "" {
			hash := sha256.Sum256([]byte(room.Password))
			room.Password = hex.EncodeToString(hash[:])
		}
		room.AddTime = time.Now().Local()
		if err := r.repo.Create(ctx, room); err != nil {
			return err
		}
		return r.roomMemberRepo.Create(ctx, &RoomMember{
			RoomMemberId: r.snowflakeId.Generate().Int64(),
			RoomId:       room.RoomId,
			UserId:       room.HostId,
			Role:         RoomMemberRoleHost,
			AddTime:      room.AddTime,
			Status:       RoomMemberStatusJoined,
		})
	})
}

func (r *RoomUseCase) ListMyRooms(ctx context.Context, userId int64, query RoomQuery, page *common.Pagination) ([]*Room, error) {
	query.MemberId = userId
	rooms, err := r.repo.ListJoinedRooms(ctx, query, page)
	if err != nil {
		return nil, err
	}
	for _, room := range rooms {
		if err := r.buildRoomDto(ctx, room); err != nil {
			return nil, err
		}
		if room.HostId == userId {
			room.IsHost = true
		}
	}
	return rooms, nil
}

func (r *RoomUseCase) ListAllRooms(ctx context.Context, query RoomQuery, page *common.Pagination) ([]*Room, error) {
	rooms, err := r.repo.ListRooms(ctx, query, page)
	if err != nil {
		return nil, err
	}
	for _, room := range rooms {
		if err := r.buildRoomDto(ctx, room); err != nil {
			return nil, err
		}
	}
	return rooms, nil
}

func (r *RoomUseCase) buildRoomDto(ctx context.Context, room *Room) error {
	count, err := r.roomMemberRepo.CountRoomMember(ctx, room.RoomId)
	if err != nil {
		return err
	}
	room.MemberCount = count
	// TODO 获取模拟器信息
	instance, _ := r.roomInstanceRepo.GetRoomInstance(ctx, room.RoomId)
	if instance != nil {
		room.EmulatorId = instance.EmulatorId
		room.EmulatorName = instance.EmulatorName
		room.GameId = instance.GameId
		room.EmulatorType = instance.EmulatorType
		room.EmulatorCode = instance.EmulatorCode
	}
	return nil
}

func (r *RoomUseCase) GetById(ctx context.Context, id int64) (*Room, error) {
	room, err := r.repo.GetById(ctx, id)
	if err != nil {
		return nil, errors.New(500, "Database Error", "查询失败")
	}
	if room == nil {
		return nil, v1.ErrorNotFound("房间不存在")
	}
	if err := r.buildRoomDto(ctx, room); err != nil {
		return nil, err
	}
	return room, nil
}

func (r *RoomUseCase) Join(ctx context.Context, roomId int64, userId int64, password string) error {
	return r.tm.Tx(ctx, func(ctx context.Context) error {
		// 判断房间是否存在
		room, err := r.repo.GetById(ctx, roomId)
		if err != nil {
			r.logger.Error("加入房间获取房间详情错误", err)
			return v1.ErrorServiceError("加入房间出错")
		}
		if room == nil {
			return v1.ErrorNotFound("房间不存在")
		}
		member, _ := r.roomMemberRepo.GetByRoomAndUser(ctx, roomId, userId)
		if member != nil {
			return nil
		}
		// 密码加入，判断密码是否正确
		if room.JoinType == RoomJoinTypePassword {
			hash := sha256.Sum256([]byte(password))
			if room.Password != hex.EncodeToString(hash[:]) {
				return v1.ErrorAccessDenied("房间密码错误")
			}
		}
		// 判断房间成员数量是否已满
		count, err := r.roomMemberRepo.CountRoomMember(ctx, roomId)
		if err != nil {
			r.logger.Error("加入房间获取房间成员数量错误", err)
			return v1.ErrorServiceError("加入房间出错")
		}
		if count >= room.MemberLimit {
			return v1.ErrorAccessDenied("房间已满")
		}
		// 数据库添加成员
		newMember := RoomMember{
			RoomMemberId: r.snowflakeId.Generate().Int64(),
			RoomId:       roomId,
			UserId:       userId,
			Role:         RoomMemberRolePlayer,
			AddTime:      time.Now(),
			Status:       RoomMemberStatusJoined,
		}
		if err := r.roomMemberRepo.Create(ctx, &newMember); err != nil {
			r.logger.Error("加入房间创建成员错误", err)
			return v1.ErrorServiceError("加入房间出错")
		}
		return nil
	})
}

func (r *RoomUseCase) UpdateRoom(ctx context.Context, room *Room, userId int64) error {
	return r.tm.Tx(ctx, func(ctx context.Context) error {
		rm, _ := r.roomMemberRepo.GetByRoomAndUser(ctx, room.RoomId, userId)
		if rm == nil || rm.Role != RoomMemberRoleHost {
			return v1.ErrorAccessDenied("没有权限修改房间")
		}
		if room.JoinType == RoomJoinTypePassword {
			if room.Password == "" {
				return errors.BadRequest("Bad Request", "请填写房间密码")
			}
			hash := sha256.Sum256([]byte(room.Password))
			room.Password = hex.EncodeToString(hash[:])
		}
		if err := r.repo.Update(ctx, room); err != nil {
			r.logger.Error("修改房间错误 ", err)
			return v1.ErrorServiceError("修改房间错误")
		}
		return nil
	})
}

func (r *RoomUseCase) Delete(ctx context.Context, roomId int64, userId int64) error {
	return r.tm.Tx(ctx, func(ctx context.Context) error {
		rm, _ := r.roomMemberRepo.GetByRoomAndUser(ctx, roomId, userId)
		if rm == nil || rm.Role != RoomMemberRoleHost {
			return v1.ErrorAccessDenied("没有权限删除房间")
		}
		// 获取分布式锁，避免删除时有用户正在创建房间实例
		mutexName := openRoomInstanceMutexName(roomId)
		mutex := r.redSync.NewMutex(mutexName, redsync.WithExpiry(time.Second*30))
		if err := mutex.Lock(); err != nil {
			r.logger.Error("删除房间获取分布式锁错误 ", err)
			return v1.ErrorServiceError("删除失败")
		}
		defer mutex.Unlock()
		// 删除房间实体
		if err := r.repo.Delete(ctx, roomId); err != nil {
			r.logger.Error("删除房间错误 ", err)
			return v1.ErrorServiceError("删除失败")
		}
		// 关闭并删除房间实例
		instance, err := r.roomInstanceRepo.GetRoomInstance(ctx, roomId)
		if err != nil {
			r.logger.Error("获取房间实例错误 ", err)
			return v1.ErrorServiceError("删除失败")
		}
		if instance != nil {
			// 关闭房间实例
			if err := r.gameServerRepo.Shutdown(ctx, instance); err != nil {
				r.logger.Error("关闭房间实例错误  ", err)
				return v1.ErrorServiceError("删除失败")
			}
			_ = r.roomInstanceRepo.DeleteRoomInstance(ctx, roomId)
		}

		// 删除所有房间成员
		if err := r.roomMemberRepo.DeleteRoomAllMembers(ctx, roomId); err != nil {
			r.logger.Error("删除所有房间成员错误 ", err)
			return v1.ErrorServiceError("删除失败")
		}
		// 删除所有存档文件
		saveIds, err := r.gameSaveRepo.ListSaveIds(ctx, roomId)
		if err != nil {
			r.logger.Error("获取房间所有存档错误 ", err)
			return v1.ErrorServiceError("删除失败")
		}
		if err := r.gameSaveRepo.DeleteFiles(ctx, saveIds); err != nil {
			r.logger.Error("删除房间存档数据失败 ", roomId, saveIds, err)
		}

		// 删除所有存档
		if err := r.gameSaveRepo.DeleteRoomAllSaves(ctx, roomId); err != nil {
			r.logger.Error("删除所有房间存档错误 ", err)
			return v1.ErrorServiceError("删除失败")
		}

		return nil
	})
}
