package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
	"github.com/tuandq2112/go-microservices/shared/interceptors"
	"github.com/tuandq2112/go-microservices/shared/logger"
	pb "github.com/tuandq2112/go-microservices/shared/proto/types/user"
	"github.com/tuandq2112/go-microservices/shared/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Mock response
	user := &pb.User{
		Id:        req.UserId,
		Username:  "john_doe",
		Email:     "john@example.com",
		FullName:  "John Doe",
		CreatedAt: 1710000000,
		UpdatedAt: 1710000000,
	}

	return &pb.GetUserResponse{
		User: user,
	}, nil
}

func main() {
	logger := logger.GetLogger()

	g := &run.Group{}

	start, stop := buildGRPCServer(logger)

	g.Add(start, stop)

	signalStart, signalStop := buildSignalHandler(logger)
	g.Add(signalStart, signalStop)

	if err := g.Run(); err != nil {
		logger.Error("Server exited", zap.Error(err))
	}
}

func buildGRPCServer(logger logger.Logger) (func() error, func(error)) {
	unaryInts := []grpc.UnaryServerInterceptor{
		interceptors.UnaryLoggerInterceptor(logger),
		interceptors.UnaryRecoveryInterceptor(logger),
	}
	streamInts := []grpc.StreamServerInterceptor{
		interceptors.StreamLoggerInterceptor(logger),
		interceptors.StreamRecoveryInterceptor(logger),
	}
	_, start, stop, err := server.BuildGRPCServer(
		":50052",
		func(s *grpc.Server) {
			logger.Info("Registering UserServiceServer")
			pb.RegisterUserServiceServer(s, &UserServer{})
		},
		unaryInts,
		streamInts,
	)
	if err != nil {
		return nil, nil
	}

	return start, stop
}

func buildSignalHandler(logger logger.Logger) (func() error, func(error)) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	return func() error {
			sig := <-sigs
			return fmt.Errorf("signal received: %v", sig)
		}, func(err error) {
			close(sigs)
		}
}
