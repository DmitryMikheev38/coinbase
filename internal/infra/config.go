package infra

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	CoinbaseURL string
	DB          ConfigDB
}

func NewConfig() *Config {
	// load values from ..env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No ..env file found")
	}
	viper.AutomaticEnv()

	cfg := &Config{}
	cfg.CoinbaseURL = viper.GetString("COINBASE_URL")
	cfg.DB.Host = viper.GetString("DB_HOST")
	cfg.DB.Port = viper.GetString("DB_PORT")
	cfg.DB.User = viper.GetString("DB_USER")
	cfg.DB.Password = viper.GetString("DB_PWD")
	cfg.DB.Name = viper.GetString("DB_NAME")

	return cfg
}
