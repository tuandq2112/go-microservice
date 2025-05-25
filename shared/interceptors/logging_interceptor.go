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

        logger.Info("Unary gRPC",
            zap.String("method", info.FullMethod),
            zap.Any("request", req),
            zap.Duration("duration", time.Since(start)),
            zap.Error(err),
        )

        return resp, err
    }
}
