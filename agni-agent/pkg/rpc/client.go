package rpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var rpcconn *grpc.ClientConn

func routerConnect(router string) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		router,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		), // replace with TLS later
	)

	if err != nil {
		panic(err)
	}

	rpcconn = conn

	return conn

}

func GetRouter() *grpc.ClientConn {
	return rpcconn
}
