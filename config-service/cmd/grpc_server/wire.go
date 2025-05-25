//go:build wireinject
// +build wireinject

package grpcserver

import (
	"github.com/google/wire"
	"github.com/tuandq2112/go-microservices/config-service/internal/interface/grpc"
	"github.com/tuandq2112/go-microservices/shared/discovery"
	configRepo "github.com/tuandq2112/go-microservices/config-service/internal/domain/config"
	configInfra "github.com/tuandq2112/go-microservices/config-service/internal/infrastructure/config"
	configHandler "github.com/tuandq2112/go-microservices/config-service/internal/interface/handler/config"
	configUseCase "github.com/tuandq2112/go-microservices/config-service/internal/usecase/config"
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
func InitializeServer() (*grpc.GrpcServer, error) {
	wire.Build(
		configSet,
		discoverySet,
		grpc.NewGrpcServer,
	)
	return nil, nil
} 