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

type authConfig struct {
	JwtKey     string `mapstructure:"jwt_key"`
	Timeout    int    `mapstructure:"timeout"`
	MaxRefresh int    `mapstructure:"max_refresh"`
}

type Config struct {
	Debug    bool       `mapstructure:"debug"`
	Host     string     `mapstructure:"host"`
	Port     int        `mapstructure:"port"`
	Database dbConfig   `mapstructure:"database"`
	Auth     authConfig `mapstructure:"auth"`
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
	Auth: authConfig{
		JwtKey:     "CHANGE_ME",
		Timeout:    60 * 60,
		MaxRefresh: 60 * 60 * 2,
	},
}

func InitConfig() {
	logrus.Infoln("init config")
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
	if Cfg.Debug {
		logrus.Debugln("config:")
		viper.Debug()
	}
}
