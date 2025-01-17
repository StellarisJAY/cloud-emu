package biz

import (
	"cmp"
	"context"
	"fmt"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redsync/redsync/v4"
	"net/url"
	"slices"
	"strconv"
	"time"
)

type RoomInstanceUseCase struct {
	repo             RoomInstanceRepo
	snowflakeId      *snowflake.Node
	redsync          *redsync.Redsync
	gameServerRepo   GameServerRepo
	roomRepo         RoomRepo
	roomMemberRepo   RoomMemberRepo
	tm               Transaction
	logger           *log.Helper
	emulatorGameRepo EmulatorGameRepo
	emulatorRepo     EmulatorRepo
}

type RoomInstanceRepo interface {
	Create(ctx context.Context, roomInstance *RoomInstance) error
	Update(ctx context.Context, roomInstance *RoomInstance) error
	ListOnlineRoomMembers(ctx context.Context, roomInstance *RoomInstance) ([]int64, error)
	SaveRoomInstance(ctx context.Context, roomInstance *RoomInstance) error
	GetRoomInstance(_ context.Context, roomId int64) (*RoomInstance, error)
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
	GameId         int64     `json:"gameId"`
	SessionKey     string
	EmulatorType   string
	EmulatorCode   string
}

const (
	RoomInstanceStatusActive int32 = iota + 1
	RoomInstanceStatusEnd
)

const (
	OpenRoomInstanceMutexExpire = time.Second * 10
)

func NewRoomInstanceUseCase(repo RoomInstanceRepo, snowflakeId *snowflake.Node, redsync *redsync.Redsync,
	gameServerRepo GameServerRepo, roomRepo RoomRepo, roomMemberRepo RoomMemberRepo, tm Transaction,
	emulatorRepo EmulatorRepo, emulatorGameRepo EmulatorGameRepo, logger log.Logger) *RoomInstanceUseCase {
	return &RoomInstanceUseCase{
		repo:             repo,
		snowflakeId:      snowflakeId,
		redsync:          redsync,
		gameServerRepo:   gameServerRepo,
		roomRepo:         roomRepo,
		roomMemberRepo:   roomMemberRepo,
		tm:               tm,
		emulatorRepo:     emulatorRepo,
		emulatorGameRepo: emulatorGameRepo,
		logger:           log.NewHelper(logger),
	}
}

type OpenRoomInstanceResult struct {
	RoomInstance
	AccessToken string
}

