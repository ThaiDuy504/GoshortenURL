package configs

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Port string `env:"APP_PORT" env-default:"8080"`
}

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

func NewConfig() *Config {
	var cfg Config
	//using config.env
	err := cleanenv.ReadConfig("configs/config.env", &cfg)
	if err != nil {
		log.Fatal("Error loading config file", err)
	}
	//using env in docker compose
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal("Error reading config file", err)
	}
	return &cfg
}
