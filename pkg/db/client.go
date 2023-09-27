package db

import (
	"context"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	ctx    context.Context
	client redis.Client
}

var lock = &sync.Mutex{}

var instance *Client
var redisClient *redis.Client

func GetClient(ctx context.Context) *Client {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()

		if instance == nil {
			opt, _ := redis.ParseURL(os.Getenv("REDIS_URL"))
			redisClient = redis.NewClient(opt)

			instance = &Client{
				ctx:    ctx,
				client: *redisClient,
			}
		}
	}

	return instance
}

func (c *Client) Store(key, value string) error {
	err := c.client.Set(c.ctx, key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Get(key string) (string, error) {
	val, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
