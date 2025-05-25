package interceptors

import (
	"time"

	"github.com/tuandq2112/go-microservices/shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func StreamLoggerInterceptor(logger logger.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		err := handler(srv, ss)

		logger.Info("Stream gRPC",
			zap.String("method", info.FullMethod),
			zap.Bool("isClientStream", info.IsClientStream),
			zap.Bool("isServerStream", info.IsServerStream),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)

		return err
	}
}
