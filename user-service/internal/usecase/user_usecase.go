package usecase

import (
	"context"
	"time"

	"github.com/tuandq2112/go-microservices/user-service/internal/domain"
)

// UserUsecase handles the business logic for user operations
type UserUsecase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase creates a new instance of UserUsecase
func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

// GetUser retrieves a user by ID
func (s *UserUsecase) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// CreateUser creates a new user
func (s *UserUsecase) CreateUser(ctx context.Context, user *domain.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	return s.userRepo.Create(ctx, user)
}

// UpdateUser updates an existing user
func (s *UserUsecase) UpdateUser(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(ctx, user)
}

// DeleteUser deletes a user by ID
func (s *UserUsecase) DeleteUser(ctx context.Context, userID string) error {
	return s.userRepo.Delete(ctx, userID)
}
