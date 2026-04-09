package repository

import (
	"Go_shortenURL/configs"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const cacheTTL = 24 * time.Hour

type URLRepository struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewURLRepository(dbConfig configs.DatabaseConfig, redisConfig configs.RedisConfig) *URLRepository {
	ctx := context.Background()
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
	)

	db, err := pgxpool.New(ctx, connectionString)

	if err != nil {
		log.Fatalf("connect postgres failed: %v", err)
	}

	if err := ensureSchema(ctx, db); err != nil {
		log.Fatalf("init schema failed: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("ping redis failed: %v", err)
	}

	return &URLRepository{db: db, redis: redisClient}
}

func ensureSchema(ctx context.Context, db *pgxpool.Pool) error {
	const createExtensionQuery = `
		CREATE EXTENSION IF NOT EXISTS pgcrypto
	`

	const createTableQuery = `
		CREATE TABLE IF NOT EXISTS urls (
			id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
			short_code VARCHAR(10) UNIQUE NOT NULL,
			original_url TEXT NOT NULL,
			click_count BIGINT NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`

	const addIDColumnQuery = `
		ALTER TABLE urls
		ADD COLUMN IF NOT EXISTS id UUID DEFAULT gen_random_uuid()
	`

	const addColumnQuery = `
		ALTER TABLE urls
		ADD COLUMN IF NOT EXISTS click_count BIGINT NOT NULL DEFAULT 0
	`

	const addUniqueShortCodeQuery = `
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'urls_short_code_key'
			) THEN
				ALTER TABLE urls ADD CONSTRAINT urls_short_code_key UNIQUE (short_code);
			END IF;
		END
		$$;
	`

	if _, err := db.Exec(ctx, createExtensionQuery); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, createTableQuery); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, addIDColumnQuery); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, addUniqueShortCodeQuery); err != nil {
		return err
	}

	_, err := db.Exec(ctx, addColumnQuery)
	return err
}

func (r *URLRepository) GetURL(ctx context.Context, shortCode string) (string, error) {
	cachedURL, err := r.redis.Get(ctx, shortCode).Result()
	if err == nil {
		go r.incrementClickCount(ctx, shortCode)
		return cachedURL, nil
	}

	if err != redis.Nil {
		return "", err
	}

	var originalURL string
	err = r.db.QueryRow(ctx, "SELECT original_url FROM urls WHERE short_code = $1", shortCode).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("URL not found")
		}
		return "", err
	}

	if err := r.redis.Set(ctx, shortCode, originalURL, cacheTTL).Err(); err != nil {
		log.Printf("set redis cache failed for shortCode=%s: %v", shortCode, err)
	}

	go r.incrementClickCount(ctx, shortCode)

	return originalURL, nil
}

func (r *URLRepository) incrementClickCount(ctx context.Context, shortCode string) {
	_, err := r.db.Exec(ctx, "UPDATE urls SET click_count = click_count + 1 WHERE short_code = $1", shortCode)
	if err != nil {
		log.Printf("update click_count failed for shortCode=%s: %v", shortCode, err)
	}
}

func (r *URLRepository) SetURL(ctx context.Context, shortCode, url string) error {

	_, err := r.db.Exec(ctx, `
		INSERT INTO urls (short_code, original_url)
		VALUES ($1, $2)
		ON CONFLICT (short_code)
		DO UPDATE SET original_url = EXCLUDED.original_url
	`, shortCode, url)
	if err != nil {
		return err
	}

	return r.redis.Set(ctx, shortCode, url, cacheTTL).Err()
}

func (r *URLRepository) GetShortCodeByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	var shortCode string
	err := r.db.QueryRow(ctx, "SELECT short_code FROM urls WHERE original_url = $1 LIMIT 1", originalURL).Scan(&shortCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("URL not found")
		}
		return "", err
	}
	return shortCode, nil
}
