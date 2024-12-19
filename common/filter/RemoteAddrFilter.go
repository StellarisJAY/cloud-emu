package filter

import (
	kratos "github.com/go-kratos/kratos/v2/transport/http"
	"net"
	"net/http"
)

type ContextRemoteAddrKey struct{}

type RemoteAddrFilter struct {
	h http.Handler
}

func (r *RemoteAddrFilter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	remoteAddr := request.RemoteAddr
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		r.h.ServeHTTP(writer, request)
		return
	}
	request.Header.Set("X-Remote-Ip", ip)
	r.h.ServeHTTP(writer, request)
}

func NewRemoteAddrFilter() kratos.FilterFunc {
	return func(handler http.Handler) http.Handler {
		return &RemoteAddrFilter{h: handler}
	}
}
