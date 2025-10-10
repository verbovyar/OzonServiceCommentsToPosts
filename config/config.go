package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort            string `mapstructure:"APP_PORT"`
	PersistanceEnabled bool   `mapstructure:"PERSISTANCE_ENABLED"`
	DbHost             string `mapstructure:"DB_HOST"`
	DbPort             string `mapstructure:"DB_PORT"`
	DbUser             string `mapstructure:"DB_USER"`
	DbPassword         string `mapstructure:"DB_PASSWORD"`
	DbName             string `mapstructure:"DB_NAME"`
}

func Load() (config Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../..")
	viper.AddConfigPath("../../config")
	viper.AddConfigPath("../../..")
	viper.AddConfigPath("../../../config")

	err = viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			log.Println("Config file not found")
		} else {
			log.Println("Reading error config file $HOME/")
		}
	}

	viper.AutomaticEnv()

	log.Printf("Config loaded: %s", viper.AllSettings())

	err = viper.Unmarshal(&config)

	return
}
