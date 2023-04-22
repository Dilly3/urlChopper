package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dilly3/urlshortner/internal"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

const TIMEOUT time.Duration = time.Duration(time.Second * 10)

func (r *RedisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *RedisRepository) Find(code string) (*internal.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	redirect := internal.Redirect{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.redis.redisDb.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(internal.ErrRedirectNotFound, "repository.redis.redisDb.Find")
	}
	created_at, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.redis.redisDb.Find")
	}

	redirect.Code = data["code"]
	redirect.CreatedAt = created_at
	redirect.Url = data["url"]
	return &redirect, nil
}

func (r *RedisRepository) Store(redirect *internal.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.Url,
		"created_at": redirect.CreatedAt,
	}
	_, err := r.client.HMSet(ctx, key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.redis.redisDb.Store")
	}
	return nil

}

func NewRedisRepository(redisUrl string) (*RedisRepository, error) {
	repo := &RedisRepository{}
	client, err := newRedisClient(redisUrl)
	if err != nil {
		return nil, errors.Wrap(err, "repository.redis.redisDb.NewRedisRepository")
	}
	repo.client = client
	return repo, nil
}

func newRedisClient(redisUrl string) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, errors.Wrap(err, "repository.redis.redisDb.newRedisclient")

	}
	client := redis.NewClient(opts)
	_, err = client.Ping(ctx).Result()
	return client, err
}
