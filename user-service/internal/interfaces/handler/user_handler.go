package handler

import (
	"context"

	pb "github.com/tuandq2112/go-microservices/shared/proto/types/user"
	"github.com/tuandq2112/go-microservices/user-service/internal/usecase"
)

// UserHandler implements the gRPC user service
type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService *usecase.UserUsecase
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser handles the GetUser gRPC request
func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.userService.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return &pb.GetUserResponse{}, nil
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			CreatedAt: user.CreatedAt.Unix(),
			UpdatedAt: user.UpdatedAt.Unix(),
		},
	}, nil
}
