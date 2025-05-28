package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Common error message IDs for localization
const (
	ErrUnauthorized     = "error.unauthorized"
	ErrForbidden        = "error.forbidden"
	ErrNotFound         = "error.not_found"
	ErrInternalServer   = "error.internal_server"
	ErrBadRequest       = "error.bad_request"
	ErrValidationFailed = "error.validation_failed"
	ErrConnectionFailed = "error.connection_failed"
)

func UnauthorizedError() error {
	return status.New(codes.Unauthenticated, ErrUnauthorized).Err()
}

func ForbiddenError() error {
	return status.New(codes.PermissionDenied, ErrForbidden).Err()
}

func NotFoundError() error {
	return status.New(codes.NotFound, ErrNotFound).Err()
}

func InternalServerError() error {
	return status.New(codes.Internal, ErrInternalServer).Err()
}

func BadRequestError() error {
	return status.New(codes.InvalidArgument, ErrBadRequest).Err()
}

func ValidationFailedError() error {
	return status.New(codes.InvalidArgument, ErrValidationFailed).Err()
}

func ConnectionFailedError() error {
	return status.New(codes.Unavailable, ErrConnectionFailed).Err()
}
