package interceptors

import (
	"github.com/tuandq2112/go-microservices/shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func StreamRecoveryInterceptor(logger logger.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered in stream interceptor",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
				)
			}
		}()
		return handler(srv, ss)
	}
}
