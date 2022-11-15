package redis

import (
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	addresses map[string]string

	Client *redis.Ring
	Cache  *cache.Cache
}

func New(opts ...Option) (*Redis, error) {
	r := &Redis{
		addresses: map[string]string{
			"server1": ":6379",
		},
	}

	for _, opt := range opts {
		opt(r)
	}

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: r.addresses,
	})

	c := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &Redis{
		Client: ring,
		Cache:  c,
	}, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}
