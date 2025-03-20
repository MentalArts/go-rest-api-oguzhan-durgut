package client

import (
	"context"
	"mentalartsapi/pkg/runner/pb"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RunnerClient struct {
	client pb.RunnerServiceClient
	conn   *grpc.ClientConn
}

func NewRunnerClient(address string) (*RunnerClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewRunnerServiceClient(conn)
	return &RunnerClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *RunnerClient) Close() error {
	return c.conn.Close()
}

func (c *RunnerClient) RunCode(ctx context.Context, code, language string, timeout int32) (*pb.RunResponse, error) {
	req := &pb.RunRequest{
		Code:     code,
		Language: language,
		Timeout:  timeout,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	return c.client.Run(ctx, req)
}
