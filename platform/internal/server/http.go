package server

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common/filter"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/StellrisJAY/cloud-emu/platform/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	kratosjwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
)

func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/v1.User/Register"] = struct{}{}
	whiteList["/v1.User/Login"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, ac *conf.Auth, userSrv v1.UserServer, roomSrv v1.RoomServer,
	roomInstanceSrv v1.RoomInstanceServer, notificationServer v1.NotificationServer, roomMemberServer v1.RoomMemberServer,
	emulatorServer v1.EmulatorServer, emulatorGameUC *biz.EmulatorGameUseCase, buttonLayoutServer v1.ButtonLayoutServer,
	keyboardBindingServer v1.KeyboardBindingServer, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			selector.Server( // jwt 验证
				kratosjwt.Server(func(token *jwt.Token) (interface{}, error) {
					return []byte(ac.JwtSecret), nil
				}, kratosjwt.WithSigningMethod(jwt.SigningMethodHS256), kratosjwt.WithClaims(func() jwt.Claims {
					return &biz.LoginClaims{}
				})),
			).Match(NewWhiteListMatcher()).Build(),
			logging.Server(logger),
		),
		http.Filter(handlers.CORS( // 浏览器跨域
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.ExposedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		), filter.NewRemoteAddrFilter()),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	route := srv.Route("/api/v1")
	route.POST("/game/upload", emulatorGameUC.Upload)
	v1.RegisterUserHTTPServer(srv, userSrv)
	v1.RegisterRoomHTTPServer(srv, roomSrv)
	v1.RegisterRoomInstanceHTTPServer(srv, roomInstanceSrv)
	v1.RegisterNotificationHTTPServer(srv, notificationServer)
	v1.RegisterRoomMemberHTTPServer(srv, roomMemberServer)
	v1.RegisterEmulatorHTTPServer(srv, emulatorServer)
	v1.RegisterButtonLayoutHTTPServer(srv, buttonLayoutServer)
	v1.RegisterKeyboardBindingHTTPServer(srv, keyboardBindingServer)
	return srv
}
