package config

import (
	"context"

	configpb "github.com/tuandq2112/go-microservices/shared/proto/types/config"
	"google.golang.org/protobuf/types/known/structpb"
)

type ConfigHandler struct {
	configpb.UnimplementedConfigServiceServer
	configUseCase ConfigUseCase
}

func NewConfigHandler(configUseCase ConfigUseCase) *ConfigHandler {
	return &ConfigHandler{
		configUseCase: configUseCase,
	}
}

func (h *ConfigHandler) GetConfig(ctx context.Context, req *configpb.GetConfigRequest) (*configpb.GetConfigResponse, error) {
	config, err := h.configUseCase.GetConfig(ctx, req.ServiceName, req.Env)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	structValue, err := MapToStruct(config.Config)
	if err != nil {
		return nil, err
	}

	return &configpb.GetConfigResponse{
		Value: structValue,
	}, nil
}

func MapToStruct(m map[string]interface{}) (*structpb.Struct, error) {
	return structpb.NewStruct(m)
}
