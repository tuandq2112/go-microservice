package config

import (
	"path/filepath"

	"github.com/tuandq2112/go-microservices/config-service/appconfig"
	"github.com/tuandq2112/go-microservices/shared/logger"
	"go.uber.org/zap"
)

type ConfigLoader struct {
	config map[string]map[string]any
	logger logger.Logger
}	

func NewConfigLoader() *ConfigLoader {
	loader := &ConfigLoader{
		config: make(map[string]map[string]any),
		logger: logger.GetLogger(),
	}

	configPath := filepath.Join(appconfig.ConfigFolderPath)
	if err := loader.LoadConfigFromFolder(configPath); err != nil {
		loader.logger.Error("Failed to load config from folder", zap.Any("error", err))
	}

	return loader
}