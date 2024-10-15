package redis

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	ctx                   context.Context
	clinet                *redis.Client
	pipe                  redis.Pipeliner
	sub                   *redis.PubSub
	channels              map[string]chan any
	mu                    *sync.Mutex
	channelsMu            *sync.RWMutex
	timeSleepForSendBatch time.Duration
	stop                  chan struct{}
	onceClose             sync.Once
}

func NewClient(ctx context.Context, redisConfig *redis.Options, timeSleepForSendBatch time.Duration) *RedisClient {

	client := redis.NewClient(redisConfig)
	pipe := client.Pipeline()
	sub := client.Subscribe(ctx)

	rc := &RedisClient{
		clinet:                client,
		pipe:                  pipe,
		sub:                   sub,
		channels:              make(map[string]chan any, redisConfig.MaxActiveConns),
		mu:                    &sync.Mutex{},
		channelsMu:            &sync.RWMutex{},
		timeSleepForSendBatch: timeSleepForSendBatch,
		stop:                  make(chan struct{}),
		onceClose:             sync.Once{},
	}

	return rc
}

func (c *RedisClient) Start() {

	go func() {
		if err := c.sendBatch(); err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		if err := c.read(); err != nil {
			fmt.Println(err)
		}
	}()

}

func (c *RedisClient) Close() (err error) {

	c.onceClose.Do(func() {

		time.Sleep(1 * time.Minute)
		close(c.stop)

		if cErr := c.clinet.Close(); cErr != nil {
			err = errors.Join(err, cErr)
		}

		if sErr := c.sub.Close(); sErr != nil {
			err = errors.Join(err, sErr)
		}

	})

	return err
}

func (c *RedisClient) sendBatch() error {

	for {

		if c.isStop() {
			return nil
		}

		time.Sleep(c.timeSleepForSendBatch)

		c.mu.Lock()
		if _, err := c.pipe.Exec(c.ctx); err != nil {
			c.Close()
			return err
		}
		c.mu.Unlock()

	}

}

func (c *RedisClient) read() error {

	for {

		if c.isStop() {
			return nil
		}

		tm, canc := context.WithTimeout(c.ctx, 1*time.Minute)
		defer canc()
		msg, err := c.sub.Receive(tm)

		if err != nil {
			c.Close()
			return err
		}

		switch msg := msg.(type) {

		case *redis.Subscription:

		case *redis.Pong:

		case *redis.Message:

			c.channelsMu.RLock()

			if rcChan, ok := c.channels[msg.Channel]; ok {

				select {
				case rcChan <- []byte(msg.Payload):
				default:

				}

			}

			c.channelsMu.RUnlock()

		}

	}

}

func (c *RedisClient) isStop() bool {

	select {

	case _, ok := <-c.stop:
		if !ok {
			return true
		}

	default:
	}

	return false

}

func (c *RedisClient) Publish(ctx context.Context, channel string, message any) error {

	c.mu.Lock()
	subs := c.pipe.Publish(c.ctx, channel, message)
	c.mu.Unlock()

	if subs.Err() != nil {
		return subs.Err()
	}

	return nil

}

func (c *RedisClient) Subscribe(ctx context.Context, channel string, ch chan any) error {

	if err := c.sub.Subscribe(ctx, channel); err != nil {
		return err
	}

	c.channelsMu.Lock()
	defer c.channelsMu.Unlock()

	if _, ok := c.channels[channel]; ok {
		return nil
	}

	c.channels[channel] = ch

	return nil
}
