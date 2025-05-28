package appconfig

import (
	"fmt"

	bootstrapLoader "github.com/tuandq2112/go-microservices/shared/loader"
)

var (
	Host                          string
	Port                          string
	UserServiceHost               string
	UserServicePort               string
	WHITELIST_METHODS_GET_PATH    []string
	WHITELIST_METHODS_POST_PATH   []string
	WHITELIST_METHODS_PUT_PATH    []string
	WHITELIST_METHODS_DELETE_PATH []string
	USER_CONTEXT_KEY              string = "user_context"
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
	if val, ok := configs["USER_SERVICE_HOST"].(string); ok {
		UserServiceHost = val
	}
	if val, ok := configs["USER_SERVICE_PORT"].(string); ok {
		UserServicePort = val
	}
	if val, ok := configs["WHITELIST_METHODS_GET_PATH"].([]interface{}); ok {
		WHITELIST_METHODS_GET_PATH = make([]string, len(val))
		for i, v := range val {
			WHITELIST_METHODS_GET_PATH[i] = fmt.Sprintf("%v", v)
		}
	}
	if val, ok := configs["WHITELIST_METHODS_POST_PATH"].([]interface{}); ok {
		WHITELIST_METHODS_POST_PATH = make([]string, len(val))
		for i, v := range val {
			WHITELIST_METHODS_POST_PATH[i] = fmt.Sprintf("%v", v)
		}
	}
	if val, ok := configs["WHITELIST_METHODS_PUT_PATH"].([]interface{}); ok {
		WHITELIST_METHODS_PUT_PATH = make([]string, len(val))
		for i, v := range val {
			WHITELIST_METHODS_PUT_PATH[i] = fmt.Sprintf("%v", v)
		}
	}
	if val, ok := configs["WHITELIST_METHODS_DELETE_PATH"].([]interface{}); ok {
		WHITELIST_METHODS_DELETE_PATH = make([]string, len(val))
		for i, v := range val {
			WHITELIST_METHODS_DELETE_PATH[i] = fmt.Sprintf("%v", v)
		}
	}
}
