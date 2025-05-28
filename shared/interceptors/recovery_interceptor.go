package interceptors

import (
	"context"

	"github.com/tuandq2112/go-microservices/shared/errors"
	"github.com/tuandq2112/go-microservices/shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func UnaryRecoveryInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered in unary interceptor",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
				)
				err = errors.InternalServerError()
			}
		}()
		return handler(ctx, req)
	}
}

func StreamRecoveryInterceptor(logger logger.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var err error
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered in stream interceptor",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
				)
				err = errors.InternalServerError()
			}
		}()
		err = handler(srv, ss)
		return err
	}
}
