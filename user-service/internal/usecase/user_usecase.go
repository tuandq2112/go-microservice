package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/tuandq2112/go-microservices/user-service/internal/domain"
	customerrors "github.com/tuandq2112/go-microservices/user-service/internal/errors"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (s *UserUsecase) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	return nil, customerrors.ErrInvalidId(map[string]interface{}{
		"user_id": userID,
		"reason":  "empty_id",
	})

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, customerrors.ErrUserNotFound(map[string]interface{}{
				"user_id": userID,
			})
		}
		return nil, customerrors.ErrInternal(map[string]interface{}{
			"error": err.Error(),
		})
	}
	return user, nil
}

func (s *UserUsecase) CreateUser(ctx context.Context, user *domain.User) error {
	if user == nil {
		return customerrors.ErrInvalidId(map[string]interface{}{
			"reason": "nil_user",
		})
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := s.userRepo.Create(ctx, user); err != nil {
		return customerrors.ErrInternal(map[string]interface{}{
			"error": err.Error(),
		})
	}
	return nil
}

func (s *UserUsecase) UpdateUser(ctx context.Context, user *domain.User) error {
	if user == nil || user.ID == "" {
		return customerrors.ErrInvalidId(map[string]interface{}{
			"user_id": user.ID,
			"reason":  "invalid_user",
		})
	}

	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(ctx, user); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return customerrors.ErrUserNotFound(map[string]interface{}{
				"user_id": user.ID,
			})
		}
		return customerrors.ErrInternal(map[string]interface{}{
			"error": err.Error(),
		})
	}
	return nil
}

func (s *UserUsecase) DeleteUser(ctx context.Context, userID string) error {
	if userID == "" {
		return customerrors.ErrInvalidId(map[string]interface{}{
			"user_id": userID,
			"reason":  "empty_id",
		})
	}

	if err := s.userRepo.Delete(ctx, userID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return customerrors.ErrUserNotFound(map[string]interface{}{
				"user_id": userID,
			})
		}
		return customerrors.ErrInternal(map[string]interface{}{
			"error": err.Error(),
		})
	}
	return nil
}
