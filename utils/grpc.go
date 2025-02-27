package utils

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (u *Utils) GrpcConnect(addr string) (*grpc.ClientConn, context.Context, context.CancelFunc) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		u.Logger.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return conn, ctx, cancel
}

func (u *Utils) GrpcDisConnect(conn *grpc.ClientConn, cancel context.CancelFunc) {
	cancel()

	if err := conn.Close(); err != nil {
		u.Logger.Fatalf("Failed to close the gRPC connection to server: %v", err)
	}
}
