package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	ctx        context.Context
	clinet     *redis.Client
	pipe       redis.Pipeliner
	sub        *redis.PubSub
	channels   map[string]chan []byte
	mu         *sync.Mutex
	channelsMu *sync.RWMutex
}

func NewClient(ctx context.Context, redisConfig *redis.Options, timeSleepForSendBatch time.Duration) *RedisClient {

	client := redis.NewClient(redisConfig)
	pipe := client.Pipeline()
	sub := client.Subscribe(ctx)

	rc := &RedisClient{
		clinet:     client,
		ctx:        ctx,
		pipe:       pipe,
		sub:        sub,
		channels:   make(map[string]chan []byte, redisConfig.MaxActiveConns),
		mu:         &sync.Mutex{},
		channelsMu: &sync.RWMutex{},
	}

	go func() {

		for {

			time.Sleep(timeSleepForSendBatch)

			rc.mu.Lock()
			if _, err := pipe.Exec(ctx); err != nil {
				fmt.Println(err)
			}
			rc.mu.Unlock()

		}

	}()

	go func() {

		for {

			msg, err := sub.Receive(ctx)

			if err != nil {
				return
			}

			switch msg := msg.(type) {

			case *redis.Subscription:

			case *redis.Pong:

			case *redis.Message:

				rc.channelsMu.RLock()

				if rcChan, ok := rc.channels[msg.Channel]; ok {
					rcChan <- []byte(msg.Payload)
				}

				rc.channelsMu.RUnlock()

			}

		}

	}()

	return rc
}

func (c *RedisClient) PublishMessageToSend(receiverUuid string, message []byte) error {

	c.mu.Lock()
	subs := c.pipe.Publish(c.ctx, fmt.Sprintf("chat:message:%s", receiverUuid), message)
	c.mu.Unlock()

	if subs.Err() != nil {
		return subs.Err()
	}

	return nil

}

func (c *RedisClient) SubscribeOnGetMessage(receiverUuid string, receiveCahn chan []byte) error {

	channel := fmt.Sprintf("chat:message:%s", receiverUuid)

	if err := c.sub.Subscribe(c.ctx, channel); err != nil {
		return err
	}

	c.channelsMu.Lock()
	defer c.channelsMu.Unlock()

	if _, ok := c.channels[channel]; ok {
		return nil
	}

	c.channels[channel] = receiveCahn

	return nil
}
