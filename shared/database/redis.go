package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SureshAmal/NimbusU-backend/shared/config"
	"github.com/SureshAmal/NimbusU-backend/shared/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	// Parse Redis URL
	redisURL := cfg.URL
	if !strings.HasPrefix(redisURL, "redis://") {
		redisURL = "redis://" + redisURL
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Redis URL: %w", err)
	}

	// Override with config values if provided
	if cfg.Password != "" {
		opt.Password = cfg.Password
	}
	opt.DB = cfg.DB

	// Create client
	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("unable to ping Redis: %w", err)
	}

	logger.Info("Redis connection established",
		zap.String("url", redisURL),
		zap.Int("db", cfg.DB),
	)

	return client, nil
}

// CloseRedisClient closes the Redis client
func CloseRedisClient(client *redis.Client) error {
	if client != nil {
		err := client.Close()
		if err != nil {
			logger.Error("Error closing Redis client", zap.Error(err))
			return err
		}
		logger.Info("Redis connection closed")
	}
	return nil
}
