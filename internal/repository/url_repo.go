package repository

import (
	// "Go_shortenURL/internal/model"
	"Go_shortenURL/configs"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

const (
	TIME_CACHE = 60 * 60 * 24 // 24 hours
)

type URLRepository struct {
	db *pgx.Conn
	redis *redis.Client
}

func NewURLRepository(dbConfig configs.DatabaseConfig, redisConfig configs.RedisConfig) *URLRepository {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)
	db, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		log.Fatal(err)
	}
	redis := redis.NewClient(&redis.Options{
		Addr: redisConfig.Addr,
		Password: redisConfig.Password,
		DB: redisConfig.DB,
	})
	return &URLRepository{db: db, redis: redis}
}


func (r *URLRepository) GetURL(shortCode string) (string, error) {
	url, err := r.redis.Get(context.Background(), shortCode).Result()
	if err == redis.Nil {
		return "", errors.New("URL not found")
	}
	return url, nil
}

func (r *URLRepository) SetURL(shortCode, url string) error {
	return r.redis.Set(context.Background(), shortCode, url, TIME_CACHE).Err()
}