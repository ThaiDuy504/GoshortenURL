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

	// Thử đọc từ tệp config.env (có thể thất bại nếu chạy trong Docker không có file này)
	if err := cleanenv.ReadConfig("configs/config.env", &cfg); err != nil {
		log.Printf("Warning: Could not load config file (configs/config.env): %v. Relying on environment variables.", err)
	}

	// Đọc từ biến môi trường (Ghi đè hoặc bổ sung)
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("Error reading environment variables: ", err)
	}

	return &cfg
}
