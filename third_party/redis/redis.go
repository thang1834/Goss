package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/thang1834/go-goss/config"
)

func New(cfg config.Cache) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Pass,
		DB:       cfg.Name,
	})
}

func NewCluster(cfg config.Cache) *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: cfg.Hosts,

		// To route commands by latency or randomly, enable one of the following.
		RouteByLatency: true,
		//RouteRandomly: true,
	})
}
