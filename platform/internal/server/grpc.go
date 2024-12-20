package server

import (
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/platform/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, userSrv v1.UserServer, roomSrv v1.RoomServer,
	roomInstanceSrv v1.RoomInstanceServer, notificationServer v1.NotificationServer, roomMemberServer v1.RoomMemberServer,
	emulatorServer v1.EmulatorServer, buttonLayoutServer v1.ButtonLayoutServer, keyboardBindingServer v1.KeyboardBindingServer,
	logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterUserServer(srv, userSrv)
	v1.RegisterRoomServer(srv, roomSrv)
	v1.RegisterRoomInstanceServer(srv, roomInstanceSrv)
	v1.RegisterNotificationServer(srv, notificationServer)
	v1.RegisterRoomMemberServer(srv, roomMemberServer)
	v1.RegisterEmulatorServer(srv, emulatorServer)
	v1.RegisterButtonLayoutServer(srv, buttonLayoutServer)
	v1.RegisterKeyboardBindingServer(srv, keyboardBindingServer)
	return srv
}
