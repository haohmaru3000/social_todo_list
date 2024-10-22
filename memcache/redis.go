package memcache

import (
	"context"
	"time"
	"to_do_list/common"

	"github.com/go-redis/cache/v9"
	goservice "github.com/haohmaru3000/go_sdk"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	store *cache.Cache
}

func NewRedisCache(sc goservice.ServiceContext) *redisCache {
	rdClient := sc.MustGet(common.PluginRedis).(*redis.Client)

	c := cache.New(&cache.Options{
		Redis:      rdClient,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &redisCache{store: c}
}

func (rdc *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return rdc.store.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   0,
	})
}

func (rdc *redisCache) Get(ctx context.Context, key string, value interface{}) error {
	return rdc.store.Get(ctx, key, value)
}

func (rdc *redisCache) Delete(ctx context.Context, key string) error {
	return rdc.store.Delete(ctx, key)
}
