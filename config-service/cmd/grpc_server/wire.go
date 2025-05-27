//go:build wireinject
// +build wireinject

package grpcserver

import (
	"github.com/google/wire"
	configInfra "github.com/tuandq2112/go-microservices/config-service/infrastructure/config_loader"
	"github.com/tuandq2112/go-microservices/config-service/infrastructure/server"
	configRepo "github.com/tuandq2112/go-microservices/config-service/internal/domain/config"
	configHandler "github.com/tuandq2112/go-microservices/config-service/internal/interface/handler/config"
	configUseCase "github.com/tuandq2112/go-microservices/config-service/internal/usecase/config"
	"github.com/tuandq2112/go-microservices/shared/discovery"
)

var configSet = wire.NewSet(
	configInfra.NewConfigLoader,
	configUseCase.NewConfigUseCase,
	configHandler.NewConfigHandler,
	wire.Bind(new(configRepo.ConfigRepository), new(*configInfra.ConfigLoader)),
	wire.Bind(new(configHandler.ConfigUseCase), new(*configUseCase.ConfigUseCase)),
)

var discoverySet = wire.NewSet(
	discovery.NewConsulRegistrar,
)

func InitializeServer() (*server.GrpcServer, error) {
	wire.Build(
		configSet,
		discoverySet,
		server.NewGrpcServer,
	)
	return nil, nil
}
