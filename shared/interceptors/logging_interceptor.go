package interceptors

import (
	"context"
	"time"

	"github.com/tuandq2112/go-microservices/shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func UnaryLoggerInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		logger.Debug("Unary gRPC",
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)

		return resp, err
	}
}

func StreamLoggerInterceptor(logger logger.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		err := handler(srv, ss)

		logger.Debug("Stream gRPC",
			zap.String("method", info.FullMethod),
			zap.Bool("isClientStream", info.IsClientStream),
			zap.Bool("isServerStream", info.IsServerStream),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)

		return err
	}
}
