package main

import (
	"context"
	"os"

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

	unaryInts := []grpc.UnaryServerInterceptor{
		interceptors.UnaryLoggerInterceptor(logger),
		interceptors.UnaryRecoveryInterceptor(logger),
	}
	streamInts := []grpc.StreamServerInterceptor{
		interceptors.StreamLoggerInterceptor(logger),
		interceptors.StreamRecoveryInterceptor(logger),
	}
	g := &run.Group{}
	_, start, stop, err := server.BuildGRPCServer(
		":50052",
		func(s *grpc.Server) {
			pb.RegisterUserServiceServer(s, &UserServer{})
		},
		unaryInts,
		streamInts,
	)
	if err != nil {
		logger.Error("Failed to build grpc server", zap.Error(err))
		os.Exit(1)
	}
	g.Add(func() error {
		logger.Info("User Service starting on :50052")
		return start()
	}, func(err error) {
		logger.Info("Shutting down gRPC server...")
		stop()
	})

	if err := g.Run(); err != nil {
		logger.Error("Server exited", zap.Error(err))
	}
}
