package appconfig

import (
	"os"

	bootstrapLoader "github.com/tuandq2112/go-microservices/shared/loader"
)

var (
	Port             = 50051
	Host             = "127.0.0.1"
	ConfigFolderPath = "./config"
	ServiceName      = "config-service"
)

func InitConfig() {
	bootstrapLoader.InitBootstrap()
	loadAppConfig()
}

func loadAppConfig() {
	cfg := bootstrapLoader.GetConfig()

	if configFolderPath := os.Getenv("CONFIG_FOLDER_PATH"); configFolderPath != "" {
		ConfigFolderPath = configFolderPath
	}

	Port = cfg.GetInt("PORT")
	Host = cfg.GetString("HOST")
	ServiceName = cfg.GetString("SERVICE_NAME")
}
