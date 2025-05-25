package config

import "context"

type ConfigRepository interface {
	GetConfig(ctx context.Context, serviceName string, env string) (map[string]any, error)
}