package appconfig

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

var cfg *viper.Viper

func InitBootstrap() {
	file := "bootstrap.yaml"

	if envFile := os.Getenv("CONFIG_FILE"); envFile != "" {
		file = envFile
	}
	cfg = viper.New()

	cfg.SetConfigType("yaml")
	cfg.SetConfigName(path.Base(file))
	cfg.AddConfigPath(path.Dir(file))
	cfg.AddConfigPath("../../config/")
	cfg.AddConfigPath("./config/")

	if err := cfg.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func GetConfig() *viper.Viper {
	return cfg
}
