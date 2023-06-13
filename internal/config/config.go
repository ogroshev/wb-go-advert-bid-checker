package config

import (
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port int
	LogLevel string
}

func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")

    viper.AutomaticEnv()
	viper.BindEnv("Port", "PORT")
	err = viper.BindEnv("LogLevel", "LOG_LEVEL")
	if err != nil {
		log.Warnf("could not bind LogLevel: %s", err)
	}
	
	viper.SetDefault("Port", 80)
	viper.SetDefault("LogLevel", "info")

    viper.ReadInConfig()
	
    err = viper.Unmarshal(&config)
    return
}