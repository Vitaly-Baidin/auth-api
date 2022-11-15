package redis

import (
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	Client *redis.Ring
	Cache  *cache.Cache
}

func New(url string) (*Redis, error) {
	r := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": ":6379", // TODO вынести в конфиг
		},
	})
	c := cache.New(&cache.Options{
		Redis:      r,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &Redis{
		Client: r,
		Cache:  c,
	}, nil
}