// OpenRoomInstance 创建房间实例，并获取连接token
// TODO 1. 游戏服务器负载均衡; 2. 房间实例记录实时更新; 3. 所有玩家连接断开后释放房间实例资源
func (uc *RoomInstanceUseCase) OpenRoomInstance(ctx context.Context, roomId int64, auth RoomMemberAuth) (*OpenRoomInstanceResult, error) {
	result := &OpenRoomInstanceResult{}
	err := uc.tm.Tx(ctx, func(ctx context.Context) error {
		member, err := uc.roomMemberRepo.GetByRoomAndUser(ctx, roomId, auth.UserId)
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
			uc.logger.Error("创建房间实例分布式锁错误:", err)
			return v1.ErrorServiceError("创建房间实例出错")
		}
		defer mutex.Unlock()
		// 房间实例已经存在
		instance, err := uc.repo.GetRoomInstance(ctx, roomId)
		if instance != nil {
			// 获取连接token
			token, err := uc.gameServerRepo.GetRoomInstanceToken(ctx, instance, roomId, auth)
			if err != nil {
				return v1.ErrorServiceError("获取token失败")
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
			RoomInstanceId: uc.snowflakeId.Generate().Int64(),
			ServerIp:       u.Hostname(),
			RoomId:         roomId,
			AddTime:        time.Now(),
			Status:         RoomInstanceStatusActive,
			EndTime:        time.Now(),
			RpcPort:        int32(port),
		}

		emulator, err := uc.emulatorRepo.GetByType(ctx, "DUMMY")
		if err != nil || emulator == nil {
			return v1.ErrorServiceError("无法获取默认模拟器配置")
		}
		instance.EmulatorId = emulator.EmulatorId
		instance.EmulatorName = emulator.EmulatorName
		instance.EmulatorType = emulator.EmulatorType
		instance.EmulatorCode = emulator.EmulatorCode

		game, err := uc.emulatorGameRepo.GetByEmulatorTypeAndName(ctx, instance.EmulatorType, "gopher.gif")
		if err != nil || game == nil {
			return v1.ErrorServiceError("无法获取默认游戏配置")
		}
		instance.GameId = game.GameId
		gameData, err := uc.emulatorGameRepo.Download(ctx, game)
		if err != nil || len(gameData) == 0 {
			return v1.ErrorServiceError("无法获取默认游戏配置")
		}
		// 在游戏服务器创建房间实例，并获取连接token
		token, sessionKey, err := uc.gameServerRepo.OpenRoomInstance(ctx, instance, auth, gameData)
		if err != nil {
			return v1.ErrorServiceError("创建房间实例出错", err)
		}
		instance.SessionKey = sessionKey
		err = uc.repo.SaveRoomInstance(ctx, instance)
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

func (uc *RoomInstanceUseCase) OpenGameConnection(ctx context.Context, roomId int64, token string, auth RoomMemberAuth) (string, error) {
	mutex := uc.redsync.NewMutex(openRoomInstanceMutexName(roomId), redsync.WithExpiry(OpenRoomInstanceMutexExpire))
	if err := mutex.Lock(); err != nil {
		uc.logger.Error("创建连接分布式锁错误:", err)
		return "", v1.ErrorServiceError("连接失败")
	}
	defer mutex.Unlock()
	instance, err := uc.repo.GetRoomInstance(ctx, roomId)
	if err != nil {
		uc.logger.Error("创建连接获取房间实例错误:", err)
		return "", v1.ErrorServiceError("连接失败")
	}
	if instance == nil {
		uc.logger.Error("创建连接获取房间实例不存在:", roomId)
		return "", v1.ErrorServiceError("连接失败")
	}
	sdpOffer, err := uc.gameServerRepo.OpenGameConnection(ctx, instance, token, auth)
	if err != nil {
		uc.logger.Error("游戏服务创建连接错误:", err)
		return "", err
	}
	return sdpOffer, nil
}

func (uc *RoomInstanceUseCase) SdpAnswer(ctx context.Context, roomId int64, token string, auth RoomMemberAuth, sdpAnswer string) error {
	mutex := uc.redsync.NewMutex(openRoomInstanceMutexName(roomId), redsync.WithExpiry(OpenRoomInstanceMutexExpire))
	if err := mutex.Lock(); err != nil {
		uc.logger.Error("sdpAnswer分布式锁错误:", err)
		return v1.ErrorServiceError("连接失败")
	}
	defer mutex.Unlock()
	instance, err := uc.repo.GetRoomInstance(ctx, roomId)
	if err != nil {
		uc.logger.Error("sdpAnswer获取房间实例错误:", err)
		return v1.ErrorServiceError("连接失败")
	}
	if instance == nil {
		uc.logger.Error("sdpAnswer获取房间实例不存在:", roomId)
		return v1.ErrorServiceError("连接失败")
	}
	err = uc.gameServerRepo.SdpAnswer(ctx, instance, token, auth, sdpAnswer)
	if err != nil {
		uc.logger.Error("游戏服务sdpAnswer错误:", err)
		return err
	}
	return nil
}

func (uc *RoomInstanceUseCase) AddICECandidate(ctx context.Context, roomId int64, token string, auth RoomMemberAuth, candidate string) error {
	mutex := uc.redsync.NewMutex(openRoomInstanceMutexName(roomId), redsync.WithExpiry(OpenRoomInstanceMutexExpire))
	if err := mutex.Lock(); err != nil {
		uc.logger.Error("ice分布式锁错误:", err)
		return v1.ErrorServiceError("连接失败")
	}
	defer mutex.Unlock()
	instance, err := uc.repo.GetRoomInstance(ctx, roomId)
	if err != nil {
		uc.logger.Error("ice获取房间实例错误:", err)
		return v1.ErrorServiceError("连接失败")
	}
	if instance == nil {
		uc.logger.Error("ice获取房间实例不存在:", roomId)
		return v1.ErrorServiceError("连接失败")
	}
	err = uc.gameServerRepo.AddICECandidate(ctx, instance, token, auth, candidate)
	if err != nil {
		uc.logger.Error("游戏服务ice错误:", err)
		return err
	}
	return nil
}

func (uc *RoomInstanceUseCase) GetServerICECandidates(ctx context.Context, roomId int64, token string, auth RoomMemberAuth) ([]string, error) {
	mutex := uc.redsync.NewMutex(openRoomInstanceMutexName(roomId), redsync.WithExpiry(OpenRoomInstanceMutexExpire))
	if err := mutex.Lock(); err != nil {
		uc.logger.Error("ice分布式锁错误:", err)
		return nil, v1.ErrorServiceError("连接失败")
	}
	defer mutex.Unlock()
	instance, err := uc.repo.GetRoomInstance(ctx, roomId)
	if err != nil {
		uc.logger.Error("ice获取房间实例错误:", err)
		return nil, v1.ErrorServiceError("连接失败")
	}
	if instance == nil {
		uc.logger.Error("ice获取房间实例不存在:", roomId)
		return nil, v1.ErrorServiceError("连接失败")
	}
	candidates, err := uc.gameServerRepo.GetServerICECandidate(ctx, instance, token, auth)
	if err != nil {
		uc.logger.Error("游戏服务ice错误:", err)
		return nil, err
	}
	return candidates, nil
}

func (uc *RoomInstanceUseCase) Restart(ctx context.Context, roomId, userId, emulatorId, gameId int64) error {
	return uc.tm.Tx(ctx, func(ctx context.Context) error {
		mutex := uc.redsync.NewMutex(openRoomInstanceMutexName(roomId), redsync.WithExpiry(OpenRoomInstanceMutexExpire))
		if err := mutex.Lock(); err != nil {
			uc.logger.Error("重启获取分布式锁错误:", err)
			return v1.ErrorServiceError("重启失败")
		}
		defer mutex.Unlock()
		// 检查操作人是否是房主
		member, err := uc.roomMemberRepo.GetByRoomAndUser(ctx, roomId, userId)
		if err != nil {
			uc.logger.Error("重启获取房间成员错误:", err)
			return v1.ErrorServiceError("重启失败")
		}
		if member == nil || member.Role != RoomMemberRoleHost {
			return v1.ErrorAccessDenied("当前用户没有权限重启模拟器")
		}
		// 获取房间实例
		instance, err := uc.repo.GetRoomInstance(ctx, roomId)
		if err != nil {
			uc.logger.Error("重启获取房间实例错误:", err)
			return v1.ErrorServiceError("无法重启，请先连接模拟器")
		}
		if instance == nil {
			uc.logger.Error("重启获取房间实例不存在:", roomId)
			return v1.ErrorServiceError("无法重启，请先连接模拟器")
		}
		// 获取游戏信息和模拟器信息
		emulator, _ := uc.emulatorRepo.GetById(ctx, emulatorId)
		if emulator == nil {
			return v1.ErrorServiceError("重启失败，模拟器不存在")
		}
		game, _ := uc.emulatorGameRepo.GetById(ctx, gameId)
		if game == nil || game.EmulatorType != emulator.EmulatorType {
			uc.logger.Error("重启获取游戏信息错误:", game.EmulatorType, emulator.EmulatorType)
			return v1.ErrorServiceError("重启失败，无法加载该游戏")
		}

		params := RestartParams{
			UserId:       userId,
			EmulatorId:   emulatorId,
			GameId:       gameId,
			EmulatorCode: emulator.EmulatorCode,
			GameName:     game.GameName,
			EmulatorType: emulator.EmulatorType,
		}

		// 重启后游戏发生改变，需要获取游戏文件并发送给游戏服务器
		gameData, err := uc.emulatorGameRepo.Download(ctx, game)
		if err != nil {
			uc.logger.Error("重启下载游戏文件错误:", err)
			return v1.ErrorServiceError("重启失败，无法加载该游戏")
		}
		params.GameData = gameData

		// 游戏服务器重启模拟器
		err = uc.gameServerRepo.RestartGameInstance(ctx, instance, params)
		if err != nil {
			uc.logger.Error("重启游戏实例错误:", err)
			return v1.ErrorServiceError("无法重启，请检查模拟器连接")
		}
		// 更新房间实例信息
		instance.EmulatorId = emulator.EmulatorId
		instance.GameId = game.GameId
		instance.EmulatorName = emulator.EmulatorName
		instance.EmulatorType = emulator.EmulatorType
		err = uc.repo.SaveRoomInstance(ctx, instance)
		if err != nil {
			uc.logger.Error("重启更新房间实例错误:", err)
			return v1.ErrorServiceError("重启失败")
		}
		return nil
	})
}

func openRoomInstanceMutexName(roomId int64) string {
	return fmt.Sprintf("/room_instance/open/%d", roomId)
}

func (uc *RoomInstanceUseCase) GetControllerPlayers(ctx context.Context, roomId int64) ([]*ControllerPlayer, error) {
	ri, _ := uc.repo.GetRoomInstance(ctx, roomId)
	if ri == nil {
		return nil, v1.ErrorNotFound("房间不存在")
	}
	return uc.gameServerRepo.GetControllerPlayers(ctx, ri)
}

func (uc *RoomInstanceUseCase) SetControllerPlayers(ctx context.Context, roomId int64, players []*ControllerPlayer, userId int64) error {
	ri, _ := uc.repo.GetRoomInstance(ctx, roomId)
	if ri == nil {
		return v1.ErrorNotFound("房间不存在")
	}
	currUser, _ := uc.roomMemberRepo.GetByRoomAndUser(ctx, roomId, userId)
	if currUser == nil || currUser.Role != RoomMemberRoleHost {
		return v1.ErrorAccessDenied("当前用户没有权限设置控制器")
	}
	for _, pl := range players {
		if pl.UserId == 0 {
			continue
		}
		member, _ := uc.roomMemberRepo.GetByRoomAndUser(ctx, roomId, pl.UserId)
		if member == nil || member.Role == RoomMemberRoleObserver {
			return v1.ErrorServiceError("无法将控制权交给该用户")
		}
	}
	return uc.gameServerRepo.SetControllerPlayer(ctx, players, ri)
}

func (uc *RoomInstanceUseCase) GetGraphicOptions(ctx context.Context, roomId int64) (*GraphicOptions, error) {
	ri, _ := uc.repo.GetRoomInstance(ctx, roomId)
	if ri == nil {
		return nil, v1.ErrorNotFound("房间不存在")
	}
	emulator, _ := uc.emulatorRepo.GetById(ctx, ri.EmulatorId)
	if emulator == nil || !emulator.SupportGraphicSetting {
		return &GraphicOptions{false}, nil
	}
	return uc.gameServerRepo.GetGraphicOptions(ctx, ri)
}

func (uc *RoomInstanceUseCase) SetGraphicOptions(ctx context.Context, roomId int64, opts *GraphicOptions, userId int64) error {
	ri, _ := uc.repo.GetRoomInstance(ctx, roomId)
	if ri == nil {
		return v1.ErrorNotFound("房间不存在")
	}
	currUser, _ := uc.roomMemberRepo.GetByRoomAndUser(ctx, roomId, userId)
	if currUser == nil || currUser.Role != RoomMemberRoleHost {
		return v1.ErrorAccessDenied("没有权限")
	}
	emulator, _ := uc.emulatorRepo.GetById(ctx, ri.EmulatorId)
	if emulator == nil || !emulator.SupportGraphicSetting {
		return errors.New(500, "Service Error", "模拟器不支持画面设置")
	}
	return uc.gameServerRepo.SetGraphicOptions(ctx, ri, opts)
}
