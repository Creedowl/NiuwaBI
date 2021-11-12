package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type dbConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	DefaultDB string `mapstructure:"default_db"`
}

type Config struct {
	Debug    bool     `mapstructure:"debug"`
	Host     string   `mapstructure:"host"`
	Port     int      `mapstructure:"port"`
	Database dbConfig `mapstructure:"database"`
}

var Cfg = Config{
	Debug: false,
	Host:  "127.0.0.1",
	Port:  8080,
	Database: dbConfig{
		Host:      "localhost",
		Port:      3306,
		Username:  "root",
		Password:  "password",
		DefaultDB: "niuwaa",
	},
}

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("failed to read config file: %+v\n", err)
	}

	err = viper.Unmarshal(&Cfg)
	if err != nil {
		logrus.Fatalf("failed to unmarshal config file: %+v\n", err)
	}
}
