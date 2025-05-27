package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"context"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func (c *ConfigLoader) GetConfig(ctx context.Context, serviceName string, env string) (map[string]any, error) {
	if config, ok := c.config[GetServiceKey(serviceName, env)]; ok {
		return config, nil
	}

	return nil, errors.New("config not found")
}

func (c *ConfigLoader) LoadConfigFromFolder(folderPath string) error {
	c.logger.Info("Loading config from folder", zap.String("folderPath", folderPath))
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		// Parse filename to get serviceName and env
		// Expected format: servicename.env.yaml
		parts := strings.Split(strings.TrimSuffix(file.Name(), ".yaml"), ".")
		if len(parts) != 2 {
			continue
		}
		serviceName, env := parts[0], parts[1]

		// Read and parse YAML file
		filePath := filepath.Join(folderPath, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", file.Name(), err)
		}

		var config map[string]any
		if err := yaml.Unmarshal(content, &config); err != nil {
			return fmt.Errorf("failed to parse YAML file %s: %v", file.Name(), err)
		}

		// Store in config map
		c.config[GetServiceKey(serviceName, env)] = config
		c.logger.Info("Loaded config", zap.String("serviceName", serviceName), zap.String("env", env), zap.Any("config", config))
	}

	return nil
}

func GetServiceKey(serviceName string, env string) string {
	return fmt.Sprintf("%s.%s", serviceName, env)
}