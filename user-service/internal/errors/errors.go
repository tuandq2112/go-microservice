package customerrors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidIdKey    = "user.invalid_id"
	ErrUserNotFoundKey = "user.not_found"
	ErrInternalKey     = "user.internal_error"
)

// Error represents a custom error with translation support

func ErrInvalidId(data map[string]interface{}) error {
	return status.Errorf(codes.NotFound, ErrInvalidIdKey)
}

func ErrUserNotFound(data map[string]interface{}) error {
	return status.Errorf(codes.NotFound, ErrUserNotFoundKey)
}

func ErrInternal(data map[string]interface{}) error {
	return status.Errorf(codes.Internal, ErrInternalKey)
}
