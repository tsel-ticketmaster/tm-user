package redis

import (
	"sync"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/tsel-ticketmaster/tm-user/config"
)

var (
	client         redis.UniversalClient
	clientSyncOnce sync.Once
)

func buildConnection() redis.UniversalClient {
	cfg := config.Get()
	c := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    cfg.Redis.Addrs,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	redisotel.InstrumentTracing(c)
	redisotel.InstrumentMetrics(c)

	return c
}

func GetClient() redis.UniversalClient {
	clientSyncOnce.Do(func() {
		client = buildConnection()
	})

	return client
}
