package config

import (
	"context"

	"github.com/tuandq2112/go-microservices/config-service/internal/domain/config"
)

type ConfigUseCase interface {
	GetConfig(ctx context.Context, serviceName string, env string) (*config.ConfigModel, error)
}