package domain

import (
	"context"
	"errors"
	"time"
)

// Error definitions
var (
	ErrUserNotFound = errors.New("user not found")
)

// User represents the core user entity
type User struct {
	ID        string
	Username  string
	Email     string
	FullName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}
