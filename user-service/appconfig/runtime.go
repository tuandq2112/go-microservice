package appconfig

import (
	"fmt"

	bootstrapLoader "github.com/tuandq2112/go-microservices/shared/loader"
)

var (
	Host              string
	Port              string
	JWTSecret         string = "SUPER_SECRET"
	WHITELIST_METHODS []string
	USER_CONTEXT_KEY  string = "user_context"
)

func InitConfig() {
	bootstrapLoader.InitBootstrap()
	loadAppConfig()
}

func loadAppConfig() {
	configServiceHost := bootstrapLoader.GetConfig().GetString("CONFIG_SERVICE_HOST")
	configServicePort := bootstrapLoader.GetConfig().GetString("CONFIG_SERVICE_PORT")
	serviceName := bootstrapLoader.GetConfig().GetString("SERVICE_NAME")
	env := bootstrapLoader.GetConfig().GetString("ENV")
	configs, err := bootstrapLoader.LoadConfigFromConfigService(configServiceHost, configServicePort, serviceName, env)
	if err != nil {
		panic(err)
	}
	if val, ok := configs["HOST"].(string); ok {
		Host = val
	}
	if val, ok := configs["PORT"].(string); ok {
		Port = val
	}
	if val, ok := configs["JWT_SECRET"].(string); ok {
		JWTSecret = val
	}

	if val, ok := configs["WHITELIST_METHODS"].([]interface{}); ok {
		WHITELIST_METHODS = make([]string, len(val))
		for i, v := range val {
			WHITELIST_METHODS[i] = fmt.Sprintf("%v", v)
		}
	}
}
