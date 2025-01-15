package common

import (
	"fmt"
	"google.golang.org/grpc"
)

func NewGRPCClient(ip string, port int) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure(), grpc.WithDefaultCallOptions(
		grpc.MaxCallSendMsgSize(32*1024*1024), grpc.MaxCallRecvMsgSize(32*1024*1024)))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
