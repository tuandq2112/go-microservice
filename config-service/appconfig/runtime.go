package appconfig

import bootstrapLoader "github.com/tuandq2112/go-microservices/shared/loader"

var (
	Port = 50051
	Host = "localhost"
	ConfigFolderPath = "./config"
	ServiceName = "config-service"
)
func InitConfig() {
	bootstrapLoader.InitBootstrap()
	loadAppConfig()
}

func loadAppConfig() {
	cfg := bootstrapLoader.GetConfig()
	Port = cfg.GetInt("PORT")
	Host = cfg.GetString("HOST")
	ConfigFolderPath = cfg.GetString("CONFIG_FOLDER_PATH")
	ServiceName = cfg.GetString("SERVICE_NAME")
}