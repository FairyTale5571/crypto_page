package redis

import (
	"context"
	"fmt"

	"github.com/fairytale5571/crypto_page/pkg/logger"
	"github.com/fairytale5571/crypto_page/pkg/storage"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	db     *redis.Client
	logger *logger.LoggerWrapper
}

var ctx = context.Background()

func New(uri string) (*Redis, error) {
	res := &Redis{
		logger: logger.New("redis"),
	}

	opt, err := redis.ParseURL(uri)
	if err != nil {
		res.logger.Fatalf("cant parse redis url %s", err.Error())
		return nil, err
	}
	res.db = redis.NewClient(opt)
	return res, nil
}

func (r *Redis) Get(key string, bucket storage.Bucket) (string, error) {
	return r.db.Get(ctx, fmt.Sprintf("%s::%s", bucket, key)).Result()
}

func (r *Redis) Set(key, value string, bucket storage.Bucket) error {
	return r.db.Set(ctx, fmt.Sprintf("%s::%s", bucket, key), value, 0).Err()
}

func (r *Redis) Delete(id string, payments storage.Bucket) {
	if err := r.db.Del(ctx, fmt.Sprintf("%s::%s", payments, id)).Err(); err != nil {
		r.logger.Errorf("delete redis key %s", err.Error())
	}
}
