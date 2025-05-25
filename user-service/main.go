package main

import (
	"context"
	"log"
	"net"

	pb "github.com/tuandq2112/go-microservices/shared/proto/types/user"
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
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, &UserServer{})

	log.Println("User Service starting on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
} 