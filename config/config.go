package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"log"
)

var config *Config

// Config .
type Config struct {
	Port  string `envconfig:"PROJECT_PORT"`
	Debug bool   `envconfig:"DEBUG"`

	MySQL struct {
		Host   string `envconfig:"DB_HOST"`
		Port   string `envconfig:"DB_PORT"`
		DBName string `envconfig:"DB_NAME"`
		User   string `envconfig:"DB_USER"`
		Pass   string `envconfig:"DB_PASSWORD"`
	} `yaml:"mysql"`

	SecretKey string `envconfig:"SECRET_KEY"`
}

func init() {
	config = &Config{}
	// read from env
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load godotenv", err)
	}

	err = envconfig.Process("", config)
	if err != nil {
		err = errors.Wrap(err, "Failed to decode config env")
		log.Println(err)
	}

	// default value
	if len(config.Port) == 0 {
		config.Port = "3000"
	}
}

// GetConfig .
func GetConfig() *Config {
	return config
}
