package customerrors

import (
	"context"

	"github.com/tuandq2112/go-microservices/shared/interceptors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidIdKey    = "user.invalid_id"
	ErrUserNotFoundKey = "user.not_found"
	ErrInternalKey     = "user.internal_error"
)

// Error represents a custom error with translation support
type Error struct {
	Code codes.Code
	Key  string
	Data map[string]interface{}
	ctx  context.Context
}

func (e *Error) Error() string {
	return interceptors.T(e.ctx, e.Key, e.Data)
}

func (e *Error) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Error())
}

func ErrInvalidId(ctx context.Context, data map[string]interface{}) error {
	return &Error{
		Code: codes.NotFound,
		Key:  ErrInvalidIdKey,
		Data: data,
		ctx:  ctx,
	}
}

func ErrUserNotFound(ctx context.Context, data map[string]interface{}) error {
	return &Error{
		Code: codes.NotFound,
		Key:  ErrUserNotFoundKey,
		Data: data,
		ctx:  ctx,
	}
}

func ErrInternal(ctx context.Context, data map[string]interface{}) error {
	return &Error{
		Code: codes.Internal,
		Key:  ErrInternalKey,
		Data: data,
		ctx:  ctx,
	}
}
