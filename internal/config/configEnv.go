package config

import (
	"github.com/spf13/viper"
	"go_kit_project/internal/middleware"
	"go_kit_project/internal/static"
)

func ConfigEnv() (error, bool) {
	viper.SetConfigName(static.CONFIG_FILE_NAME)
	viper.AddConfigPath(middleware.RootDir())
	viper.AutomaticEnv()
	viper.SetConfigType(static.CONFIG_FILE_TYPE)
	if err := viper.ReadInConfig(); err != nil {
		return err, false
	}
	return nil, true
}
