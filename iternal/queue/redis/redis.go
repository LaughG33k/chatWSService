package redis

import (
	"context"
	"fmt"

	"github.com/LaughG33k/chatWSService/iternal/client/redis"
	"github.com/LaughG33k/chatWSService/iternal/queue"
)

type RedisQueue struct {
	redis *redis.RedisClient
}

// Publish implements queue.Queue.
func (r *RedisQueue) Publish(ctx context.Context, receiver string, message any) error {
	return r.redis.Publish(ctx, fmt.Sprintf("message:%s", receiver), message)
}

// Subscribe implements queue.Queue.
func (r *RedisQueue) Subscribe(ctx context.Context, receiver string, ch chan any) error {
	return r.redis.Subscribe(ctx, fmt.Sprintf("message:%s", receiver), ch)
}

func NewRedisQueue(redisClient *redis.RedisClient) queue.Queue {
	return &RedisQueue{
		redis: redisClient,
	}
}
