//go:build wireinject
// +build wireinject

package rungrpc

import (
	"github.com/google/wire"
	"github.com/tuandq2112/go-microservices/shared/logger"
	"github.com/tuandq2112/go-microservices/user-service/infrastructure/repository/memory"
	"github.com/tuandq2112/go-microservices/user-service/infrastructure/server"
	"github.com/tuandq2112/go-microservices/user-service/internal/domain"
	"github.com/tuandq2112/go-microservices/user-service/internal/interfaces/handler"
	"github.com/tuandq2112/go-microservices/user-service/internal/usecase"
)

var userSet = wire.NewSet(
	memory.NewUserRepository,
	usecase.NewUserUsecase,
	handler.NewUserHandler,
	wire.Bind(new(domain.UserRepository), new(*memory.UserRepository)),
)

func InitializeApp() (*server.GRPCServer, error) {
	wire.Build(
		logger.InitLogger,
		userSet,
		wire.Struct(new(server.GRPCServer), "*"),
	)
	return &server.GRPCServer{}, nil
}
