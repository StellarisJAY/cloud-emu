package biz

import (
	"cmp"
	"context"
	"fmt"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/bwmarrin/snowflake"
	"github.com/go-redsync/redsync/v4"
	"net/url"
	"slices"
	"strconv"
	"time"
)

type RoomInstanceUseCase struct {
	repo           RoomInstanceRepo
	snowflakeId    *snowflake.Node
	redsync        *redsync.Redsync
	gameServerRepo GameServerRepo
	roomRepo       RoomRepo
	roomMemberRepo RoomMemberRepo
	tm             Transaction
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
	ServerIp       string    `json:"serverIp"`
	RpcPort        int32     `json:"rpcPort"`
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
	gameServerRepo GameServerRepo, roomRepo RoomRepo, roomMemberRepo RoomMemberRepo, tm Transaction) *RoomInstanceUseCase {
	return &RoomInstanceUseCase{
		repo:           repo,
		snowflakeId:    snowflakeId,
		redsync:        redsync,
		gameServerRepo: gameServerRepo,
		roomRepo:       roomRepo,
		roomMemberRepo: roomMemberRepo,
		tm:             tm,
	}
}

type OpenRoomInstanceResult struct {
	RoomInstance
	AccessToken string
}

func (uc *RoomInstanceUseCase) OpenRoomInstance(ctx context.Context, roomId int64, userId int64) (*OpenRoomInstanceResult, error) {
	result := &OpenRoomInstanceResult{}
	err := uc.tm.Tx(ctx, func(ctx context.Context) error {
		member, err := uc.roomMemberRepo.GetByRoomAndUser(ctx, roomId, userId)
		if err != nil {
			return v1.ErrorServiceError("获取房间成员出错")
		}
		if member == nil {
			return v1.ErrorAccessDenied("无法访问该房间")
		}
		mutexName := openRoomInstanceMutexName(roomId)
		mutex := uc.redsync.NewMutex(mutexName, redsync.WithExpiry(OpenRoomInstanceMutexExpire))
		// 分布式锁，防止同时创建房间实例
		if err := mutex.Lock(); err != nil {
			return v1.ErrorServiceError("创建房间实例出错")
		}
		defer mutex.Unlock()
		// 房间实例已经存在
		instance, err := uc.repo.GetActiveInstanceByRoomId(ctx, roomId)
		if instance != nil {
			// 获取连接token
			token, err := uc.gameServerRepo.GetRoomInstanceToken(ctx, instance, roomId, userId)
			if err != nil {
				return v1.ErrorServiceError("获取房间实例出错")
			}
			result.RoomInstance = *instance
			result.AccessToken = token
			return nil
		}
		if err != nil {
			return v1.ErrorServiceError("获取房间实例出错")
		}

		// 观战角色无法创建房间实例
		if member.Role == RoomMemberRoleObserver {
			return v1.ErrorAccessDenied("当前用户无法获取房间实例，请等待房主开始游戏")
		}
		// 获取所有可用的游戏服务器
		servers, err := uc.gameServerRepo.ListActiveGameServers(ctx)
		if err != nil {
			return v1.ErrorServiceError("获取可用游戏服务器出错", err)
		}
		if len(servers) == 0 {
			return v1.ErrorServiceError("无法找到可用的游戏服务器")
		}
		// TODO 游戏服务器负载均衡选择策略
		slices.SortStableFunc(servers, func(a, b *GameServer) int {
			return cmp.Compare(a.Weight, b.Weight)
		})
		u, _ := url.Parse(servers[0].Address)
		port, _ := strconv.Atoi(u.Port())
		instance = &RoomInstance{
			ServerIp: u.Hostname(),
			RoomId:   roomId,
			AddTime:  time.Now(),
			Status:   RoomInstanceStatusActive,
			EndTime:  time.Now(),
			RpcPort:  int32(port),
		}

		// 在游戏服务器创建房间实例，并获取连接token
		token, err := uc.gameServerRepo.OpenRoomInstance(ctx, instance)
		if err != nil {
			return v1.ErrorServiceError("创建房间实例出错", err)
		}
		err = uc.repo.Create(ctx, instance)
		if err != nil {
			return v1.ErrorServiceError("创建房间实例出错", err)
		}

		result.RoomInstance = *instance
		result.AccessToken = token
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
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
