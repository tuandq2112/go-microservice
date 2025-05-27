package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
	"github.com/tuandq2112/go-microservices/shared/interceptors"
	"github.com/tuandq2112/go-microservices/shared/logger"
	pb "github.com/tuandq2112/go-microservices/shared/proto/types/user"
	"github.com/tuandq2112/go-microservices/shared/server"
	"github.com/tuandq2112/go-microservices/user-service/internal/interfaces/handler"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	Handler *handler.UserHandler
	Logger  logger.Logger
}

func (s *GRPCServer) Start() {
	g := &run.Group{}

	start, stop := buildGRPCServer(s.Logger, s.Handler)

	g.Add(start, stop)

	signalStart, signalStop := buildSignalHandler(s.Logger)
	g.Add(signalStart, signalStop)

	if err := g.Run(); err != nil {
		s.Logger.Error("Server exited", zap.Error(err))
	}
}

func buildGRPCServer(logger logger.Logger, handler *handler.UserHandler) (func() error, func(error)) {
	unaryInts := []grpc.UnaryServerInterceptor{
		interceptors.UnaryLoggerInterceptor(logger),
		interceptors.UnaryRecoveryInterceptor(logger),
	}
	streamInts := []grpc.StreamServerInterceptor{
		interceptors.StreamLoggerInterceptor(logger),
		interceptors.StreamRecoveryInterceptor(logger),
	}
	_, start, stop, err := server.BuildGRPCServer(
		":50052",
		func(s *grpc.Server) {
			logger.Info("Registering UserServiceServer")
			pb.RegisterUserServiceServer(s, handler)
		},
		unaryInts,
		streamInts,
	)
	if err != nil {
		return nil, nil
	}

	return start, stop
}

func buildSignalHandler(logger logger.Logger) (func() error, func(error)) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	return func() error {
			sig := <-sigs
			return fmt.Errorf("signal received: %v", sig)
		}, func(err error) {
			close(sigs)
		}
}
