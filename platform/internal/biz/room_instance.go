package biz

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/bwmarrin/snowflake"
	"github.com/go-redsync/redsync/v4"
	"slices"
	"time"
)

type RoomInstanceUseCase struct {
	repo           RoomInstanceRepo
	snowflakeId    *snowflake.Node
	redsync        *redsync.Redsync
	gameServerRepo GameServerRepo
	roomRepo       RoomRepo
	roomMemberRepo RoomMemberRepo
}

type RoomInstanceRepo interface {
	Create(ctx context.Context, roomInstance *RoomInstance) error
	Update(ctx context.Context, roomInstance *RoomInstance) error
	GetActiveInstanceByRoomId(ctx context.Context, roomId int64) (*RoomInstance, error)
	ListInstanceByRoomId(ctx context.Context, roomId int64, p *common.Pagination) ([]*RoomInstance, error)
	ListOnlineRoomMembers(ctx context.Context, roomInstance *RoomInstance) ([]*RoomMember, error)
}

type RoomInstance struct {
	RoomInstanceId int64     `json:"roomInstanceId"`
	RoomId         int64     `json:"roomId"`
	RoomName       string    `json:"roomName"`
	HostId         int64     `json:"hostId"`
	HostName       string    `json:"hostName"`
	EmulatorId     int64     `json:"emulatorId"`
	EmulatorName   string    `json:"emulatorName"`
	AddTime        time.Time `json:"addTime"`
	ServerUrl      string    `json:"serverUrl"`
	EndTime        time.Time `json:"endTime"`
	Status         int32     `json:"status"`
}

const (
	RoomInstanceStatusActive int32 = iota
	RoomInstanceStatusEnd
)

const (
	OpenRoomInstanceMutexExpire = time.Second * 10
)

func NewRoomInstanceUseCase(repo RoomInstanceRepo, snowflakeId *snowflake.Node, redsync *redsync.Redsync,
	gameServerRepo GameServerRepo, roomRepo RoomRepo, roomMemberRepo RoomMemberRepo) *RoomInstanceUseCase {
	return &RoomInstanceUseCase{
		repo:           repo,
		snowflakeId:    snowflakeId,
		redsync:        redsync,
		gameServerRepo: gameServerRepo,
		roomRepo:       roomRepo,
		roomMemberRepo: roomMemberRepo,
	}
}

type OpenRoomInstanceResult struct {
	RoomInstance
	AccessToken string
}

func (uc *RoomInstanceUseCase) OpenRoomInstance(ctx context.Context, roomId int64, userId int64) (*OpenRoomInstanceResult, error) {
	member, err := uc.roomMemberRepo.GetByRoomAndUser(ctx, roomId, userId)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.New("not member of this room")
	}
	mutexName := openRoomInstanceMutexName(roomId)
	mutex := uc.redsync.NewMutex(mutexName, redsync.WithExpiry(OpenRoomInstanceMutexExpire))
	// 分布式锁，防止同时创建房间实例
	if err := mutex.Lock(); err != nil {
		return nil, err
	}
	defer mutex.Unlock()
	// 房间实例已经存在
	instance, err := uc.repo.GetActiveInstanceByRoomId(ctx, roomId)
	if instance != nil {
		// 获取连接token
		token, err := uc.gameServerRepo.GetRoomInstanceToken(ctx, instance, roomId, userId)
		if err != nil {
			return nil, err
		}
		return &OpenRoomInstanceResult{
			RoomInstance: *instance,
			AccessToken:  token,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	// 观战角色无法创建房间实例
	if member.Role == RoomMemberRoleObserver {
		return nil, errors.New("no authority")
	}
	// 获取所有可用的游戏服务器
	servers, err := uc.gameServerRepo.ListActiveGameServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list active game servers: %w", err)
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("no active game servers available")
	}
	// TODO 游戏服务器负载均衡选择策略
	slices.SortStableFunc(servers, func(a, b *GameServer) int {
		return cmp.Compare(a.Weight, b.Weight)
	})
	instance = &RoomInstance{
		ServerUrl: servers[0].Address,
		RoomId:    roomId,
		AddTime:   time.Now(),
		Status:    RoomInstanceStatusActive,
		EndTime:   time.Now(),
	}

	// 在游戏服务器创建房间实例，并获取连接token
	token, err := uc.gameServerRepo.OpenRoomInstance(ctx, instance)
	if err != nil {
		return nil, fmt.Errorf("failed to open room instance on target server: %w", err)
	}
	err = uc.repo.Create(ctx, instance)
	if err != nil {
		return nil, fmt.Errorf("failed to create room instance: %w", err)
	}

	return &OpenRoomInstanceResult{
		RoomInstance: *instance,
		AccessToken:  token,
	}, nil
}

func (uc *RoomInstanceUseCase) ListRoomGameHistory(ctx context.Context, roomId int64, p *common.Pagination) ([]*RoomInstance, error) {
	instances, err := uc.repo.ListInstanceByRoomId(ctx, roomId, p)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func openRoomInstanceMutexName(roomId int64) string {
	return fmt.Sprintf("/room_instance/open/%d", roomId)
}
