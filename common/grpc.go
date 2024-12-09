package common

import (
	"google.golang.org/grpc"
	"net/url"
)

func NewGRPCClient(serverUrl string) (*grpc.ClientConn, error) {
	u, err := url.Parse(serverUrl)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(u.Host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
