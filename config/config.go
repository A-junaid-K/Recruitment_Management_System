package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort       string `mapstructure:"APP_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBName        string `mapstructure:"DB_NAME"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	Email         string `mapstructure:"EMAIL"`
	EmailPassword string `mapstructure:"EMAIL_PASSWORD"`

	UserAccessTokenSecret     string `mapstructure:"USER_ACCESS_TOKEN_SECRET"`
	UserAccessTokenExpiryHour int    `mapstructure:"USER_ACCESS_TOKEN_EXPIRY_HOUR"`

	AdminAccessTokenSecret     string `mapstructure:"ADMIN_ACCESS_TOKEN_SECRET"`
	AdminAccessTokenExpiryHour int    `mapstructure:"ADMIN_ACCESS_TOKEN_EXPIRY_HOUR"`

	AwsBucket          string `mapstructure:"AWS_BUCKET"`
	AwsRegion          string `mapstructure:"AWS_REGION"`
	AwsAccessKey       string `mapstructure:"AWS_ACCESS_KEY"`
	AwsSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
}

var cfg Config

func InitConfig() Config {

	viper.AddConfigPath("../")
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err.Error())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err.Error())
	}
	return cfg
}

func GetConfig() Config {
	return cfg
}
