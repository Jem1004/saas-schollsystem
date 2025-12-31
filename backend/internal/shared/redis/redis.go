package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/school-management/backend/internal/config"
)

// Client wraps the Redis client with additional functionality
type Client struct {
	rdb *redis.Client
}

// Connect establishes a connection to Redis
func Connect(cfg config.RedisConfig) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.rdb.Close()
}

// GetClient returns the underlying Redis client
func (c *Client) GetClient() *redis.Client {
	return c.rdb
}

// Queue Operations for Notification System

// QueueName constants
const (
	NotificationQueue = "notifications:queue"
	RetryQueue        = "notifications:retry"
)

// Enqueue adds an item to a queue
func (c *Client) Enqueue(ctx context.Context, queue string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	return c.rdb.RPush(ctx, queue, jsonData).Err()
}

// Dequeue removes and returns an item from a queue (blocking)
func (c *Client) Dequeue(ctx context.Context, queue string, timeout time.Duration) (string, error) {
	result, err := c.rdb.BLPop(ctx, timeout, queue).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // No item available
		}
		return "", fmt.Errorf("failed to dequeue: %w", err)
	}

	if len(result) < 2 {
		return "", nil
	}

	return result[1], nil
}

// DequeueNonBlocking removes and returns an item from a queue (non-blocking)
func (c *Client) DequeueNonBlocking(ctx context.Context, queue string) (string, error) {
	result, err := c.rdb.LPop(ctx, queue).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // No item available
		}
		return "", fmt.Errorf("failed to dequeue: %w", err)
	}

	return result, nil
}

// QueueLength returns the number of items in a queue
func (c *Client) QueueLength(ctx context.Context, queue string) (int64, error) {
	return c.rdb.LLen(ctx, queue).Result()
}

// Cache Operations

// Set stores a value with optional expiration
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.rdb.Set(ctx, key, jsonData, expiration).Err()
}

// Get retrieves a value by key
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	result, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Key not found
		}
		return "", fmt.Errorf("failed to get value: %w", err)
	}

	return result, nil
}

// Delete removes a key
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}

	return result > 0, nil
}

// SetNX sets a value only if the key doesn't exist (for distributed locks)
func (c *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.rdb.SetNX(ctx, key, jsonData, expiration).Result()
}

// Increment increments a counter
func (c *Client) Increment(ctx context.Context, key string) (int64, error) {
	return c.rdb.Incr(ctx, key).Result()
}

// Expire sets expiration on a key
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.rdb.Expire(ctx, key, expiration).Err()
}

// Pub/Sub Operations

// Publish publishes a message to a channel
func (c *Client) Publish(ctx context.Context, channel string, message interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return c.rdb.Publish(ctx, channel, jsonData).Err()
}

// Subscribe subscribes to a channel
func (c *Client) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.rdb.Subscribe(ctx, channels...)
}
