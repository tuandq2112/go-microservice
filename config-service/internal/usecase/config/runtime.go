package config

import (
	"context"

	"github.com/tuandq2112/go-microservices/config-service/internal/domain/config"
)

type ConfigUseCase struct {
	repo config.ConfigRepository
}

func NewConfigUseCase(repo config.ConfigRepository) *ConfigUseCase {
	return &ConfigUseCase{repo: repo}
}

func (c *ConfigUseCase) GetConfig(ctx context.Context, serviceName string, env string) (*config.ConfigModel, error) {
	configData, err := c.repo.GetConfig(ctx, serviceName, env)
	if err != nil {
		return nil, err
	}

	return &config.ConfigModel{
		ServiceName: serviceName,
		Env:         env,
		Config:      configData,
	}, nil
}