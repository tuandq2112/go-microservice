package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func BuildGRPCServer(
	addr string,
	registerFunc func(*grpc.Server),
	unaryInterceptors []grpc.UnaryServerInterceptor,
	streamInterceptors []grpc.StreamServerInterceptor,
) (
	server *grpc.Server,
	start func() error,
	stop func(err error),
	err error,
) {
	var opts []grpc.ServerOption

	if len(unaryInterceptors) > 0 {
		opts = append(opts, grpc.ChainUnaryInterceptor(unaryInterceptors...))
	}

	if len(streamInterceptors) > 0 {
		opts = append(opts, grpc.ChainStreamInterceptor(streamInterceptors...))
	}

	server = grpc.NewServer(opts...)
	registerFunc(server)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to listen: %w", err)
	}

	start = func() error {
		return server.Serve(lis)
	}

	stop = func(err error) {
		server.GracefulStop()
	}

	return server, start, stop, nil
}
